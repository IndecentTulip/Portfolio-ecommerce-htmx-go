#!/bin/bash

rm ./database/products.db
touch ./database/products.db

echo "info related to google auth"
touch ./src/authcred.txt
echo "provide Client ID"
read clientid
echo "$clientid" > ./src/authcred.txt
echo "provide Client secret"
read clientsecret
echo "$clientsecret" >> ./src/authcred.txt

cd ./src/
go run main.go
