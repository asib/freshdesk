package freshdesk

import (
	"testing"

	"github.com/leebenson/conform"
)

func TestTicket(t *testing.T) {
	t.Skip("Need account info to test this")

	client, err := NewClient("", "")
	if err != nil {
		t.Fatalf("Could not create client: %s\n", err)
	}

	ticket := &Ticket{
		Email:       "testuser@example.com",
		Name:        "Cool Testuser",
		Subject:     "this is a test feature request",
		Type:        FeatureRequest,
		Description: "the content of the ticket would go here",
		Status:      Open,
		Priority:    Medium,
		Source:      Portal,
	}

	// optionally check the ticket with conform
	conform.Strings(ticket)

	if _, err := client.CreateTicket(ticket); err != nil {
		t.Fatalf("failed to create ticket: %s", err)
	}
}
