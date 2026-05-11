package main

import (
	"flag"
	"log"
)

func main() {
	mode := flag.String("mode", "create", "Operation mode: create, update, delete, or multitenancy")
	clientID := flag.String("clientID", "", "Aruba Cloud client ID")
	clientSecret := flag.String("clientSecret", "", "Aruba Cloud client secret")
	projectID := flag.String("projectID", "", "Project ID (required for update and delete modes)")
	debug := flag.Bool("debug", false, "Enable verbose HTTP logging")
	flag.Parse()

	namePrefix = "sdk-example"
	nameSuffix = generateRandomSuffix()

	switch *mode {
	case "create":
		runCreateExample(*clientID, *clientSecret, *debug)
	case "update":
		if *projectID == "" {
			log.Fatal("update mode requires -projectID")
		}
		runUpdateExample(*clientID, *clientSecret, *projectID, *debug)
	case "delete":
		if *projectID == "" {
			log.Fatal("delete mode requires -projectID")
		}
		runDeleteExample(*clientID, *clientSecret, *projectID, *debug)
	case "multitenancy":
		runMultitenancyExample()
	default:
		log.Fatalf("unknown mode %q; use create, update, delete, or multitenancy", *mode)
	}
}
