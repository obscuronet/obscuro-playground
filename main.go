package main

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/rpc"
)

const url = "ws://127.0.0.1:37500/" // The websocket address for the first Obscuro node in the full network simulation.

func main() {
	var client *rpc.Client
	for {
		var err error
		client, err = rpc.DialWebsocket(context.Background(), url, "")
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	defer client.Close()

	ch := make(chan *types.Log)

	logSubscription := common.LogSubscription{
		Accounts: []*common.SubscriptionAccount{},
	}
	encodedLogSubscription, err := rlp.EncodeToBytes(logSubscription)
	if err != nil {
		panic(err)
	}

	sub, err := client.Subscribe(context.Background(), "eth", ch, "logs", encodedLogSubscription)
	if err != nil {
		panic(err)
	}

	sub.Unsubscribe()

	//for {
	//	select {
	//	case log := <-ch:
	//		println(fmt.Sprintf("Received logs. Block number: %d. Index: %d. Data: %s", log.BlockNumber, log.Index, string(log.Data)))
	//	case err = <-sub.Err():
	//		panic(err)
	//	}
	//}
}
