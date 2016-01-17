# Freshdesk API (Go, golang)

Basic API submitting a ticket to [Freshdesk.com](https://freshdesk.com)

## How to use

```go

import "github.com/fundedcity/freshdesk"

r := &freshdesk.Request{
	Domain: "your-freshdesk-domain",
	API:    "api-key-here",
}

tk := &freshdesk.NewTicket{
  Ticket: &freshdesk.Ticket{
  	Email:       "email@example.com",
  	Name:        "your name",
  	Subject:     "this is a test",
  	Description: "the content of the ticket would go here",
  	Status:      freshdesk.Open,
  	Priority:    freshdesk.Medium,
  	Source:      freshdesk.Portal,
  }
}

conform.Strings(tk.Ticket) // <-- optionally use "conform" library

tk.Create(r) // <-- returns error || nil

```

## Conform library compatibility

Ticket struct fields work with the optional [conform library](https://github.com/leebenson/conform).
