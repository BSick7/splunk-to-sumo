package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BSick7/splunk-to-sumo/agent"
)

var Version string

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(0)
	}

	bind := os.Args[1]
	host := os.Args[2]

	config := agent.HttpConfig{
		Address:      fmt.Sprintf(bind),
		ReadTimeout:  readDurationEnvVar("STS_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: readDurationEnvVar("STS_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  readDurationEnvVar("STS_IDLE_TIMEOUT", 120*time.Second),
	}

	server := agent.NewServer(host)
	fmt.Printf("Http Server Config: %+v\n", config)
	fmt.Printf("Starting splunk-to-sumo v%s on %s with hostname=%s...\n", Version, bind, host)
	if err := server.ListenAndServe(config); err != nil {
		log.Fatalf(err.Error())
	}
}

func usage() {
	fmt.Println(`usage: splunk-to-sumo [bind-address] [sumo-host]

  bind-address      This specifies which ip:port to bind the server.

  sumo-host         This will forward messages to sumologic with
                    this set as the X-Sumo-Host header.

Environment Variables:

  STS_READ_TIMEOUT  This specifies the ReadTimeout on the listening http server.
                    This must be parsable by golang time.ParseDuration.
                    default=5s

  STS_WRITE_TIMEOUT This specifies the WriteTimeout on the listening http server.
                    This must be parsable by golang time.ParseDuration.
                    default=10s

  STS_IDLE_TIMEOUT  This specifies the IdleTimeout on the listening http server.
                    This must be parsable by golang time.ParseDuration.
                    default=120s
`)
}

func readDurationEnvVar(name string, fallback time.Duration) time.Duration {
	val := os.Getenv(name)
	if val == "" {
		return fallback
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("[WARN] could not parse duration %s=%s: %s\n", name, val, err)
		return fallback
	}

	return dur
}
