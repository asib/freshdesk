package freshdesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// Source of the ticket
	Source int
	// Status of the ticket
	Status int
	// Priority of the ticket
	Priority int
	// TicketType corresponds to the type of the ticket
	TicketType string
)

const (
	// Email ...
	Email Source = 1 + iota
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
	Open Status = 2 + iota
	// Pending ...
	Pending
	// Resolved ...
	Resolved
	// Closed ...
	Closed
)

const (
	// Low ...
	Low Priority = 1 + iota
	// Medium ...
	Medium
	// High ...
	High
	// Urgent ...
	Urgent
)

const (
	// Question ...
	Question TicketType = "Question"
	// Incident ...
	Incident = "Incident"
	// Problem ...
	Problem = "Problem"
	// FeatureRequest ...
	FeatureRequest = "Feature Request"
	// Lead ...
	Lead = "Lead"
)

// API corresponds to the FreshdeskAPI. It can be used for mock
// implementations for testing.
type API interface {
	CreateTicket(ticket *Ticket) (*Ticket, error)
}

// Client is a freshdesk client that allows access to the freshdesk
// API. It is save to use the client from different goroutines.
type Client struct {
	Domain string
	API    string

	httpClient *http.Client
}

// Ticket represents a single helpdesk ticket
type Ticket struct {
	Email       string     `conform:"email" json:"email"`
	Name        string     `conform:"name" json:"name"`
	Subject     string     `conform:"trim,title" json:"subject"`
	Description string     `conform:"trim" json:"description"`
	Type        TicketType `conform:"trim" json:"type"`

	Tags         []string          `json:"tags,omitempty"`
	CustomFields map[string]string `json:"custom_fields,omitempty"`
	Status       Status            `json:"status"`
	Priority     Priority          `json:"priority"`
	Source       Source            `json:"source"`
}

// NewClient returns a new freshdesk client that uses the
// http.DefaultClient under the hood.
func NewClient(domain, api string) (*Client, error) {
	return &Client{Domain: domain, API: api, httpClient: http.DefaultClient}, nil
}

// CreateTicket creates a new ticket.
func (c *Client) CreateTicket(ticket *Ticket) (*Ticket, error) {
	var ret Ticket

	b, err := json.Marshal(&ticket)
	if err != nil {
		return nil, err
	}

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

	if res.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("freshdesk server didn't like the request (status code: %d): %s", res.StatusCode, bs)
	}

	if err = json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
