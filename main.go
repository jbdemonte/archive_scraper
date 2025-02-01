package main

import (
	"archive-dl/scraper"
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Directory where files will be downloaded
const defaultDownloadDir = "downloads"

func main() {
	var itemID string
	var destination string
	var noChecksum bool

	for _, arg := range os.Args[1:] {
		if arg == "--no-checksum" {
			noChecksum = true
		} else if itemID == "" {
			itemID = getItemID(arg)
		} else if destination == "" {
			destination = arg
		}

	}

	if itemID == "" {
		log.Fatal("Usage: ./archive_scraper <URL or TokenID> [<DESTINATION>]")
	}

	if destination == "" {
		destination = defaultDownloadDir
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	err := scraper.DownloadAllFilesFromItemID(ctx, itemID, destination, noChecksum)
	if err != nil {
		log.Fatal(err)
	}
}

// extract item ID from URL or a token
// example :
//
//	https://archive.org/details/karous-eng => karous-eng
//	karous-eng
func getItemID(p string) string {
	if strings.Contains(p, "https://") == false {
		return p
	}
	parsedURL, err := url.Parse(p)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	return pathSegments[1]
}
