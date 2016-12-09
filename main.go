package main

import (
	"fmt"
	"github.com/BSick7/splunk-to-sumo/agent"
	"net/http"
	"os"
)

var Version string

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(0)
	}

	bind := os.Args[1]
	host := os.Args[2]

	server := agent.NewServer(host)
	fmt.Printf("Starting splunk-to-sumo v%s on %s with hostname=%s...\n", Version, bind, host)
	http.ListenAndServe(fmt.Sprintf(bind), server)
}

func usage() {
	fmt.Println(`usage: splunk-to-sumo [bind-address] [sumo-host]

  bind-address      This specifies which ip:port to bind the server.

  sumo-host         This will forward messages to sumologic with
                    this set as the X-Sumo-Host header.
`)
}
