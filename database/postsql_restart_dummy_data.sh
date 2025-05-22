#!/bin/bash

# ======
DB_NAME="ecommerce"
DB_USER="postgres"
DB_HOST="localhost"
# ======

./postsql_del_create.sh
psql -U "$DB_USER" -d "$DB_NAME" -f ./postsql_schema.sql
./postsql_insert.sh
psql -U "$DB_USER" -d "$DB_NAME" -f ./postsql_addon.sql
