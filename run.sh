#!/bin/bash

rm ./database/products.db
touch ./database/products.db

touch ./src/authcred.txt
echo "<><><><><><><><><><>"
echo "info related to google auth"
echo "<><><><><><><><><><>"
echo "provide Client ID"
read clientid
echo "$clientid" > ./src/authcred.txt
echo "provide Client secret"
read clientsecret
echo "$clientsecret" >> ./src/authcred.txt
echo "<><><><><><><><><><>"
echo "info related to github auth"
echo "<><><><><><><><><><>"
echo "provide Client ID"
read clientidgit
echo "$clientidgit" >> ./src/authcred.txt
echo "provide Client secret"
read clientsecretgit
echo "$clientsecretgit" >> ./src/authcred.txt


cd ./src/
go run main.go
