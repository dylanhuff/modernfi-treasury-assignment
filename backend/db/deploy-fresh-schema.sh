#!/bin/bash

# ============================================================================
# Treasury Management Application - Fresh Schema Deployment Script
# ============================================================================
# This script deploys the consolidated schema and seed data to a PostgreSQL
# database. It drops and recreates the database for a clean start.
#
# Usage:
#   ./deploy-fresh-schema.sh
# ============================================================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="treasury_db"
DB_USER="postgres"
DB_PASS="postgres"

echo -e "${GREEN}Deploying to LOCAL database...${NC}"

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SCHEMA_FILE="$SCRIPT_DIR/schema.sql"
SEED_FILE="$SCRIPT_DIR/seed.sql"

# Check if files exist
if [ ! -f "$SCHEMA_FILE" ]; then
    echo -e "${RED}Error: schema.sql not found at $SCHEMA_FILE${NC}"
    exit 1
fi

if [ ! -f "$SEED_FILE" ]; then
    echo -e "${RED}Error: seed.sql not found at $SEED_FILE${NC}"
    exit 1
fi

# Function to run psql commands
run_psql() {
    docker exec treasury_postgres psql -U "$DB_USER" "$@"
}

# Function to run psql with file
run_psql_file() {
    local file=$1
    docker exec -i treasury_postgres psql -U "$DB_USER" -d "$DB_NAME" < "$file"
}

echo -e "${YELLOW}Step 1: Terminating active connections to database...${NC}"
run_psql -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$DB_NAME' AND pid <> pg_backend_pid();" || true

echo -e "${YELLOW}Step 2: Dropping existing database (if exists)...${NC}"
run_psql -c "DROP DATABASE IF EXISTS $DB_NAME;" || true

echo -e "${YELLOW}Step 3: Creating fresh database...${NC}"
run_psql -c "CREATE DATABASE $DB_NAME WITH ENCODING='UTF8' LC_COLLATE='en_US.utf8' LC_CTYPE='en_US.utf8';"

echo -e "${YELLOW}Step 4: Applying consolidated schema...${NC}"
run_psql_file "$SCHEMA_FILE"

echo -e "${YELLOW}Step 5: Applying seed data...${NC}"
run_psql_file "$SEED_FILE"

echo -e "${YELLOW}Step 6: Verifying deployment...${NC}"

# Verify users
echo -e "\n${GREEN}Users:${NC}"
run_psql -d "$DB_NAME" -c "SELECT id, name, balance FROM users ORDER BY id;"

# Verify transaction counts
echo -e "\n${GREEN}Transaction counts by user:${NC}"
run_psql -d "$DB_NAME" -c "SELECT user_id, COUNT(*) as transaction_count FROM transactions GROUP BY user_id ORDER BY user_id;"

# Verify holdings counts
echo -e "\n${GREEN}Holdings counts by user:${NC}"
run_psql -d "$DB_NAME" -c "SELECT user_id, COUNT(*) as holdings_count FROM holdings GROUP BY user_id ORDER BY user_id;"

# Verify total counts
echo -e "\n${GREEN}Total record counts:${NC}"
run_psql -d "$DB_NAME" -c "SELECT 'users' as table_name, COUNT(*) as count FROM users UNION ALL SELECT 'transactions', COUNT(*) FROM transactions UNION ALL SELECT 'holdings', COUNT(*) FROM holdings;"

echo -e "\n${GREEN}✓ Deployment complete!${NC}"
echo -e "${GREEN}Database: $DB_NAME on $DB_HOST${NC}"
