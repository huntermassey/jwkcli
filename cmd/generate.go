package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"

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
	alg string = ""
)

func initGenerate() {
	serveCmd.Flags().StringVarP(&serveFile, "alg", "f", "", "File path to serve as /.well-known/jwks.json")
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

	jwkPubBytes, err := json.Marshal(pubKey)
	fmt.Println(string(jwkPubBytes))
}
