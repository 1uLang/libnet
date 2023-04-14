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

function generate_client() {
  CLIENT=$1
  O=$2
  OU=$3
  openssl genrsa -out ${CLIENT}.key 2048
  openssl req -new -key ${CLIENT}.key -days ${DAYS} -out ${CLIENT}.csr \
    -subj "/C=SO/ST=Earth/L=Mountain/O=$O/OU=$OU/CN=127.0.0.1"
  openssl x509  -req -in ${CLIENT}.csr \
    -extfile <(printf "subjectAltName=DNS:127.0.0.1") \
    -CA ca.crt -CAkey ca.key -out ${CLIENT}.crt -days ${DAYS} -sha256 -CAcreateserial
  rm ${CLIENT}.csr
#  cat ${CLIENT}.crt ${CLIENT}.key > ${CLIENT}.pem
}

generate_client client Client Client-OU
#generate_client client.b Client-B Client-B-OU
