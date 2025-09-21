-- ============================================================================
-- Kyber Accounting System - Initial Database Schema
-- ============================================================================
-- Database: PostgreSQL 12+
-- This migration creates the complete initial schema for the kyber accounting system
-- following industry best practices for PostgreSQL database design.
-- Uses PostgreSQL-specific features: TIMESTAMPTZ, UUID, GIN indexes, partial indexes

-- ============================================================================
-- CORE TABLES
-- ============================================================================

-- Ledgers: Main aggregate root for user's financial ledgers
CREATE TABLE ledgers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(TRIM(name)) > 0),
    description VARCHAR(1000),
    base_currency CHAR(3) NOT NULL CHECK (LENGTH(TRIM(base_currency)) > 0),
    status VARCHAR(20) NOT NULL CHECK (LENGTH(TRIM(status)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Ledgers indexes
CREATE INDEX idx_ledgers_status ON ledgers(status);
CREATE INDEX idx_ledgers_active ON ledgers(id, base_currency) WHERE status = 'ACTIVE';

-- Ledgers comment
COMMENT ON TABLE ledgers IS 'Main aggregate root for users financial ledgers with multi-currency support';

-- Ledger Users: RBAC for ledger access
CREATE TABLE ledger_users (
    ledger_id UUID NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (LENGTH(TRIM(role)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (ledger_id, user_id)
);

-- Ledger users indexes
CREATE INDEX idx_ledger_users_user_id ON ledger_users(user_id);
CREATE INDEX idx_ledger_users_role ON ledger_users(role);

-- Ledger users comment
COMMENT ON TABLE ledger_users IS 'RBAC table defining user access and roles within ledgers';

-- Accounts: Financial accounts within ledgers
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(TRIM(name)) > 0),
    description VARCHAR(1000),
    type VARCHAR(30) NOT NULL CHECK (LENGTH(TRIM(type)) > 0),
    currency CHAR(3) NOT NULL CHECK (LENGTH(TRIM(currency)) > 0),
    balance_amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (LENGTH(TRIM(status)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (ledger_id, name)
);

-- Accounts indexes
CREATE INDEX idx_accounts_ledger_id ON accounts(ledger_id);
CREATE INDEX idx_accounts_type ON accounts(type);
CREATE INDEX idx_accounts_status ON accounts(status);
CREATE INDEX idx_accounts_ledger_status_type ON accounts(ledger_id, status, type);
CREATE INDEX idx_accounts_active ON accounts(ledger_id, id) WHERE status = 'ACTIVE';

-- Accounts comment
COMMENT ON TABLE accounts IS 'Financial accounts within ledgers (checking, savings, credit cards, loans, etc.)';

-- Counterparties: People or organizations involved in transactions
CREATE TABLE counterparties (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(TRIM(name)) > 0),
    type VARCHAR(30) NOT NULL CHECK (LENGTH(TRIM(type)) > 0),
    description VARCHAR(1000),
    status VARCHAR(20) NOT NULL CHECK (LENGTH(TRIM(status)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (ledger_id, name)
);

-- Counterparties indexes
CREATE INDEX idx_counterparties_ledger_id ON counterparties(ledger_id);
CREATE INDEX idx_counterparties_type ON counterparties(type);
CREATE INDEX idx_counterparties_status ON counterparties(status);
CREATE INDEX idx_counterparties_name ON counterparties(name);
CREATE INDEX idx_counterparties_ledger_status_type ON counterparties(ledger_id, status, type);
CREATE INDEX idx_counterparties_active ON counterparties(ledger_id, id) WHERE status = 'ACTIVE';

-- Counterparties comment
COMMENT ON TABLE counterparties IS 'People or organizations involved in financial transactions';

-- Budget Items: Budget categories for income, expenses, transfers
CREATE TABLE budget_items (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL CHECK (LENGTH(TRIM(name)) > 0),
    description VARCHAR(1000),
    type VARCHAR(20) NOT NULL CHECK (LENGTH(TRIM(type)) > 0),
    currency CHAR(3) NOT NULL CHECK (LENGTH(TRIM(currency)) > 0),
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (ledger_id, name)
);

-- Budget items indexes
CREATE INDEX idx_budget_items_ledger_id ON budget_items(ledger_id);
CREATE INDEX idx_budget_items_type ON budget_items(type);
CREATE INDEX idx_budget_items_is_active ON budget_items(is_active);
CREATE INDEX idx_budget_items_ledger_active_type ON budget_items(ledger_id, is_active, type);
CREATE INDEX idx_budget_items_active ON budget_items(ledger_id, id, type) WHERE is_active = TRUE;

-- Budget items comment
COMMENT ON TABLE budget_items IS 'Budget items representing planned income, expenses, or transfers within ledgers';

-- Budget Tracking: Monthly budget vs actual tracking
CREATE TABLE budget_tracking (
    item_id UUID NOT NULL REFERENCES budget_items(id) ON DELETE CASCADE,
    year INTEGER NOT NULL CHECK (year >= 1900 AND year <= 3000),
    month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    target_amount BIGINT NOT NULL CHECK (target_amount >= 0),
    budgeted_amount BIGINT NOT NULL CHECK (budgeted_amount >= 0),
    actual_amount BIGINT NOT NULL CHECK (actual_amount >= 0),
    updated_at TIMESTAMPTZ NOT NULL,

    PRIMARY KEY (item_id, year, month)
);

-- Budget tracking indexes
CREATE INDEX idx_budget_tracking_year_month ON budget_tracking(year, month);
CREATE INDEX idx_budget_tracking_item_year ON budget_tracking(item_id, year);
CREATE INDEX idx_budget_tracking_year_month_item ON budget_tracking(year, month, item_id);

-- Budget tracking comment
COMMENT ON TABLE budget_tracking IS 'Monthly budget tracking data for budget items';

-- Transactions: Financial transactions
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    ledger_id UUID NOT NULL REFERENCES ledgers(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES budget_items(id) ON DELETE RESTRICT,
    counterparty_id UUID REFERENCES counterparties(id) ON DELETE SET NULL,
    amount BIGINT NOT NULL CHECK (amount != 0),
    description VARCHAR(500) NOT NULL CHECK (LENGTH(TRIM(description)) > 0),
    notes VARCHAR(2000),
    transaction_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Transactions indexes (most important for performance)
CREATE INDEX idx_transactions_ledger_id ON transactions(ledger_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_item_id ON transactions(item_id);
CREATE INDEX idx_transactions_transaction_date ON transactions(transaction_date);

-- Composite indexes for common query patterns
CREATE INDEX idx_transactions_ledger_date ON transactions(ledger_id, transaction_date DESC);
CREATE INDEX idx_transactions_account_date ON transactions(account_id, transaction_date DESC);
CREATE INDEX idx_transactions_item_date ON transactions(item_id, transaction_date DESC);
CREATE INDEX idx_transactions_counterparty_date ON transactions(counterparty_id, transaction_date DESC) WHERE counterparty_id IS NOT NULL;

-- Performance indexes for specific use cases
CREATE INDEX idx_transactions_recent ON transactions(ledger_id, transaction_date DESC, account_id)
    WHERE transaction_date >= NOW() - INTERVAL '2 years';
CREATE INDEX idx_transactions_monthly ON transactions(ledger_id, DATE_TRUNC('month', transaction_date), item_id);
CREATE INDEX idx_transactions_account_balance_calc ON transactions(account_id, transaction_date ASC, amount);

-- Transactions comment
COMMENT ON TABLE transactions IS 'Financial transactions representing money movements between accounts and budget items';

