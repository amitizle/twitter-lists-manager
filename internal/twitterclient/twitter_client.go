package twitterclient

import (
	"errors"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type Client struct {
	api               *anaconda.TwitterApi
	userId            int64
	Username          string
	AccessToken       string
	AccessTokenSecret string
	ConsumerKey       string
	ConsumerSecret    string
}

type List struct {
	Slug        string
	Name        string
	Description string
	Members     []string
	Mode        string
	listId      int64
	memberCount int64
}

// NewClient returns a new twitterclient
func NewClient() *Client {
	return &Client{
		api: &anaconda.TwitterApi{},
	}
}

// Init initializes the client
func (client *Client) Init() error {
	if client.Username == "" {
		return errors.New("username cannot be empty")
	}
	newAPI := anaconda.NewTwitterApiWithCredentials(client.AccessToken, client.AccessTokenSecret, client.ConsumerKey, client.ConsumerSecret)
	client.api = newAPI
	user, err := client.api.GetUsersShow(client.Username, nil)
	if err != nil {
		return err
	}
	client.userId = user.Id
	return nil
}

// GetOwnedLists returns the lists owned by the authenticated user
func (client *Client) GetOwnedLists(v url.Values) ([]*List, error) {
	ownedLists, err := client.api.GetListsOwnedBy(client.userId, nil)
	if err != nil {
		return []*List{}, err
	}
	ownedListsRet := make([]*List, len(ownedLists))
	for i, list := range ownedLists {
		ownedListsRet[i] = &List{
			Slug:        list.Slug,
			Name:        list.Name,
			Description: list.Description,
			Members:     []string{},
			Mode:        list.Mode,
			memberCount: list.MemberCount,
			listId:      list.Id,
		}
	}
	return ownedListsRet, nil
}

// PopulateListMembers populates the Members field in the List struct
func (client *Client) PopulateListMembers(list *List) error {
	v := url.Values{}
	var cursor anaconda.UserCursor
	var err error
	for cursor.Next_cursor_str != "0" {
		cursor, err = client.api.GetListMembers(client.Username, list.listId, v)
		if err != nil {
			return err
		}
		for _, user := range cursor.Users {
			list.Members = append(list.Members, user.ScreenName)
		}
		v.Set("cursor", cursor.Next_cursor_str)
	}
	return nil
}

// RemoveUsersFromList is not yet supported because anaconda does not support
// POST /lists/destroy yet (https://developer.twitter.com/en/docs/accounts-and-users/create-manage-lists/api-reference/post-lists-members-destroy)
func (client *Client) RemoveUsersFromList(list *List, usersToRemove []string) error {
	return nil
}

// AddUsersToList adds the given users to the given list
func (client *Client) AddUsersToList(list *List, usersToAdd []string) error {
	_, err := client.api.AddMultipleUsersToList(usersToAdd, list.listId, nil)
	if err != nil {
		return err
	}
	return nil
}
