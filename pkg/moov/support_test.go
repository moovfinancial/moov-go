package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Tickets(t *testing.T) {
	mc := NewTestClient(t)

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	var (
		ticket1 *moov.Ticket = nil
		ticket2 *moov.Ticket = nil
	)

	const (
		ticketTitle     = "moov-go SDK Test Ticket"
		ticketBody      = "Testing the moov-go SDK ticket creation"
		ticketAuthor    = "Override Author"
		ticketForeignID = "ticket-foreign-id"
	)

	t.Run("create ticket", func(t *testing.T) {

		var err error
		ticket1, err = mc.CreateTicket(BgCtx(), customer.AccountID, moov.CreateTicket{
			Title: ticketTitle,
			Body:  ticketBody,
			Contact: moov.TicketContact{
				Email: "moovbot@moov.io",
			},
			Author:    ticketAuthor,
			ForeignID: ticketForeignID,
		})

		require.NoError(t, err)
		require.NotNil(t, ticket1)
		require.Equal(t, ticketTitle, ticket1.Title)

		ticket2, err = mc.CreateTicket(BgCtx(), customer.AccountID, moov.CreateTicket{
			Title: ticketTitle,
			Body:  ticketBody,
			Contact: moov.TicketContact{
				Email: "moovbot@moov.io",
			},
		})

		require.NoError(t, err)

		t.Logf("ticket: %+v\n", ticket1)
		t.Logf("ticket: %+v\n", ticket2)
	})

	t.Run("list tickets", func(t *testing.T) {
		require.NotNil(t, ticket1)
		require.NotNil(t, ticket2)

		tickets, err := mc.ListTickets(BgCtx(), customer.AccountID, moov.WithTicketCount(1))
		require.NoError(t, err)
		require.Len(t, tickets.Items, 1)
		require.NotNil(t, tickets.NextPage)
		require.Contains(t, tickets.Items, *ticket2)

		tickets, err = mc.ListTickets(BgCtx(), customer.AccountID, moov.WithTicketCursor(tickets.NextPage.Cursor))
		require.NoError(t, err)
		require.Len(t, tickets.Items, 1)
		require.Nil(t, tickets.NextPage)
		require.Contains(t, tickets.Items, *ticket1)

		tickets, err = mc.ListTickets(BgCtx(), customer.AccountID, moov.WithTicketStatus(moov.TicketStatusNew))
		require.NoError(t, err)
		require.Len(t, tickets.Items, 2)
		require.Nil(t, tickets.NextPage)
		require.Contains(t, tickets.Items, *ticket1)
		require.Contains(t, tickets.Items, *ticket2)

		tickets, err = mc.ListTickets(BgCtx(), customer.AccountID, moov.WithTicketForeignID(ticketForeignID))
		require.NoError(t, err)
		require.Len(t, tickets.Items, 1)
		require.Nil(t, tickets.NextPage)
		require.Contains(t, tickets.Items, *ticket1)
	})

	t.Run("list ticket messages", func(t *testing.T) {
		require.NotNil(t, ticket1)

		ticketMessages, err := mc.ListTicketMessages(BgCtx(), customer.AccountID, ticket1.ID)
		require.NoError(t, err)
		require.Len(t, ticketMessages, 1)

		ticketMessage := ticketMessages[0]
		require.Equal(t, ticketBody, ticketMessage.Body)
		require.Equal(t, ticketAuthor, ticketMessage.Author)
	})

	t.Run("get ticket", func(t *testing.T) {
		require.NotNil(t, ticket1)

		foundTicket, err := mc.GetTicket(BgCtx(), customer.AccountID, ticket1.ID)
		require.NoError(t, err)
		require.NotNil(t, foundTicket)
		require.Equal(t, *ticket1, *foundTicket)
	})

	t.Run("update ticket", func(t *testing.T) {
		require.NotNil(t, ticket1)

		updatedTicket, err := mc.UpdateTicket(BgCtx(), customer.AccountID, ticket1.ID, moov.UpdateTicket{
			Status: moov.UpdateTicketStatusClosed,
		})

		require.NoError(t, err)
		require.NotNil(t, updatedTicket)
		require.Equal(t, moov.TicketStatusClosed, updatedTicket.Status)
	})
}
