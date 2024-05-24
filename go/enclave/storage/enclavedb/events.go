package enclavedb

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/tracers"
)

const (
	baseEventsJoin = "from events e join exec_tx extx on e.tx=extx.tx and e.batch=extx.batch  join tx on extx.tx=tx.id join batch b on extx.batch=b.sequence where b.is_canonical=true "
)

func StoreEventLogs(ctx context.Context, dbtx *sql.Tx, receipts []*types.Receipt, batch *core.Batch, stateDB *state.StateDB) error {
	var args []any
	totalLogs := 0
	for _, receipt := range receipts {
		for _, l := range receipt.Logs {
			logArgs, err := logDBValues(ctx, dbtx, l, stateDB)
			if err != nil {
				return err
			}
			args = append(args, logArgs...)
			txId, _ := GetTxId(ctx, dbtx, l.TxHash)
			if txId == 0 {
				args = append(args, nil)
			} else {
				args = append(args, txId)
			}
			args = append(args, batch.SeqNo().Uint64())
			totalLogs++
		}
	}
	if totalLogs > 0 {
		query := "insert into events (topic0,topic1,topic2,topic3,topic4,datablob,log_idx,address,lifecycle_event,rel_address1,rel_address2,rel_address3,rel_address4,tx,batch) values " +
			repeat("(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", ",", totalLogs)
		_, err := dbtx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

// This method stores a log entry together with relevancy metadata
// Each types.Log has 5 indexable topics, where the first one is the event signature hash
// The other 4 topics are set by the programmer
// According to the data relevancy rules, an event is relevant to accounts referenced directly in topics
// If the event is not referring any user address, it is considered a "lifecycle event", and is relevant to everyone
func logDBValues(ctx context.Context, db *sql.Tx, l *types.Log, stateDB *state.StateDB) ([]any, error) {
	// The topics are stored in an array with a maximum of 5 entries, but usually less
	var t0, t1, t2, t3, t4 []byte

	// these are the addresses to which this event might be relevant to.
	var addr1, addr2, addr3, addr4 *gethcommon.Address
	var a1, a2, a3, a4 []byte

	// start with true, and as soon as a user address is discovered, it becomes false
	isLifecycle := true

	// internal variable
	var isUserAccount bool

	n := len(l.Topics)
	if n > 0 {
		t0 = l.Topics[0].Bytes()
	}
	var err error
	// for every indexed topic, check whether it is an end user account
	// if yes, then mark it as relevant for that account
	if n > 1 {
		t1 = l.Topics[1].Bytes()
		isUserAccount, addr1, err = isEndUserAccount(ctx, db, l.Topics[1], stateDB)
		if err != nil {
			return nil, err
		}
		isLifecycle = isLifecycle && !isUserAccount
		if addr1 != nil {
			a1 = addr1.Bytes()
		}
	}
	if n > 2 {
		t2 = l.Topics[2].Bytes()
		isUserAccount, addr2, err = isEndUserAccount(ctx, db, l.Topics[2], stateDB)
		if err != nil {
			return nil, err
		}
		isLifecycle = isLifecycle && !isUserAccount
		if addr2 != nil {
			a2 = addr2.Bytes()
		}
	}
	if n > 3 {
		t3 = l.Topics[3].Bytes()
		isUserAccount, addr3, err = isEndUserAccount(ctx, db, l.Topics[3], stateDB)
		if err != nil {
			return nil, err
		}
		isLifecycle = isLifecycle && !isUserAccount
		if addr3 != nil {
			a3 = addr3.Bytes()
		}
	}
	if n > 4 {
		t4 = l.Topics[4].Bytes()
		isUserAccount, addr4, err = isEndUserAccount(ctx, db, l.Topics[4], stateDB)
		if err != nil {
			return nil, err
		}
		isLifecycle = isLifecycle && !isUserAccount
		if addr4 != nil {
			a4 = addr4.Bytes()
		}
	}

	// normalise the data field to nil to avoid duplicates
	data := l.Data
	if len(data) == 0 {
		data = nil
	}

	return []any{
		t0, t1, t2, t3, t4,
		data, l.Index,
		l.Address.Bytes(),
		isLifecycle,
		a1, a2, a3, a4,
	}, nil
}

func FilterLogs(
	ctx context.Context,
	db *sql.DB,
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	batchHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	queryParams := []any{}
	query := ""
	if batchHash != nil {
		query += " AND b.hash = ? "
		queryParams = append(queryParams, batchHash.Bytes())
	}

	// ignore negative numbers
	if fromBlock != nil && fromBlock.Sign() > 0 {
		query += " AND b.height >= ?"
		queryParams = append(queryParams, fromBlock.Int64())
	}
	if toBlock != nil && toBlock.Sign() > 0 {
		query += " AND b.height <= ?"
		queryParams = append(queryParams, toBlock.Int64())
	}

	if len(addresses) > 0 {
		cond := repeat("(address=?)", " OR ", len(addresses))
		query += " AND (" + cond + ")"
		for _, address := range addresses {
			queryParams = append(queryParams, address.Bytes())
		}
	}
	if len(topics) > 5 {
		return nil, fmt.Errorf("invalid filter. Too many topics")
	}
	if len(topics) > 0 {
		for i, sub := range topics {
			// empty rule set == wildcard
			if len(sub) > 0 {
				topicColumn := fmt.Sprintf("topic%d", i)
				cond := repeat(fmt.Sprintf("(%s=? )", topicColumn), " OR ", len(sub))
				query += " AND (" + cond + ")"
				for _, topic := range sub {
					queryParams = append(queryParams, topic.Bytes())
				}
			}
		}
	}

	return loadLogs(ctx, db, requestingAccount, query, queryParams)
}

func DebugGetLogs(ctx context.Context, db *sql.DB, txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	var queryParams []any

	query := "select rel_address1, rel_address2, rel_address3, rel_address4, lifecycle_event, topic0, topic1, topic2, topic3, topic4, datablob, b.hash, b.height, tx.hash, tx.idx, log_idx, address " +
		baseEventsJoin +
		" AND tx.hash = ? "

	queryParams = append(queryParams, txHash.Bytes())

	result := make([]*tracers.DebugLogs, 0)

	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := tracers.DebugLogs{
			Log: types.Log{
				Topics: []gethcommon.Hash{},
			},
			LifecycleEvent: false,
		}

		var t0, t1, t2, t3, t4 sql.NullString
		var relAddress1, relAddress2, relAddress3, relAddress4 []byte
		err = rows.Scan(
			&relAddress1,
			&relAddress2,
			&relAddress3,
			&relAddress4,
			&l.LifecycleEvent,
			&t0, &t1, &t2, &t3, &t4,
			&l.Data,
			&l.BlockHash,
			&l.BlockNumber,
			&l.TxHash,
			&l.TxIndex,
			&l.Index,
			&l.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("could not load log entry from db: %w", err)
		}

		for _, topic := range []sql.NullString{t0, t1, t2, t3, t4} {
			if topic.Valid {
				l.Topics = append(l.Topics, stringToHash(topic))
			}
		}

		l.RelAddress1 = bytesToAddress(relAddress1)
		l.RelAddress2 = bytesToAddress(relAddress2)
		l.RelAddress3 = bytesToAddress(relAddress3)
		l.RelAddress4 = bytesToAddress(relAddress4)

		result = append(result, &l)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func bytesToAddress(b []byte) *gethcommon.Address {
	if b != nil {
		addr := gethcommon.BytesToAddress(b)
		return &addr
	}
	return nil
}

// Of the log's topics, returns those that are (potentially) user addresses. A topic is considered a user address if:
//   - It has at least 12 leading zero bytes (since addresses are 20 bytes long, while hashes are 32) and at most 22 leading zero bytes
//   - It does not have associated code (meaning it's a smart-contract address)
//   - It has a non-zero nonce (to prevent accidental or malicious creation of the address matching a given topic,
//     forcing its events to become permanently private (this is not implemented for now)
//
// todo - find a more efficient way
func isEndUserAccount(ctx context.Context, db *sql.Tx, topic gethcommon.Hash, stateDB *state.StateDB) (bool, *gethcommon.Address, error) {
	potentialAddr := common.ExtractPotentialAddress(topic)
	if potentialAddr == nil {
		return false, nil, nil
	}
	addrBytes := potentialAddr.Bytes()
	// Check the database if there are already entries for this address
	var count int
	query := "select count(*) from events where (rel_address1=?) OR (rel_address2=?) OR (rel_address3=? ) OR (rel_address4=? )"
	err := db.QueryRowContext(ctx, query, addrBytes, addrBytes, addrBytes, addrBytes).Scan(&count)
	if err != nil {
		// exit here
		return false, nil, err
	}

	if count > 0 {
		return true, potentialAddr, nil
	}

	// TODO A user address must have a non-zero nonce. This prevents accidental or malicious sending of funds to an
	// address matching a topic, forcing its events to become permanently private.
	// if db.GetNonce(potentialAddr) != 0

	// If the address has code, it's a smart contract address instead.
	if stateDB.GetCode(*potentialAddr) == nil {
		return true, potentialAddr, nil
	}

	return false, nil, nil
}

// utility function that knows how to load relevant logs from the database
// todo always pass in the actual batch hashes because of reorgs, or make sure to clean up log entries from discarded batches
func loadLogs(ctx context.Context, db *sql.DB, requestingAccount *gethcommon.Address, whereCondition string, whereParams []any) ([]*types.Log, error) {
	if requestingAccount == nil {
		return nil, fmt.Errorf("logs can only be requested for an account")
	}

	result := make([]*types.Log, 0)
	query := "select topic0, topic1, topic2, topic3, topic4, datablob, b.hash, b.height, tx.hash, tx.idx, log_idx, address" + " " + baseEventsJoin
	var queryParams []any

	// Add relevancy rules
	//  An event is considered relevant to all account owners whose addresses are used as topics in the event.
	//	In case there are no account addresses in an event's topics, then the event is considered relevant to everyone (known as a "lifecycle event").
	query += " AND (lifecycle_event OR (rel_address1=? OR rel_address2=? OR rel_address3=? OR rel_address4=?)) "
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())

	query += whereCondition
	queryParams = append(queryParams, whereParams...)

	query += " order by b.height, tx.idx asc"

	rows, err := db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := types.Log{
			Topics: []gethcommon.Hash{},
		}
		var t0, t1, t2, t3, t4 []byte
		err = rows.Scan(&t0, &t1, &t2, &t3, &t4, &l.Data, &l.BlockHash, &l.BlockNumber, &l.TxHash, &l.TxIndex, &l.Index, &l.Address)
		if err != nil {
			return nil, fmt.Errorf("could not load log entry from db: %w", err)
		}

		for _, topic := range [][]byte{t0, t1, t2, t3, t4} {
			if len(topic) > 0 {
				l.Topics = append(l.Topics, byteArrayToHash(topic))
			}
		}

		result = append(result, &l)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func stringToHash(ns sql.NullString) gethcommon.Hash {
	value, err := ns.Value()
	if err != nil {
		return [32]byte{}
	}
	s, ok := value.(string)
	if !ok {
		return [32]byte{}
	}
	result := gethcommon.Hash{}
	result.SetBytes([]byte(s))
	return result
}

func byteArrayToHash(b []byte) gethcommon.Hash {
	result := gethcommon.Hash{}
	result.SetBytes(b)
	return result
}
