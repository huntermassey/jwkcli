# JWK CLI

A CLI for interacting with JSON Web Keys (JWK) and JSON Web Key Sets (JWKS).

! WARNING !

THIS IS HEAVILY WORK IN PROGRESS, AND CURRENTLY INTENDED MOSTLY FOR SELF-EDUCATION

## Target

Implement a CLI for generating, inspecting, and serving JWK/JWKS files with the following intents:
* For use in integration testing of services verifying JWT tokens against JWKS
* For understanding the use of the various JWKS fields such as x5t, x5c, etc
* For being able to rapidly inspect the contents of a remote JWKS URL

## Usage

Sketch of desired usage below

### Generate a new JWK

jwkcli generate -type enc -alg=RSA256 

jwkcli generate -type sig -alg=...

### Inspect a remote JWKS json

jwkcli inspect -remote https://example.com/.well-known/jwks.json

### Serve JWKS json

jwkcli serve -file <file>

## Notes

### Inspecting x5c

Save x5c value base64 decoded to file, and inspect with:

```
openssl x509 -inform der -in test.der -text -noout
```