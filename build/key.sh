#!/bin/bash

set -e

pushd `dirname $0` > /dev/null
SCRIPTPATH=`pwd -P`
popd > /dev/null
SCRIPTFILE=`basename $0`
CERTPATH=${SCRIPTPATH}/../certs
mkdir -p ${CERTPATH}

cd ${CERTPATH}

DAYS=3650

# generate a self-signed rootCA file that would be used to sign both the server and client cert.
# Alternatively, we can use different CA files to sign the server and client, but for our use case, we would use a single CA.
openssl req -newkey rsa:2048 \
  -new -nodes -x509 \
  -days ${DAYS} \
  -out ca.crt \
  -keyout ca.key \
  -subj "/C=SO/ST=Earth/L=Mountain/O=MegaEase/OU=MegaCloud/CN=localhost"

#create a key for server
openssl genrsa -out server.key 2048

#generate the Certificate Signing Request
openssl req -new -key server.key -days ${DAYS} -out server.csr \
  -subj "/C=SO/ST=Earth/L=Mountain/O=MegaEase/OU=MegaCloud/CN=localhost"

#sign it with Root CA
# https://stackoverflow.com/questions/64814173/how-do-i-use-sans-with-openssl-instead-of-common-name
openssl x509  -req -in server.csr \
  -extfile <(printf "subjectAltName=DNS:localhost") \
  -CA ca.crt -CAkey ca.key  \
  -days ${DAYS} -sha256 -CAcreateserial \
  -out server.crt

rm server.csr
#cat server.crt server.key > server.pem