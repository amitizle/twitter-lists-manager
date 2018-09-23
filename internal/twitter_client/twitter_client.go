package twitter_client

import (
	"errors"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
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

func NewClient() *Client {
	return &Client{
		api: &anaconda.TwitterApi{},
	}
}

func (client *Client) Init() error {
	if client.Username == "" {
		return errors.New("username cannot be empty")
	}
	newApi := anaconda.NewTwitterApiWithCredentials(client.AccessToken, client.AccessTokenSecret, client.ConsumerKey, client.ConsumerSecret)
	client.api = newApi
	user, err := client.api.GetUsersShow(client.Username, nil)
	if err != nil {
		return err
	}
	client.userId = user.Id
	return nil
}

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
