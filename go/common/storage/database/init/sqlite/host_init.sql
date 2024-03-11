create table if not exists block_host
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    hash        binary(16)      NOT NULL UNIQUE,
    header      blob            NOT NULL,
    rollup_hash binary(32)      NOT NULL
);

create index IDX_BLOCK_HASH_HOST on block_host (hash);

create table if not exists rollup_host
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    hash              binary(16) NOT NULL UNIQUE,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    time_stamp        int        NOT NULL,
    ext_rollup        blob       NOT NULL,
    compression_block binary(32) NOT NULL
);

create index IDX_ROLLUP_HASH_HOST on rollup_host (hash);
create index IDX_ROLLUP_PROOF_HOST on rollup_host (compression_block);
create index IDX_ROLLUP_SEQ_HOST on rollup_host (start_seq, end_seq);

create table if not exists batch_body_host
(
    id          int        NOT NULL primary key,
    content     mediumblob NOT NULL
);

create table if not exists batch_host
(
    sequence       int primary key,
    full_hash      binary(32) NOT NULL,
    hash           binary(16) NOT NULL unique,
    height         int        NOT NULL,
    tx_count       int        NOT NULL,
    header         blob       NOT NULL,
    body_id        int        NOT NULL REFERENCES batch_body_host
);
create index IDX_BATCH_HEIGHT_HOST on batch_host (height);

create table if not exists transactions_host
(
    hash           binary(16) primary key,
    full_hash      binary(32) NOT NULL,
    body_id        int REFERENCES batch_body_host
);

create table if not exists transaction_count
(
    id          int  NOT NULL primary key,
    total       int  NOT NULL
);

insert into transaction_count (id, total)
values (1, 0) on CONFLICT (id) DO NOTHING;