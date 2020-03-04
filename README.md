## Purpose

Demonstrate about ssl/tls with simple example.

## Keywords

- Openssl - A tool to manage ssl/tls stuff.
- CA - Certificate Authority.
- x509 - A standard format of public key certificates. Which suffix is `.pem`.
- PKI - Public Key Infrastructure.

- Certificate - Certificate.
- CSR - Certificate Signing Request.

- mTLS - Mutual TLS. Both side verify certificate of the other.

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
## Set common name to 'Root CA'.
openssl req -new -x509 -key cakey.pem -out cacert.pem
```

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
## Set common name to '*.example.com', which is a wildcard name.
openssl req -new -key svrkey.pem -out svrcsr.pem

```

Sign CSR with CA:

```
## Sign with CA
openssl ca -in svrcsr.pem -out svrcert.pem
```

## Simple test

Run server:

```
./tls -l
```

Run client:

```
./tls 
``` 

## mTLS test

With mTLS mode, you should create client key and certificate as the server does. For simplify, you can copy the server key and certificate.

```
cp svrkey.pem cltkey.pem
cp svrcert.pem cltcert.pem
```

Run server:

```
./tls -l -mtls
```

Run client:

```
./tls -mtls
```

## Show certificate

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
