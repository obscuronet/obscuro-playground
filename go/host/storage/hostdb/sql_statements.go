package hostdb

// SQLStatements struct holds SQL statements for a specific database type
type SQLStatements struct {
	InsertBatch        string
	InsertTransactions string
	InsertTxCount      string
	InsertRollup       string
	InsertBlock        string
	Pagination         string
	Placeholder        string
}

func SQLiteSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:        "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES (?, ?, ?, ?, ?)",
		InsertTransactions: "REPLACE INTO transaction_host (hash, full_hash, b_sequence) VALUES (?, ?, ?)",
		InsertTxCount:      "INSERT INTO transaction_count (id, total) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET total = EXCLUDED.total",
		InsertRollup:       "INSERT INTO rollup_host (hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block) values (?,?,?,?,?,?)",
		InsertBlock:        "REPLACE INTO block_host (hash, header, rollup_hash) values (?,?,?)",
		Pagination:         "LIMIT ? OFFSET ?",
		Placeholder:        "?",
	}
}

func PostgresSQLStatements() *SQLStatements {
	return &SQLStatements{
		InsertBatch:        "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES ($1, $2, $3, $4, $5)",
		InsertTransactions: "INSERT INTO transaction_host (hash, full_hash, b_sequence) VALUES ($1, $2, $3) ON CONFLICT (hash) DO NOTHING",
		InsertTxCount:      "INSERT INTO transaction_count (id, total) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET total = EXCLUDED.total",
		InsertRollup:       "INSERT INTO rollup_host (hash, start_seq, end_seq, time_stamp, ext_rollup, compression_block) values ($1, $2, $3, $4, $5, $6)",
		InsertBlock:        "INSERT INTO block_host (hash, header, rollup_hash) VALUES ($1, $2, $3) ON CONFLICT (hash) DO NOTHING",
		Pagination:         "LIMIT $1 OFFSET $2",
		Placeholder:        "$1",
	}
}
