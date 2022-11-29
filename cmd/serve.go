package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// rootCmd represents the base command when called without any subcommands
var serveCmd = &cobra.Command{
	Use:    "serve -file <file>",
	PreRun: toggleDebug,
	Short:  "Serve a JWKS file for local use",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: runServe,
}

var (
	servePort int    = 9000
	serveFile string = "./out/jwks.pub"
)

// initServe sets up flags for the serve command
func initServe() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 9000, "port to serve jwks file on")
	serveCmd.Flags().StringVarP(&serveFile, "file", "f", "", "File path to serve as /.well-known/jwks.json")
}

// runServe serves the provided file
func runServe(cmd *cobra.Command, args []string) {
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Errorf("File %v not found", file)
		os.Exit(1)
	}

	if file == "" {
		log.Errorf("File must be provided")
		os.Exit(1)
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		log.Errorf("Port %d not available", port)
		os.Exit(1)
	}

	if port <= 0 || port > 65535 {
		log.Errorf("Port must be in range 1-65535")
		os.Exit(1)
	}

	log.Fatal(serveJWKS(file, port))
}

func serveJWKS(filePath string, port int) error {
	http.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.ServeFile(w, r, filePath)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
