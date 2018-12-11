## Purpose

Demonstrate SSL/TLS.

## Requirement

[Openssl](https://www.openssl.org/) is required to manage certification.

## CA

Firstly, you need a CA(Certificate Authority).
If you already have a one, skip this section.

Generate CA key and certification:

```
## Generate rsa key
openssl genrsa -out cakey.pem

## Generate x509 certification
openssl req -new -x509 -key cakey.pem -out cacert.pem
```

Copy to CA home directory:

```
## Default CA home directory is '/etc/pki/CA/'
## We use $ca_home
cp cakey.pem $ca_home/private/
cp cacert.pem $ca_home/
```

## Server certification

Generate server key and CSR(Certificate Signature Request):

```
## Generate rsa key
openssl genrsa -out svrkey.pem

## Generate CSR(Certificate Signature Request)
openssl req -new -key svrkey.pem -out svrcsr.pem
```

Sign CSR with CA:

```
## Sign with CA
openssl ca -in svrcsr.pem -out svrcert.pem
```

## Test

Build and run `tls`.

Run server with `./tls -l`.
Run client with `./tls`.

## Show certificate information

```
openssl x509 -in svrcert.pem -noout -text
```
