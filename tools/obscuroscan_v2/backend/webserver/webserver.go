package webserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type WebServer struct {
	engine      *gin.Engine
	backend     *backend.Backend
	bindAddress string
	logger      log.Logger
	server      *http.Server
}

func New(backend *backend.Backend, bindAddress string, logger log.Logger) *WebServer {
	r := gin.New()
	r.RedirectTrailingSlash = false
	gin.SetMode(gin.ReleaseMode)

	// todo this should be reviewed as anyone can access the api right now
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	r.Use(cors.New(config))

	server := &WebServer{
		engine:      r,
		backend:     backend,
		bindAddress: bindAddress,
		logger:      logger,
	}

	// routes
	r.GET("/health/", server.health)
	r.GET("/count/contracts/", server.getTotalContractCount)
	r.GET("/count/transactions/", server.getTotalTransactionCount)
	r.GET("/items/batch/latest/", server.getLatestBatch)
	r.GET("/items/rollup/latest/", server.getLatestRollupHeader)
	r.GET("/batch/:hash", server.getBatch)
	r.GET("/block/:hash", server.getBatch)
	r.GET("/tx/:hash", server.getTransaction)
	r.GET("/items/transactions/", server.getPublicTransactions)
	r.GET("/items/batches/", server.getBatchListing)
	r.GET("/items/blocks/", server.getBlockListing)

	return server
}

func (w *WebServer) Start() error {
	w.server = &http.Server{
		Addr:              w.bindAddress,
		Handler:           w.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling
	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// todo don't panic
			panic(err)
		}
	}()

	return nil
}

func (w *WebServer) Stop() error {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return w.server.Shutdown(ctx)
}

func (w *WebServer) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"healthy": true})
}

func (w *WebServer) getTotalContractCount(c *gin.Context) {
	count, err := w.backend.GetTotalContractCount()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (w *WebServer) getTotalTransactionCount(c *gin.Context) {
	count, err := w.backend.GetTotalTransactionCount()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (w *WebServer) getLatestBatch(c *gin.Context) {
	batch, err := w.backend.GetLatestBatch()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getLatestRollupHeader(c *gin.Context) {
	block, err := w.backend.GetLatestRollupHeader()
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": block})
}

func (w *WebServer) getBatch(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetBatch(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getTransaction(c *gin.Context) {
	hash := c.Param("hash")
	parsedHash := gethcommon.HexToHash(hash)
	batch, err := w.backend.GetTransaction(parsedHash)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": batch})
}

func (w *WebServer) getPublicTransactions(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	publicTxs, err := w.backend.GetPublicTransactions(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": publicTxs})
}

func (w *WebServer) getBatchListing(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	batchesListing, err := w.backend.GetBatchesListing(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": batchesListing})
}

func (w *WebServer) getBlockListing(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	sizeStr := c.DefaultQuery("size", "10")

	offset, err := strconv.ParseUint(offsetStr, 10, 32)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	parseUint, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	batchesListing, err := w.backend.GetBlockListing(offset, parseUint)
	if err != nil {
		errorHandler(c, fmt.Errorf("unable to execute request %w", err), w.logger)
		return
	}

	for _, data := range batchesListing.BlocksData {
		if data.RollupHash.Hex() != "0x0000000000000000000000000000000000000000000000000000000000000000" {
			fmt.Println("output found rollupHash has at block: ", data.BlockHeader.Number.Int64(), " - ", data.RollupHash)
		}
	}

	fmt.Println("Fetched block starting ")

	c.JSON(http.StatusOK, gin.H{"result": batchesListing})
}

func errorHandler(c *gin.Context, err error, logger log.Logger) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
	logger.Error(err.Error())
}
