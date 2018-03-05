# Freshdesk API (Go, golang)

Basic API submitting a ticket to [Freshdesk.com](https://freshdesk.com)

## How to use

```go

import (
	"fmt"
	"os"

	"github.com/leebenson/conform"
)

func main() {
	client, err := NewClient("your-domain", "your-api-key")
	if err != nil {
		fmt.Printf("Could not create client: %s\n", err)
		os.Exit(1)
	}

	ticket := &Ticket{
	Email:       "email@example.com",
	Name:        "your name",
	Subject:     "this is a test",
	Type:        "Question",
	Description: "the content of the ticket would go here",
	Status:      freshdesk.Open,
	Priority:    freshdesk.Medium,
	Source:      freshdesk.Portal,
	}

	// optionally check the ticket with conform
	conform.Strings(ticket)

	if _, err := client.CreateTicket(ticket); err != nil {
		fmt.Printf("failed to create ticket: %s", err)
		os.Exit(1)
	}
}

```

## Conform library compatibility

Ticket struct fields work with the optional [conform library](https://github.com/leebenson/conform).
