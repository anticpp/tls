## Purpose

Demonstrate SSL/TLS.

## Requirement

[Openssl](https://www.openssl.org/) is required to manage certification.

## CA

First of all, you need a CA(Certificate Authority).
If you already have a one, skip this section.

Generate CA key and certification:

```
## Generate rsa key
openssl genrsa -out cakey.pem

## Generate x509 certification
openssl req -new -x509 -key cakey.pem -out cacert.pem
```

> Note: Set `Common Name` to `example.com` when generating certification.

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

> Note: Set `Common Name` to `example.com` when generating CSR.

Sign CSR with CA:

```
## Sign with CA
## Default using key and certificate in directory `$ca_home`, etc. '/etc/pki/CA/'
openssl ca -in svrcsr.pem -out svrcert.pem

## Or sign with specify key and certificate, using `-keyfile` and `-cert`
openssl ca -keyfile cakey.pem -cert cacert.pem -in svrcsr.pem -out svrcert.pem
```

## Add host

Add host line to your `/etc/hosts`, as below.

```
127.0.0.1 example.com
```

## Test

Build and run `tls`.

Run server with `./tls -l`.
Run client with `./tls -addr=example.com:20012` 
    or using specify `CA certification` `./tls -addr=example.com:20012 -ca=./cacert.pem`

## Show certificate information

```
openssl x509 -in svrcert.pem -noout -text
```

## Errors

Error:

```
failed to update database
TXT_DB error number 2
```

Problem: 
If you generate a signed certificate with the same CN (Common Name) information that the CA certificate that you've generated before.

Solution:
Use `openssl ca -revoke <path-to-certfile>` to revoke the exsisting certificate. You can find it at `$ca_home/newcerts/`.
