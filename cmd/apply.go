package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/amitizle/twitter_lists_manager/internal/printer"
	"github.com/amitizle/twitter_lists_manager/internal/twitterclient"
)

// ApplyLists applying the JSON formatted inFile, modifying, adding and removing lists and
// members from the lists
func ApplyLists(client *twitterclient.Client, inFile string) {
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
	localListsMap := make(map[string]*twitterclient.List)
	ownedListsMap := make(map[string]*twitterclient.List)
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
			client.AddUsersToList(ownedList, screenNamesToAdd)
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

func readJSONLists(jsonFilePath string) ([]twitterclient.List, error) {
	fullPath, err := filepath.Abs(jsonFilePath)
	if err != nil {
		return []twitterclient.List{}, err
	}
	jsonFile, err := os.Open(fullPath)
	if err != nil {
		return []twitterclient.List{}, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []twitterclient.List{}, err
	}

	lists := make([]twitterclient.List, 0)
	err = json.Unmarshal(byteValue, &lists)
	if err != nil {
		return []twitterclient.List{}, err
	}

	return lists, nil
}
