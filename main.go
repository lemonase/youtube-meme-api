package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lemonase/youtube-meme-api/client"
	"github.com/lemonase/youtube-meme-api/server"
)

var (
	port       = flag.String("port", "8000", "Port to listen on (default is 8000)")
	apiKey     = flag.String("key", "", "API key to access Google resources")
	secretFile = flag.String("secretFile", "", "Credentials file downloaded from GCP (/path/to/credentials.json)")
)

func handleArgs() {

	// flag parsing
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
	}
	flag.Parse()

	// client/api parameters
	if *apiKey != "" {
		client.InitClientsWithAPIKey(*apiKey)
	} else if os.Getenv("YT_API_KEY") != "" {
		client.InitClientsWithAPIKey(os.Getenv("YT_API_KEY"))
	} else if *secretFile != "" {
		client.InitClientsWithSecretJSONFile(*secretFile)
	} else {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "Please specify either an API Key or a credentials.json file\n")
		os.Exit(1)
	}

	// server parameters
	if os.Getenv("PORT") != "" {
		*port = os.Getenv("PORT")
	}

	// prepend ":" for port if not already
	if strings.Index(*port, ":") == -1 {
		*port = ":" + *port
	}

}

func main() {
	handleArgs()
	server.FetchInitResources()
	server.InitServer(*port)
}
