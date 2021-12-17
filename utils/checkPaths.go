package utils

import (
	"log"
	"os"
)

func CheckPaths() {
	log.Println("Checking paths...")
	// if exist
	_, err := os.Stat(os.Getenv("STORAGE_DIR"))
	if err != nil {
		err := os.MkdirAll(os.Getenv("STORAGE_DIR"), 0777) // create directories if don't exist
		if err != nil {
			log.Fatal(err)
		}
	}
}
