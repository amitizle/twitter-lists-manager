package main

// TODO refactor some of the funcs here to different packages/files

import (
	"encoding/json"
	"fmt"
	"github.com/amitizle/twitter_lists_manager/internal/printer"
	"github.com/amitizle/twitter_lists_manager/internal/twitter_client"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	app               = kingpin.New("twitter", "Manage Twitter lists").DefaultEnvars().Author("Amit Goldberg").Version("0.1.0")
	user              = app.Flag("user", "Twitter username (without the '@')").Required().Short('u').String()
	accessToken       = app.Flag("access-token", "Twitter API access token").Required().Short('a').String()
	accessTokenSecret = app.Flag("access-token-secret", "Twitter API access token secret").Required().Short('s').String()
	consumerKey       = app.Flag("consumer-key", "Twitter consumer API key").Required().Short('c').String()
	consumerSecret    = app.Flag("consumer-secret", "Twitter consumer API key secret").Required().Short('x').String()

	applyCommand  = app.Command("apply", "Apply changes to the lists")
	listCommand   = app.Command("list", "List user lists")
	importCommand = app.Command("import", "Import lists from Twitter")

	importOutputFile = importCommand.Flag("out-file", "Output file for the import (JSON formatted)").Short('o').String()
	inFile           = applyCommand.Flag("file", "Input JSON formatted file").Required().Short('f').String()
)

func main() {
	kingpin.CommandLine.Help = "A CLI to manage (update/create) Twitter lists"
	kingpinParse := kingpin.MustParse(app.Parse(os.Args[1:]))
	client, err := twitterClient()
	if err != nil {
		printer.Fatalf("Error getting Twitter client: %v", err)
	}
	switch kingpinParse {
	case applyCommand.FullCommand():
		applyLists(client, *inFile)
	case importCommand.FullCommand():
		importLists(client, *importOutputFile)
	case listCommand.FullCommand():
		listLists(client)
	}
}

func readJsonLists(jsonFilePath string) ([]twitter_client.List, error) {
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

func twitterClient() (*twitter_client.Client, error) {
	c := twitter_client.NewClient()
	c.Username = *user
	c.AccessToken = *accessToken
	c.AccessTokenSecret = *accessTokenSecret
	c.ConsumerKey = *consumerKey
	c.ConsumerSecret = *consumerSecret
	if err := c.Init(); err != nil {
		return c, err
	}
	return c, nil
}

func applyLists(client *twitter_client.Client, inFile string) {
	lists, err := readJsonLists(inFile)
	if err != nil {
		printer.Fatalf("Error while reading lists from JSON file %s: %v", inFile, err)
	}
	printer.Infof("%#v", lists)
}

func importLists(client *twitter_client.Client, outputFile string) {
	ownedLists, err := client.GetOwnedLists(nil)
	if err != nil {
		printer.Fatalf("Error: %v", err)
	}
	for _, list := range ownedLists {
		client.PopulateListMembers(list)
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

func listLists(client *twitter_client.Client) {
	ownedLists, err := client.GetOwnedLists(nil)
	if err != nil {
		printer.Fatalf("Error: %v", err)
	}
	for _, l := range ownedLists {
		printer.NoColorf("Name: %s, Slug: %s, Description: %s", l.Name, l.Slug, l.Description)
	}
}
