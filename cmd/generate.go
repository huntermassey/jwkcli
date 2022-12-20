package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/lestrrat-go/jwx/v2/cert"
	"github.com/lestrrat-go/jwx/v2/jwk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var generateCmd = &cobra.Command{
	Use:    "generate",
	PreRun: toggleDebug,
	Short:  "Generate a new JWK/JWKS file",
	Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: runGenerate,
}

var (
	alg     string = ""
	useCert bool   = false
)

func initGenerate() {
	serveCmd.Flags().StringVarP(&serveFile, "alg", "f", "", "File path to serve as /.well-known/jwks.json")
	serveCmd.Flags().BoolVarP(&useCert, "use-cert", "c", false, "Whether to generate and include a self-signed cert for x5c, x5t values")
}

func runGenerate(cmd *cobra.Command, args []string) {
	// target initial implementation: JWS asymmetric RSASSA-PKCS1-v1_5 + SHA256

	// default 2048 length
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Errorf("Failure to generate RSA key, error was %v", err)
	}

	// This "FromRaw" is nice
	key, err := jwk.FromRaw(privateKey)
	if _, ok := key.(jwk.RSAPrivateKey); !ok {
		log.Errorf("expected jwk.SymmetricKey, got %T\n", key)
		return
	}

	rsaJWK := key.(jwk.RSAPrivateKey)

	// output private key
	// jwkBytes, err := json.Marshal(rsaJWK)
	// fmt.Println(string(jwkBytes))

	// output public key
	pubKey, err := rsaJWK.PublicKey()
	if err != nil {
		log.Errorf("error producing public key from private, %v", err)
	}

	_, x5c, x5t := generateCertificateRSA(privateKey)
	chain := cert.Chain{}
	chain.Add([]byte(x5c[0]))
	pubKey.Set("x5c", &chain)
	pubKey.Set("x5t", x5t)

	jwkPubBytes, err := json.Marshal(pubKey)
	fmt.Println(string(jwkPubBytes))
}

func generateCertificateRSA(privKey *rsa.PrivateKey) (cert *x509.Certificate, x5c []string, x5t string) {
	var certTemplate = x509.Certificate{
		SerialNumber: big.NewInt(1337),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"Company Co."},
			CommonName:   "Self-signed test",
		},
		NotBefore:   time.Now().Add(-10 * time.Second),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // 1 year! ... or, close enough
		KeyUsage:    x509.KeyUsageCRLSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:        false,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privKey.PublicKey, privKey)
	if err != nil {
		panic("Failed to create certificate:" + err.Error())
	}

	cert, err = x509.ParseCertificate(certBytes)
	if err != nil {
		panic("Failed to parse created certificate:" + err.Error())
	}

	x5c = []string{base64.StdEncoding.EncodeToString(certBytes)}

	sha1Thumbprint := sha1.Sum(certBytes)
	x5t = base64.URLEncoding.EncodeToString(sha1Thumbprint[:])

	return
}
