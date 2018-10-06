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
	listID      int64
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
	ownedListsRet := make([]*List, 0)
	for _, list := range ownedLists {
		ownedListsRet = append(ownedListsRet, &List{
			Slug:        list.Slug,
			Name:        list.Name,
			Description: list.Description,
			Members:     []string{},
			Mode:        list.Mode,
			memberCount: list.MemberCount,
			listID:      list.Id,
		})
	}
	return ownedListsRet, nil
}

// PopulateListMembers populates the Members field in the List struct
func (client *Client) PopulateListMembers(list *List) error {
	v := url.Values{}
	var cursor anaconda.UserCursor
	var err error
	for cursor.Next_cursor_str != "0" {
		cursor, err = client.api.GetListMembers(client.Username, list.listID, v)
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
// There is an open PR: https://github.com/ChimeraCoder/anaconda/pull/247
func (client *Client) RemoveUsersFromList(list *List, usersToRemove []string) error {
	return nil
}

// UpdateList updates the given list, however it's not implemented yet
// https://developer.twitter.com/en/docs/accounts-and-users/create-manage-lists/api-reference/post-lists-update
func (client *Client) UpdateList(list *List) error {
	return nil
}

// AddUsersToList adds the given users to the given list
// TODO support more than 100 members (should be separate API calls)
// AddMultipleUsersToList is a wrapper to create_all.json, which is limited to
// 100 users per call;
// https://developer.twitter.com/en/docs/accounts-and-users/create-manage-lists/api-reference/post-lists-members-create_all
func (client *Client) AddUsersToList(list *List, usersToAdd []string) error {
	_, err := client.api.AddMultipleUsersToList(usersToAdd, list.listID, nil)
	return err
}

// CreateList creates the list using the *List given
func (client *Client) CreateList(list *List) error {
	v := url.Values{}
	v.Add("mode", list.Mode)
	newList, err := client.api.CreateList(list.Name, list.Description, v)
	if err != nil {
		return err
	}
	list.listID = newList.Id
	return client.AddUsersToList(list, list.Members)
}
