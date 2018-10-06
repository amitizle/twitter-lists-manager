package cmd

import (
	"encoding/json"
	"fmt"
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
	ownedLists, err := client.GetOwnedLists(nil)
	if err != nil {
		printer.Fatalf("Error: %v", err)
	}
	localListsMap := make(map[string]twitterclient.List)
	ownedListsMap := make(map[string]*twitterclient.List)
	for _, list := range localLists {
		localListsMap[list.Slug] = list // TODO something?
	}
	for _, ownedList := range ownedLists {
		ownedListsMap[ownedList.Slug] = ownedList
	}
	for slug, localList := range localListsMap {
		if ownedList, ok := ownedListsMap[slug]; ok {
			printer.Infof("List %s already exists, diffing members and adding accordingly", ownedList.Name)
			ownedList.Name = localList.Name
			ownedList.Description = localList.Description
			client.PopulateListMembers(ownedList)
			screenNamesToRemove := sliceDifference(ownedList.Members, localList.Members)
			screenNamesToAdd := sliceDifference(localList.Members, ownedList.Members)
			printer.Infof("[Adding]\n%s", strings.Join(screenNamesToAdd, "\n"))
			printer.Redf("[Removing]\n%s", strings.Join(screenNamesToRemove, "\n"))
			err := client.AddUsersToList(ownedList, screenNamesToAdd)
			if err != nil {
				fmt.Errorf("Error while adding users to list %s: %v", ownedList.Slug, err)
			}
			err = client.RemoveUsersFromList(ownedList, screenNamesToRemove)
			if err != nil {
				fmt.Errorf("Error while removing users from list %s: %v", ownedList.Slug, err)
			}
			if err = client.UpdateList(ownedList); err != nil {
				printer.Redf("Error updating a list %s: %v", localList.Name, err)
			}
		} else { // new list TODO support deleting lists
			printer.Infof("List %s does not exists, creating", localList.Name)
			if err = client.CreateList(&localList); err != nil {
				printer.Redf("Error creating a list %s: %v", localList.Name, err)
			}
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
