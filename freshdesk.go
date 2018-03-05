package freshdesk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type (
	source   int
	status   int
	priority int
)

const (
	// Email ...
	Email source = 1 + iota
	// Portal ...
	Portal
	// Phone ...
	Phone
	// Forum ...
	Forum
	// Twitter ...
	Twitter
	// Facebook ...
	Facebook
	// Chat ...
	Chat
)

const (
	// Open ...
	Open status = 2 + iota
	// Pending ...
	Pending
	// Resolved ...
	Resolved
	// Closed ...
	Closed
)

const (
	// Low ...
	Low priority = 1 + iota
	// Medium ...
	Medium
	// High ...
	High
	// Urgent ...
	Urgent
)

// Client is a freshdesk client that allows access to the freshdesk
// API. It is save to use the client from different goroutines.
type Client struct {
	Domain string
	API    string

	httpClient *http.Client
}

// NewClient returns a new freshdesk client.
func NewClient(domain, api string) (*Client, error) {
	return &Client{Domain: domain, API: api, httpClient: http.DefaultClient}, nil
}

// CreateTicket creates a new ticket.
func (c *Client) CreateTicket(ticket *Ticket) (*Ticket, error) {
	var ret *Ticket

	b, err := json.Marshal(&ticket)
	if err != nil {
		return nil, err
	}

	// Post JSON request to FreshDesk
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s.freshdesk.com/api/v2/tickets", c.Domain), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.API, "")
	req.Header.Add("Content-type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	// Check the status
	if res.StatusCode != 200 {
		return nil, errors.New("Freshdesk server didn't like the request")
	}

	// Grab the JSON response
	if err = json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Ticket represents a single helpdesk ticket
type Ticket struct {
	Email       string   `conform:"email" json:"email"`
	Name        string   `conform:"name" json:"name"`
	Subject     string   `conform:"trim,title" json:"subject"`
	Description string   `conform:"trim" json:"description"`
	Type        string   `conform:"trim" json:"ticket_type"`
	Status      status   `json:"status"`
	Priority    priority `json:"priority"`
	Source      source   `json:"source"`
}
