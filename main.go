package main

// TODO refactor some of the funcs here to different packages/files

import (
	"os"

	"github.com/amitizle/twitter_lists_manager/cmd"
	"github.com/amitizle/twitter_lists_manager/internal/printer"
	"github.com/amitizle/twitter_lists_manager/internal/twitterclient"
	"gopkg.in/alecthomas/kingpin.v2"
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
		cmd.ApplyLists(client, *inFile)
	case importCommand.FullCommand():
		cmd.ImportLists(client, *importOutputFile)
	case listCommand.FullCommand():
		cmd.ListLists(client)
	}
}

func twitterClient() (*twitterclient.Client, error) {
	c := twitterclient.NewClient()
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
