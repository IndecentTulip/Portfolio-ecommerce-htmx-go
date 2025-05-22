#!/bin/bash

DB_NAME="ecommerce"
DB_USER="postgres"
DB_HOST="localhost"

# Terminate connections and drop the database
psql -U "$DB_USER" -h "$DB_HOST" -X -c "
  REVOKE CONNECT ON DATABASE $DB_NAME FROM public;
  SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$DB_NAME';
"

psql -U "$DB_USER" -h "$DB_HOST" -X -c "DROP DATABASE IF EXISTS $DB_NAME;"

echo "Database '$DB_NAME' deleted successfully."

# Create database if it doesn't exist
psql -U "$DB_USER" -h "$DB_HOST" -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
psql -U "$DB_USER" -h "$DB_HOST" -c "CREATE DATABASE $DB_NAME;"

echo "Database '$DB_NAME' created (or already exists)."

