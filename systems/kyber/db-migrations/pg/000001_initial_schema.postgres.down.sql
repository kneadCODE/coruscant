-- ============================================================================
-- Kyber Accounting System - Drop Initial Schema
-- ============================================================================
-- Database: PostgreSQL 12+
-- This migration drops the complete kyber accounting system schema
-- Order is important due to foreign key dependencies

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS budget_tracking;
DROP TABLE IF EXISTS budget_items;
DROP TABLE IF EXISTS counterparties;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS ledger_users;
DROP TABLE IF EXISTS ledgers;