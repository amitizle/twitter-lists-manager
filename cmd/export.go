package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/amitizle/twitter_lists_manager/internal/printer"
	"github.com/amitizle/twitter_lists_manager/internal/twitterclient"
)

// ExportLists exports the lists owned by the authenticated user to a JSON formatted outputFile
func ExportLists(client *twitterclient.Client, outputFile string) { // TODO rename to export
	ownedLists, err := client.GetOwnedLists(nil)
	if err != nil {
		printer.Fatalf("Error: %v", err)
	}
	for _, list := range ownedLists {
		client.PopulateListMembers(list) // TODO Should this be a member of List?
	}
	jsonBytes, err := json.MarshalIndent(ownedLists, "", "  ") // TODO add pretty out optional
	if err != nil {
		printer.Fatalf("Error while marshaling JSON: %v", err)
	}
	jsonString := string(jsonBytes)
	if outputFile == "" {
		printer.NoColor(jsonString)
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			printer.Fatalf("Error creating file %s: %v", outputFile, err)
		}
		defer file.Close()
		_, err = fmt.Fprintf(file, jsonString)
		if err != nil {
			printer.Fatalf("Error writing to file file %s: %v", outputFile, err)
		}
		printer.Infof("Successfully wrote lists to %s", outputFile)
	}
}
