package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/deichindianer/semver-go/pkg/sem"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Usage: semver <kind>\nKind can be one of major,minor,patch\n")
		os.Exit(1)
	}
	kind, err := sem.ParseKind(os.Args[1])
	if err != nil {
		fmt.Printf("Usage: semver <kind>\nKind can be one of major,minor,patch\n")
		os.Exit(1)
	}

	curDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %s\n", err)
		os.Exit(1)
	}
	latestVersion, err := sem.GetLatestVersion(curDir)
	if err != nil {
		if errors.Is(err, sem.ErrNoVersionsAvailable) {
			// This is an opinionated initial version
			fmt.Println("v0.1.0")
			os.Exit(0)
		} else {
			fmt.Printf("Failed to get latest version: %s\n", err)
			os.Exit(1)
		}
	}
	if err := latestVersion.Next(kind); err != nil {
		fmt.Printf("Failed to get next version: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(latestVersion)
}
