package genesis

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/storage/init/sqlite"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ten-protocol/go-ten/integration/datagenerator"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const testLogs = "../.build/tests/"

func TestDefaultGenesis(t *testing.T) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "unit",
		TestSubtype: "genesis",
		LogLevel:    gethlog.LvlInfo,
	})

	gen, err := New("")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if len(gen.Accounts) != 3 {
		t.Fatal("unexpected number of accounts")
	}

	backingDB, err := sqlite.CreateTemporarySQLiteDB("", "", testlog.Logger())
	if err != nil {
		t.Fatalf("unable to create temp db: %s", err)
	}
	storageDB := storage.NewStorage(backingDB, nil, gethlog.New())
	stateDB, err := gen.applyAllocations(storageDB)
	if err != nil {
		t.Fatalf("unable to apply genesis allocations")
	}

	if TestnetGenesis.Accounts[0].Amount.Cmp(stateDB.GetBalance(TestnetGenesis.Accounts[0].Address)) != 0 {
		t.Fatalf("unexpected balance")
	}
}

func TestCustomGenesis(t *testing.T) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "unit",
		TestSubtype: "genesis",
		LogLevel:    gethlog.LvlInfo,
	})

	addr1 := datagenerator.RandomAddress()
	amt1 := datagenerator.RandomUInt64()
	addr2 := datagenerator.RandomAddress()
	amt2 := datagenerator.RandomUInt64()

	gen, err := New(
		fmt.Sprintf(
			`{"Accounts": [
				{"Address": "%s", "Amount": %d},
				{"Address": "%s", "Amount": %d}	] }
				`,
			addr1.Hex(), amt1, addr2.Hex(), amt2))
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if len(gen.Accounts) != 2 {
		t.Fatal("unexpected number of accounts")
	}

	backingDB, err := sqlite.CreateTemporarySQLiteDB("", "", testlog.Logger())
	if err != nil {
		t.Fatalf("unable to create temp db: %s", err)
	}
	storageDB := storage.NewStorage(backingDB, nil, gethlog.New())
	stateDB, err := gen.applyAllocations(storageDB)
	if err != nil {
		t.Fatalf("unable to apply genesis allocations")
	}

	if big.NewInt(int64(amt1)).Cmp(stateDB.GetBalance(addr1)) != 0 {
		t.Fatalf("unexpected balance")
	}
	if big.NewInt(int64(amt2)).Cmp(stateDB.GetBalance(addr2)) != 0 {
		t.Fatalf("unexpected balance")
	}
}
