package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

func (c *Client) Updates(offset, limit int) ([]Updates, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := c.makeRequests("getUpdates", query)

	if err != nil {
		return nil, fmt.Errorf("can't get updates")
	}

	var response UpdatesResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("can't get updates")
	}

	return response.Result, nil
}

func (c *Client) SendMessage(chatId int, text string, markup ReplyKeyboardMarkup) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatId))
	query.Add("text", text)
	markupJSON, err := json.Marshal(markup)
	if err != nil {
		return fmt.Errorf("can't marshal markup: %w", err)
	}
	query.Add("keyboard", string(markupJSON))
	log.Println(query)
	_, err = c.makeRequests("sendMessage", query)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (c *Client) makeRequests(method string, query url.Values) ([]byte, error) {
	url := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't make the request: %w", err)
	}

	request.URL.RawQuery = query.Encode()

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't make the request: %w", err)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("can't make the request: %w", err)
	}

	return body, nil
}
