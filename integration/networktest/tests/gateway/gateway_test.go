package gateway

import (
	"math/big"
	"testing"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/networktest/actions"
	"github.com/ten-protocol/go-ten/integration/networktest/env"
)

var _transferAmount = big.NewInt(100_000_000)

// TestGatewayHappyPath tests ths same functionality as the smoke_test but with the gateway:
// 1. Create two test users
// 2. Allocate funds to the first user
// 3. Send funds from the first user to the second
// 4. Verify the second user has the funds
// 5. Verify the first user has the funds deducted
// To run this test with a local network use the flag to start it with the gateway enabled.
func TestGatewayHappyPath(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(
		"gateway-happy-path",
		t,
		env.LocalDevNetwork(env.WithTenGateway()),
		actions.Series(
			&actions.CreateTestUser{UserID: 0, UseGateway: true},
			&actions.CreateTestUser{UserID: 1, UseGateway: true},
			actions.SetContextValue(actions.KeyNumberOfTestUsers, 2),

			&actions.AllocateFaucetFunds{UserID: 0},
			actions.SnapshotUserBalances(actions.SnapAfterAllocation), // record user balances (we have no guarantee on how much the network faucet allocates)

			&actions.SendNativeFunds{FromUser: 0, ToUser: 1, Amount: _transferAmount},

			&actions.VerifyBalanceAfterTest{UserID: 1, ExpectedBalance: _transferAmount},
			&actions.VerifyBalanceDiffAfterTest{UserID: 0, Snapshot: actions.SnapAfterAllocation, ExpectedDiff: big.NewInt(0).Neg(_transferAmount)},
		),
	)
}
