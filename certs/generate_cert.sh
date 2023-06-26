#!/bin/bash
FILE_CERT_NAME=localhost
openssl req -x509 \
            -sha512 -days 365 \
            -nodes \
            -newkey rsa:4096 \
            -subj "/CN=my.app/C=TH/L=Bangkok" \
            -addext "subjectAltName = DNS:my.app" \
            -keyout "$FILE_CERT_NAME.key" -out "$FILE_CERT_NAME.crt"
