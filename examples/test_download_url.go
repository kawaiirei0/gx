package main

import (
	"fmt"

	"github.com/yourusername/gx/internal/downloader"
	"github.com/yourusername/gx/internal/platform"
)

func main() {
	platformAdapter := platform.NewAdapter()
	dl := downloader.NewDownloader()

	// Try actual available versions
	versions := []string{"1.25.4", "1.24.10"}

	fmt.Println("Testing download URL generation for recent Go versions:\n")

	for _, version := range versions {
		url, err := dl.GetDownloadURL(version, platformAdapter.GetOS(), platformAdapter.GetArch())
		if err != nil {
			fmt.Printf("❌ Go %s: %v\n", version, err)
		} else {
			fmt.Printf("✓ Go %s: %s\n", version, url)
		}
	}
}
