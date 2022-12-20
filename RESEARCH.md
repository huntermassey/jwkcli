# Collection of links obtained during research of potential use

I like to try and collect all the links I run across as I work on and research a project and prior work in a space. Below is that collection!

## Cool implementations

* https://github.com/square/go-jose
    * Interesting cli tool in-tree: https://github.com/square/go-jose/tree/v2/jwk-keygen
* A rather complete JOSE library: https://github.com/lestrrat-go/jwx
    * Great beginner friendly how-to! https://github.com/lestrrat-go/jwx/tree/develop/v2/docs
* Existing JOSE cli, fills some of my use case but I want to build my own and target my specific use-case https://github.com/antoniomo/jose-tool
    * Fills one desired use-case of JWK generation
* A JWT library working with certificates, which is just what I was looking for an example of https://github.com/crossedbot/simplejwt
* blog post about what x5c and a certificate is used for in a JWT https://software-factotum.medium.com/validating-rsa-signature-for-a-jws-more-about-jwk-and-certificates-e8a3932669f1
* Github gist for reading x5c https://gist.github.com/monmohan/d08d41c856a54d7e7619f8fba8afdf44
    * This function in particular of interest to me, regarding key.E property decoding:

```
//GetPublicKeyFromModulusAndExponent returns a rsa.PublicKey based on n and e members in JWK JSON
func GetPublicKeyFromModulusAndExponent(jwk *JWK) *rsa.PublicKey {
	n, _ := base64.RawURLEncoding.DecodeString(jwk.N)
	e, _ := base64.RawURLEncoding.DecodeString(jwk.E)
	z := new(big.Int)
	z.SetBytes(n)
	//decoding key.E returns a three byte slice, https://golang.org/pkg/encoding/binary/#Read and other conversions fail
	//since they are expecting to read as many bytes as the size of int being returned (4 bytes for uint32 for example)
	var buffer bytes.Buffer
	buffer.WriteByte(0)
	buffer.Write(e)
	exponent := binary.BigEndian.Uint32(buffer.Bytes())
	return &rsa.PublicKey{N: z, E: int(exponent)}
}
```

# Various

* https://www.fastly.com/blog/key-size-for-tls
    * Trying to determine modern good default for RSA key length (2048 being assumed). Turns out, 2048 good until 2030 according to NIST
    * Also a reminder of why elliptic curve (ECDSA) is the "new" hotness - way less bits required to meet "bits of security" requirement

## Stackoverflow

* A simple review with well known JWKS json https://stackoverflow.com/a/61200654
* JWK, vs JWE, vs JWS, vs JW... https://stackoverflow.com/a/74257561 
    * highly useful for starting out as a beginner
* On the x5c, x5t parameters https://stackoverflow.com/a/69232318

## Shoutouts

* Auth0 for their wonderful docs and educational material
    * including https://auth0.com/blog/json-web-token-signing-algorithms-overview/#Common-JWT-Signing-Algorithms
