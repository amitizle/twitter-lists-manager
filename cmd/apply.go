package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/amitizle/twitter_lists_manager/internal/printer"
	"github.com/amitizle/twitter_lists_manager/internal/twitter_client"
)

// ApplyLists applying the JSON formatted inFile, modifying, adding and removing lists and
// members from the lists
func ApplyLists(client *twitter_client.Client, inFile string) {
	localLists, err := readJSONLists(inFile)
	if err != nil {
		printer.Fatalf("Error while reading lists from JSON file %s: %v", inFile, err)
	}
	printer.Infof("%#v", localLists)
	ownedLists, err := client.GetOwnedLists(nil)
	if err != nil {
		printer.Fatalf("Error: %v", err)
	}
	printer.Infof("%#v", ownedLists)
	localListsMap := make(map[string]*twitter_client.List)
	ownedListsMap := make(map[string]*twitter_client.List)
	for _, list := range localLists {
		localListsMap[list.Slug] = &list // TODO something?
	}
	for _, ownedList := range ownedLists {
		ownedListsMap[ownedList.Slug] = ownedList
	}
	for slug, localList := range localListsMap {
		if ownedList, ok := ownedListsMap[slug]; ok {
			client.PopulateListMembers(ownedList)
			screenNamesToRemove := sliceDifference(ownedList.Members, localList.Members)
			screenNamesToAdd := sliceDifference(localList.Members, ownedList.Members)
			printer.Infof(strings.Join(screenNamesToAdd, "\n"))
			printer.Redf(strings.Join(screenNamesToRemove, "\n"))
			// diff members
			// update list https://developer.twitter.com/en/docs/accounts-and-users/create-manage-lists/api-reference/post-lists-update
		} else { // new list TODO support deleting lists
		}
	}
}

func sliceDifference(a, b []string) []string {

	diff := make([]string, 0)
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

func readJSONLists(jsonFilePath string) ([]twitter_client.List, error) {
	fullPath, err := filepath.Abs(jsonFilePath)
	if err != nil {
		return []twitter_client.List{}, err
	}
	jsonFile, err := os.Open(fullPath)
	if err != nil {
		return []twitter_client.List{}, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []twitter_client.List{}, err
	}

	lists := make([]twitter_client.List, 0)
	err = json.Unmarshal(byteValue, &lists)
	if err != nil {
		return []twitter_client.List{}, err
	}

	return lists, nil
}
