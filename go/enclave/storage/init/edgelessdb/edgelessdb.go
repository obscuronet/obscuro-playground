package edgelessdb

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"embed"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ten-protocol/go-ten/go/common/log"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/storage/init/migration"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common/httputil"

	"github.com/ten-protocol/go-ten/go/enclave/core/egoutils"

	"github.com/go-sql-driver/mysql"
)

/*
   The TEN Enclave (TE) needs a way to persist data into a trusted database. Trusted not to reveal that data to anyone but that particular enclave.

   To achieve this, the TE must first perform Remote Attestation (RA), which gives it confidence that it is connected to
	a trusted version of software running on trusted hardware. The result of this process is a Certificate which can be
	used to set up a trusted TLS connection into the database.

   The next step is to configure the database schema and users in such a way that the TE knows that the db engine will
	only allow itself access to it. This is achieved by creating a "Manifest" file that contains the SQL init code and a
	DBClient Certificate that is known only to the TE.

	This "DBClient" Cert is used by the database to authenticate that it is communicating to the entity that has initialised that schema.

	--------

	In more detail :
   The edgeless DB is unusable when it's first started, as an owner, we must initially:
     - do a remote attestation on the report provided by {edbAddress}/quote
     - create a ca cert that will authenticate our DB users going forwards
     - prepare a manifest.json that contains that CA cert and some SQL to initialise the DB tables and user
     - submit that manifest.json file to {edbAddress}/manifest, using the certificate provided from /quote to authenticate
     - seal and persist the manifest.json and the certs so we can reconnect if enclave is restarted

   When connecting to an already-initialized edgeless DB we must:
     - perform remote attestation on the edgeless db
     - unseal the manifest.json and get the hash of it, also unseal the certificate that edb was initialized with
	 - verify the /signature request for edgeless DB matches the manifest.json hash
     - connect to edb with the persisted cert - it's now safe to read and write to the DB

	Some useful documentation for this:
		Main edb docs: https://docs.edgeless.systems/edgelessdb/#/
		EDB demo docs: https://github.com/edgelesssys/edgelessdb/tree/main/demo
		// Note: due to an issue with the dependency, I've duplicated the relevant parts of the ERA tool into edb_attestation.go
		ERA - remote attestation tool: https://github.com/edgelesssys/era
*/

const (
	edbHTTPPort          = "8080"
	edbManifestEndpoint  = "/manifest"
	edbSignatureEndpoint = "/signature"

	dataDir         = "/data"
	certIssuer      = "tenCA"
	certSubject     = "tenUser"
	enclaveHostName = "enclave"

	dbUser = "ten"
	dbName = "tendb"

	// change this flag to true to debug issues with edgeless DB (and start EDB process with -e EDG_EDB_DEBUG=1
	//   this will give you:
	//  	- verbose logging on EDB
	//		- write the edb.pem file out for connecting to Edgeless DB services manually
	//		- versions of files created with a '.unsealed' suffix that can be used to connect to the database using mysql-client
	debugMode = false

	initFile = "001_init.sql"

	// maximum number of database connections maintained in the tx pool
	maxDBPoolSize = 100
)

var (
	//go:embed *.sql
	sqlFiles embed.FS

	edbCredentialsFilepath = filepath.Join(dataDir, "edb-credentials.json")

	edgelessDBStartTimeout = 60 * time.Second
)

type manifest struct {
	SQL   []string `json:"sql"`
	Cert  string   `json:"ca"`
	Debug bool     `json:"debug"`
}

// todo (#1474) - move more of the hardcoded config into this (attestation conf, usernames etc.)
type Config struct {
	Host string
}

type Credentials struct {
	ManifestJSON string // contains CA cert and sql statements to initialize edb and then to verify edb is setup as expected
	EDBCACertPEM string // root cert securely provided by edb enclave to encrypt all our communication with it
	CACertPEM    string // root cert we generate in our enclave and securely provide to the edb in the manifest
	UserCertPEM  string // db user cert, generated in our enclave, signed by our root cert
	UserKeyPEM   string // db user private key, generated in our enclave
}

// Connector (re-)establishes a connection to the Edgeless DB for the TEN enclave
func Connector(edbCfg *Config, config enclaveconfig.EnclaveConfig, logger gethlog.Logger) (enclavedb.EnclaveDB, error) {
	// rather than fail immediately if EdgelessDB is not available yet we wait up for `edgelessDBStartTimeout` for it to be available
	err := waitForEdgelessDBToStart(edbCfg.Host, logger)
	if err != nil {
		return nil, err
	}

	// load credentials from encrypted persistence if available, otherwise perform handshake and initialization to prepare them
	edbCredentials, err := getHandshakeCredentials(config, edbCfg, logger)
	if err != nil {
		return nil, err
	}

	tlsCfg, err := CreateTLSCfg(edbCredentials)
	if err != nil {
		return nil, err
	}

	sqlDB, err := ConnectToEdgelessDB(edbCfg.Host, tlsCfg, logger)
	if err != nil {
		return nil, err
	}

	// perform db migration
	err = migration.DBMigration(sqlDB, sqlFiles, logger.New(log.CmpKey, "DB_MIGRATION"))
	if err != nil {
		return nil, err
	}

	// wrap it in our eth-compatible key-value store layer
	return enclavedb.NewEnclaveDB(sqlDB, sqlDB, config, logger)
}

func waitForEdgelessDBToStart(edbHost string, logger gethlog.Logger) error {
	start := time.Now()
	edgelessHTTPAddr := fmt.Sprintf("%s:%s", edbHost, edbHTTPPort)
	logger.Info("Waiting to ensure Edgeless DB is available for http requests...")
	var conn net.Conn
	var err error
	for time.Since(start) < edgelessDBStartTimeout {
		conn, err = net.DialTimeout("tcp", edgelessHTTPAddr, time.Second)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("waited for %s but EdgelessDB http server (%s) was still unavailable - %w",
		edgelessDBStartTimeout, edgelessHTTPAddr, err)
}

func getHandshakeCredentials(enclaveConfig enclaveconfig.EnclaveConfig, edbCfg *Config, logger gethlog.Logger) (*Credentials, error) {
	// if we have previously performed the handshake we can retrieve the creds from disk and proceed
	edbCreds, found, err := LoadCredentialsFromFile()
	if err != nil {
		return nil, err
	}
	if !found {
		// they don't exist on disk so we have to perform the handshake and set them up
		edbCreds, err = performHandshake(enclaveConfig, edbCfg, logger)
		if err != nil {
			return nil, err
		}
	}

	return edbCreds, nil
}

// LoadCredentialsFromFile returns (credentials object, found flag, error), if file not found it will return nil error but found=false
func LoadCredentialsFromFile() (*Credentials, bool, error) {
	b, err := egoutils.ReadAndUnseal(edbCredentialsFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to read and unseal credentials file - %w", err)
	}
	var edbCreds *Credentials
	err = json.Unmarshal(b, &edbCreds)
	if err != nil {
		return nil, false, err
	}

	return edbCreds, true, nil
}

func performHandshake(enclaveConfig enclaveconfig.EnclaveConfig, edbCfg *Config, logger gethlog.Logger) (*Credentials, error) {
	// we need to make sure this dir exists before we start read/writing files in there
	err := os.MkdirAll(dataDir, 0o644)
	if err != nil {
		return nil, err
	}

	// Before we try to connect to the Edgeless DB we have to do Remote Attestation (RA) on it
	// the RA will ensure that we are connecting to a database that will not leak any data.
	// The RA will return a Certificate which we'll use for the TLS mutual authentication when we connect to the database.
	// The trust path is as follows:
	// 1. The TEN Enclave performs RA on the database enclave, and the RA object contains a certificate which only the database enclave controls.
	// 2. Connecting to the database via mutually authenticated TLS using the above certificate, will give the TEN enclave confidence that it is only giving data away to some code and hardware it trusts.
	edbPEM, err := performEDBRemoteAttestation(enclaveConfig, edbCfg.Host, defaultEDBConstraints, logger)
	if err != nil {
		return nil, err
	}

	// client used to make secure HTTP requests to Edgeless DB using the ca-cert we have received
	edbHTTPClient, err := httputil.CreateTLSHTTPClient(edbPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare http client from EdgelessDB cert PEM - %w", err)
	}
	caCertPEM, userCertPEM, userKeyPEM, err := prepareCerts()
	if err != nil {
		return nil, err
	}

	edbInitFile, err := sqlFiles.ReadFile(initFile)
	if err != nil {
		logger.Crit("Could not read the initialisation sql file", log.ErrKey, err)
	}

	manifest := &manifest{
		SQL:   createManifestFormat(string(edbInitFile)),
		Cert:  caCertPEM,
		Debug: true,
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest to json - %w", err)
	}
	logger.Info("Initialise edgelessdb with script", "script", string(manifestJSON))
	err = initialiseEdgelessDB(edbCfg.Host, manifest, edbHTTPClient, logger)
	if err != nil {
		return nil, err
	}

	edbCreds := &Credentials{
		EDBCACertPEM: edbPEM,
		CACertPEM:    caCertPEM,
		UserCertPEM:  userCertPEM,
		UserKeyPEM:   userKeyPEM,
		ManifestJSON: string(manifestJSON),
	}
	edbCredsJSON, err := json.Marshal(edbCreds)
	if err != nil {
		return nil, err
	}
	// todo (#1377) - the credentials must be sealed with the enclave unique ID in production, not just the product key
	err = egoutils.SealAndPersist(string(edbCredsJSON), edbCredentialsFilepath, true)
	if err != nil {
		return nil, err
	}
	if debugMode {
		unsealedFile, _ := os.Create(edbCredentialsFilepath + ".unsealed")
		_, err = unsealedFile.WriteString(string(edbCredsJSON))
		if err != nil {
			return nil, fmt.Errorf("failed to write unsealed credentials file when debug is enabled - %w", err)
		}
		_ = unsealedFile.Close()
	}

	return edbCreds, nil
}

func createManifestFormat(content string) (result []string) {
	lines := strings.Split(content, ";")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		trimmed = strings.ReplaceAll(trimmed, "\n", " ")
		if len(trimmed) > 0 {
			result = append(result, trimmed)
		}
	}
	return
}

func CreateTLSCfg(creds *Credentials) (*tls.Config, error) {
	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM([]byte(creds.EDBCACertPEM)); !ok {
		return nil, fmt.Errorf("failed to append edb cert to mysql CA cert pool")
	}
	cert, err := tls.X509KeyPair([]byte(creds.UserCertPEM), []byte(creds.UserKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare keypair from cert and key - %w", err)
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func prepareCerts() (string, string, string, error) {
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		IsCA:                  true,
		BasicConstraintsValid: true,
		// this subject must match the subject authorised in the manifest.json
		Subject:   pkix.Name{CommonName: certIssuer},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		DNSNames:  []string{enclaveHostName},
	}
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate key for CA cert to init Edgeless DB - %w", err)
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create CA cert - %w", err)
	}
	caPEMBuf := new(bytes.Buffer)
	err = pem.Encode(caPEMBuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create ca cert pem - %w", err)
	}

	userCert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		// the issuer and subject have to match those submitted in manifest.json
		Issuer:    pkix.Name{CommonName: certIssuer},
		Subject:   pkix.Name{CommonName: certSubject},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
	}
	certPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate private key for user cert - %w", err)
	}
	userCertBytes, err := x509.CreateCertificate(rand.Reader, userCert, caCert, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to prepare user certificate - %w", err)
	}

	userCertPEM := new(bytes.Buffer)
	err = pem.Encode(userCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: userCertBytes,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("failed to PEM encode user certificate - %w", err)
	}

	certKeyPEM := new(bytes.Buffer)
	privKeyOut, err := x509.MarshalPKCS8PrivateKey(certPrivKey)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to marshal cert priv key - %w", err)
	}
	err = pem.Encode(certKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyOut,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("failed to pem encode the user private key - %w", err)
	}

	return caPEMBuf.String(), userCertPEM.String(), certKeyPEM.String(), nil
}

// initialiseEdgelessDB sends a manifest over http to the edgeless DB with its initial config
func initialiseEdgelessDB(edbHost string, manifest *manifest, httpClient *http.Client, logger gethlog.Logger) error {
	b, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest json - %w", err)
	}
	url := fmt.Sprintf("https://%s:%s%s", edbHost, edbHTTPPort, edbManifestEndpoint)
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("faild to create manifest initialization req - %w", err)
	}

	_, err = httputil.ExecuteHTTPReq(httpClient, req)
	if err != nil {
		return fmt.Errorf("manifest initialization req failed - %w", err)
	}

	// initializing the DB takes sometime as it restarts itself (seems to be typically around 10 seconds)

	maxRetries := 12 // one minute with 5sec waits
	attempts := 0
	for ; attempts < maxRetries; attempts++ {
		time.Sleep(5 * time.Second)
		logger.Info(fmt.Sprintf("Verifying edgeless DB has initialized correctly - attempt %d", attempts))
		err = verifyEdgelessDB(edbHost, manifest, httpClient, logger)
		if err == nil {
			logger.Info("Edgeless DB initialized successfully.")
			break
		}
	}

	if err != nil {
		// give up - output the last seen error
		return fmt.Errorf("failed to verify Edgeless DB after %d attempts - %w", attempts, err)
	}

	return nil
}

// verifyEdgelessDB requests the /signature from the edb, it should match the hash of the manifest we expected
func verifyEdgelessDB(edbHost string, m *manifest, httpClient *http.Client, logger gethlog.Logger) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest to json - %w", err)
	}
	h := sha256.Sum256(b)
	expectedHash := hex.EncodeToString(h[:])

	url := fmt.Sprintf("https://%s:%s%s", edbHost, edbHTTPPort, edbSignatureEndpoint)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return fmt.Errorf("faild to create edb signature req - %w", err)
	}

	edbHash, err := httputil.ExecuteHTTPReq(httpClient, req)
	if err != nil {
		return fmt.Errorf("failed to receive edbHash from /signature request - %w", err)
	}
	if expectedHash != string(edbHash) {
		return fmt.Errorf("hash from edb /signature request didn't match expected hash of manifest.json, expected=%s, found=%s", expectedHash, edbHash)
	}
	logger.Info(fmt.Sprintf("EDB signature matched the expected hash from our manifest (%s)", expectedHash))

	return nil
}

func ConnectToEdgelessDB(edbHost string, tlsCfg *tls.Config, logger gethlog.Logger) (*sql.DB, error) {
	err := mysql.RegisterTLSConfig("custom", tlsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to register tls config for mysql connection - %w", err)
	}
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = edbHost
	cfg.User = dbUser
	cfg.DBName = dbName
	cfg.TLSConfig = "custom"
	dsn := cfg.FormatDSN()
	logger.Info(fmt.Sprintf("Configuring mysql connection: %s", dsn))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mysql connection to edb - %w", err)
	}
	db.SetMaxOpenConns(maxDBPoolSize)
	return db, nil
}
