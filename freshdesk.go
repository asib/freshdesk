package freshdesk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  TYPES, CONSTANTS
//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type source int
type status int
type priority int

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

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  REQUEST
//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Request object for API request
type Request struct {
	Domain string
	API    string
}

// Do ...
func (r *Request) Do(path string, in, out interface{}) error {
	// Create client
	client := &http.Client{}
	// Turn the struct into JSON bytes
	b, _ := json.Marshal(&in)
	// Post JSON request to FreshDesk
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://%s.freshdesk.com%s", r.Domain, path), bytes.NewReader(b))
	req.SetBasicAuth(r.API, "")
	req.Header.Add("Content-type", "application/json")
	res, e := client.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()
	// Check the status
	if res.StatusCode != 200 {
		return errors.New("Freshdesk server didn't like the request")
	}
	// Grab the JSON response
	if e = json.NewDecoder(res.Body).Decode(out); e != nil {
		return e
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  TICKETS
//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

// NewTicket represents helpdesk ticket request
type NewTicket struct {
	Ticket *Ticket `json:"helpdesk_ticket"`
	CC     string  `conform:"trim" json:"cc_emails"`
}

type createTicketFields struct {
	ID int `json:"id"`
}

type createTicketResponse struct {
	Helpdesk *createTicketFields `json:"helpdesk_ticket"`
}

// Create posts a new ticket to FreshDesk
func (h *NewTicket) Create(r *Request) error {
	out := new(createTicketResponse)
	if e := r.Do("/helpdesk/tickets.json", h, out); e != nil {
		return e
	}
	if out.Helpdesk.ID == 0 {
		return errors.New("Freshdesk didn't return valid ticket ID")
	}
	return nil
}
