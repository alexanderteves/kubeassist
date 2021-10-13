#!/bin/bash

openssl genrsa -out ca.key 2048

openssl req -new -x509 -subj "/CN=CA" -extensions v3_ca -days 3650 -key ca.key -sha256 -out ca.crt -config localhost.cnf

openssl genrsa -out localhost.key 2048

openssl req -subj "/CN=localhost" -extensions v3_req -sha256 -new -key localhost.key -out localhost.csr

openssl x509 -req -extensions v3_req -days 3650 -sha256 -in localhost.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out localhost.crt -extfile localhost.cnf
