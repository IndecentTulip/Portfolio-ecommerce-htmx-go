#!/bin/bash

# ======
DB_NAME="ecommerce"
DB_USER="postgres"
DB_HOST="localhost"
# ======

psql -U "$DB_USER" -d "$DB_NAME" -f ./postsql_schema.sql

