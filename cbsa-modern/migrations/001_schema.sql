-- CBSA Modern Schema
-- Replaces: Db2 ACCOUNT/PROCTRAN/CONTROL tables + VSAM CUSTOMER file

CREATE TABLE IF NOT EXISTS customers (
    id              BIGSERIAL PRIMARY KEY,
    customer_number VARCHAR(10) UNIQUE NOT NULL,
    sort_code       VARCHAR(6)  NOT NULL,
    title           VARCHAR(10) NOT NULL DEFAULT '',
    first_name      VARCHAR(32) NOT NULL DEFAULT '',
    last_name       VARCHAR(32) NOT NULL DEFAULT '',
    address_line1   VARCHAR(60) NOT NULL DEFAULT '',
    address_line2   VARCHAR(60) NOT NULL DEFAULT '',
    address_line3   VARCHAR(40) NOT NULL DEFAULT '',
    date_of_birth   DATE        NOT NULL DEFAULT CURRENT_DATE,
    credit_score    INTEGER     NOT NULL DEFAULT 0,
    review_date     DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts (
    id                 BIGSERIAL PRIMARY KEY,
    account_number     VARCHAR(8)    UNIQUE NOT NULL,
    sort_code          VARCHAR(6)    NOT NULL,
    customer_number    VARCHAR(10)   NOT NULL REFERENCES customers(customer_number) ON DELETE CASCADE,
    account_type       VARCHAR(8)    NOT NULL,
    interest_rate      DECIMAL(6,2)  NOT NULL DEFAULT 0,
    opened_date        DATE          NOT NULL DEFAULT CURRENT_DATE,
    overdraft_limit    INTEGER       NOT NULL DEFAULT 0,
    last_statement     DATE          NOT NULL DEFAULT CURRENT_DATE,
    next_statement     DATE          NOT NULL DEFAULT CURRENT_DATE,
    available_balance  DECIMAL(12,2) NOT NULL DEFAULT 0,
    actual_balance     DECIMAL(12,2) NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_accounts_customer ON accounts(sort_code, customer_number);

CREATE TABLE IF NOT EXISTS processed_transactions (
    id             BIGSERIAL PRIMARY KEY,
    sort_code      VARCHAR(6)    NOT NULL,
    account_number VARCHAR(8)    NOT NULL,
    txn_date       DATE          NOT NULL,
    txn_time       TIME          NOT NULL,
    reference      VARCHAR(12)   NOT NULL DEFAULT '',
    txn_type       VARCHAR(3)    NOT NULL,
    description    VARCHAR(40)   NOT NULL DEFAULT '',
    amount         DECIMAL(12,2) NOT NULL,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_proctran_account ON processed_transactions(sort_code, account_number);

CREATE TABLE IF NOT EXISTS controls (
    name       VARCHAR(32) PRIMARY KEY,
    value_num  INTEGER     NOT NULL DEFAULT 0,
    value_str  VARCHAR(40) NOT NULL DEFAULT ''
);

-- Sequences replacing CICS Named Counters
CREATE SEQUENCE IF NOT EXISTS customer_number_seq START WITH 1;
CREATE SEQUENCE IF NOT EXISTS account_number_seq START WITH 1;

-- Seed control records (equivalent to base COBOL install)
INSERT INTO controls (name, value_num, value_str)
VALUES
    ('ACCOUNT-LAST', 0, ''),
    ('CUSTOMER-LAST', 0, '')
ON CONFLICT (name) DO NOTHING;
