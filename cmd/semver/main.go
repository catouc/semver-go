package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/catouc/semver-go/pkg/sem"
)

var (
	ignoreNonSemVerTags = flag.Bool("i", false, "Ignore all tags that are not SemVer compliant instead of failing")
)

func main() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: semver <kind>\nKind can be one of major,minor,patch,current\n")
		flag.PrintDefaults()
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	curDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %s\n", err)
		os.Exit(1)
	}

	latestVersion, err := sem.GetLatestVersion(curDir, *ignoreNonSemVerTags)
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

	k := flag.Args()[0]

	if k == "current" {
		fmt.Println(latestVersion)
		os.Exit(0)
	}

	kind, err := sem.ParseKind(flag.Args()[0])
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	if err := latestVersion.Next(kind); err != nil {
		fmt.Printf("Failed to get next version: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(latestVersion)
}
