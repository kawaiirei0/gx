package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/gx/pkg/constants"
	"github.com/yourusername/gx/pkg/interfaces"
)

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	fmt.Println("Fetching available Go versions from official API...\n")

	resp, err := client.Get(constants.GoVersionsAPIURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var versions []interfaces.RemoteVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return
	}

	fmt.Printf("Found %d versions\n\n", len(versions))
	fmt.Println("Latest 10 stable versions:")

	count := 0
	for _, v := range versions {
		if v.Stable && count < 10 {
			fmt.Printf("  - %s\n", v.Version)
			
			// Show available platforms for first version
			if count == 0 {
				fmt.Println("    Available platforms:")
				platformMap := make(map[string]bool)
				for _, f := range v.Files {
					key := fmt.Sprintf("%s/%s", f.OS, f.Arch)
					if !platformMap[key] {
						platformMap[key] = true
						fmt.Printf("      - %s\n", key)
					}
				}
			}
			count++
		}
	}
}
