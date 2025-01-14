package l1

import (
	"context"
	"errors"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"google.golang.org/grpc/status"
)

type (
	ForkUniqueID = gethcommon.Hash
	RollupNumber = uint64
)

type CrossChainStateMachine interface {
	GetRollupData(number RollupNumber) (RollupInfo, error)
	Synchronize() error
	PublishNextBundle() error
	host.Service
}

// crossChainStateMachine - responsible for maintaining a view of the submitted cross chain bundles for the rollups on the L1.
// Whenever a reorg happens, the state machine will revert to the latest known common ancestor rollup.
// Bundles are only submitted after a rollup is pushed on the L1. The state machine will keep track of the latest rollup and the bundles that have been submitted.
type crossChainStateMachine struct {
	latestRollup  RollupInfo
	rollupHistory map[RollupNumber]RollupInfo
	currentRollup RollupNumber

	enclaveClient   common.Enclave
	publisher       host.L1Publisher
	ethClient       ethadapter.EthClient
	mgmtContractLib mgmtcontractlib.MgmtContractLib // Library to handle Management Contract lib operations
	logger          gethlog.Logger
	hostStopper     *stopcontrol.StopControl
}

type RollupInfo struct {
	ForkUID ForkUniqueID
	Number  RollupNumber
}

func NewCrossChainStateMachine(
	publisher host.L1Publisher,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	ethClient ethadapter.EthClient,
	enclaveClient common.Enclave,
	logger gethlog.Logger,
	hostStopper *stopcontrol.StopControl,
) CrossChainStateMachine {
	return &crossChainStateMachine{
		latestRollup: RollupInfo{
			ForkUID: gethutil.EmptyHash,
			Number:  0,
		},
		rollupHistory:   make(map[RollupNumber]RollupInfo),
		currentRollup:   0,
		publisher:       publisher,
		ethClient:       ethClient,
		mgmtContractLib: mgmtContractLib,
		enclaveClient:   enclaveClient,
		logger:          logger,
		hostStopper:     hostStopper,
	}
}

func (c *crossChainStateMachine) Start() error {
	return nil
}

func (c *crossChainStateMachine) Stop() error {
	return nil
}

func (c *crossChainStateMachine) HealthStatus(context.Context) host.HealthStatus {
	errMsg := ""
	if c.hostStopper.IsStopping() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

func (c *crossChainStateMachine) GetRollupData(number RollupNumber) (RollupInfo, error) {
	if number == c.latestRollup.Number {
		return c.latestRollup, nil
	} else if number > c.latestRollup.Number {
		return RollupInfo{}, errutil.ErrNotFound
	} else {
		return c.rollupHistory[number], nil
	}
}

func (c *crossChainStateMachine) PublishNextBundle() error {
	// If all bundles for the rollups have been published, there is nothing to do.
	if c.currentRollup >= c.latestRollup.Number {
		return nil
	}

	// Get the bundle range from the management contract
	nextForkUID, begin, end, err := c.publisher.GetBundleRangeFromManagementContract(big.NewInt(0).SetUint64(c.currentRollup), c.rollupHistory[c.currentRollup].ForkUID)
	if err != nil {
		return err
	}

	data, err := c.GetRollupData(c.currentRollup + 1)
	if err != nil {
		return err
	}
	if data.ForkUID != *nextForkUID {
		return errutil.ErrRollupForkMismatch
	}

	bundle, err := c.enclaveClient.ExportCrossChainData(context.Background(), begin.Uint64(), end.Uint64())
	if err != nil {
		s, ok := status.FromError(err)
		if ok && errors.Is(s.Err(), errutil.ErrCrossChainBundleNoBatches) {
			c.currentRollup++
			return nil
		}
		return err
	}

	alreadyPublished, err := c.IsBundleAlreadyPublished(bundle)
	if err != nil {
		return err
	}

	if alreadyPublished {
		c.currentRollup++
		return nil
	}

	err = c.publisher.PublishCrossChainBundle(bundle, big.NewInt(0).SetUint64(data.Number), data.ForkUID)
	if err != nil {
		return err
	}

	// Move the current rollup to the next rollup
	c.currentRollup++

	return nil
}

func (c *crossChainStateMachine) IsBundleAlreadyPublished(bundle *common.ExtCrossChainBundle) (bool, error) {
	managementContract, err := ManagementContract.NewManagementContract(*c.mgmtContractLib.GetContractAddr(), c.ethClient.EthClient())
	if err != nil {
		return false, err
	}

	return managementContract.IsBundleAvailable(&bind.CallOpts{}, bundle.CrossChainRootHashes)
}

// Synchronize - checks if there are any new rollups or forks and moves the tracking needle to the latest common ancestor.
func (c *crossChainStateMachine) Synchronize() error {
	forkUID, _, _, err := c.publisher.GetBundleRangeFromManagementContract(big.NewInt(0).SetUint64(c.latestRollup.Number), c.latestRollup.ForkUID)
	if err != nil {
		if errors.Is(err, errutil.ErrNoNextRollup) {
			c.logger.Debug("No new rollup or fork found")
			return nil
		}

		if errors.Is(err, errutil.ErrRollupForkMismatch) {
			return c.revertToLatestKnownCommonAncestorRollup()
		}

		c.logger.Error("Failed to get bundle range from management contract", "error", err)
		return err
	}

	c.rollupHistory[c.latestRollup.Number] = c.latestRollup
	c.latestRollup = RollupInfo{
		ForkUID: *forkUID,
		Number:  c.latestRollup.Number + 1,
	}

	c.logger.Info("Synchronized rollup state machine", "latestRollup", c.latestRollup.Number, "forkUID", c.latestRollup.ForkUID.String())
	return nil
}

func (c *crossChainStateMachine) revertToLatestKnownCommonAncestorRollup() error {
	managementContract, err := ManagementContract.NewManagementContract(*c.mgmtContractLib.GetContractAddr(), c.ethClient.EthClient())
	if err != nil {
		return err
	}

	hashBytes, _, err := managementContract.GetUniqueForkID(&bind.CallOpts{}, big.NewInt(0).SetUint64(c.latestRollup.Number))
	if err != nil {
		return err
	}

	var forkHash gethcommon.Hash
	forkHash = gethcommon.BytesToHash(hashBytes[:])

	for forkHash != c.latestRollup.ForkUID {
		// Revert to previous rollup; No need to wipe the map as the synchronization reinserts the latest rollup
		c.latestRollup = c.rollupHistory[c.latestRollup.Number-1] // go to previous rollup

		hashBytes, _, err = managementContract.GetUniqueForkID(&bind.CallOpts{}, big.NewInt(0).SetUint64(c.latestRollup.Number))
		if err != nil {
			return err
		}

		forkHash = gethcommon.BytesToHash(hashBytes[:])
	}

	// Rollback current rollup if it was dumped due to a fork.
	if c.currentRollup > c.latestRollup.Number {
		c.currentRollup = c.latestRollup.Number
	}

	return nil
}
