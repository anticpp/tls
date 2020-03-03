## Purpose

Simple demonstration on ssl/tls for better understanding.

## Keywords

Openssl - A tool to manage ssl/tls stuff.
CA - Certificate Authority.
x509 - A standard format of public key certificates. Which suffix is `.pem`.
PKI - Public Key Infrastructure.

Certificate - Certificate.
CSR - Certificate Signing Request.

## CA

Let's build a CA(Certificate Authority). If you already have a one, skip this section.

You should have installed [Openssl](https://www.openssl.org/) on your host, then you can find the pki(Public Key Infrastructure) directory at `/etc/pki/`. The CA files will be located at `/etc/pki/CA/`.

Create CA `index.txt` and `serial`:

```
cd /etc/pki/CA/
touch index.txt serial
echo 01 > serial
```

Generate CA key and certification:

```
## Generate rsa key
openssl genrsa -out cakey.pem

## Generate x509 certification
openssl req -new -x509 -key cakey.pem -out cacert.pem
```

> Note: Set `Common Name` to something like `Root CA`.

Copy to CA home directory:

```
cp cakey.pem /etc/pki/CA/private/
cp cacert.pem /etc/pki/CA/
```

## Server certification

Generate server key and CSR(Certificate Signature Request):

```
## Generate rsa key
openssl genrsa -out svrkey.pem

## Generate CSR(Certificate Signature Request)
openssl req -new -key svrkey.pem -out svrcsr.pem

```

> Note: Set `Common Name` to `*.example.com` when generating CSR, which is a wildcard name.

Sign CSR with CA:

```
## Sign with CA
## Default using key and certificate in directory `$ca_home`, etc. '/etc/pki/CA/'
openssl ca -in svrcsr.pem -out svrcert.pem

## Or sign with specify key and certificate, using `-keyfile` and `-cert`
openssl ca -keyfile cakey.pem -cert cacert.pem -in svrcsr.pem -out svrcert.pem
```

## Add host

Add line to your `/etc/hosts`.

```
127.0.0.1 www.example.com
```

## Test

Server side:

```
./tls -l
```

Client side:

```
./tls -addr=www.example.com:20012
``` 

## How to show certificate information

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
