package moov

import (
	"time"
)

type Ticket struct {
	ID string `json:"ticketID"`

	Number  int           `json:"number"`
	Title   string        `json:"title"`
	Contact TicketContact `json:"contact"`
	Status  TicketStatus  `json:"status"`

	CreatedOn       time.Time  `json:"createdOn"`
	UpdatedOn       time.Time  `json:"updatedOn"`
	LatestMessageOn *time.Time `json:"latestMessageOn"`
	ClosedOn        *time.Time `json:"closedOn"`
}

type TicketContact struct {
	Email string `json:"email"`

	// Optional name of the contact
	Name string `json:"name,omitempty"`
}

type TicketStatus string

const (
	TicketStatusNew        TicketStatus = "new"
	TicketStatusInProgress TicketStatus = "in-progress"
	TicketStatusOnHold     TicketStatus = "on-hold"
	TicketStatusClosed     TicketStatus = "closed"
)

type TicketMessage struct {
	Author string    `json:"author"`
	Body   string    `json:"body"`
	SentOn time.Time `json:"sentOn"`
}

type ListTicket struct {
	Items []Ticket `json:"items"`

	// NextPage is nil if there are no more pages available.
	NextPage *ListTicketNextPage `json:"nextPage,omitempty"`
}

type ListTicketNextPage struct {
	Cursor string `json:"cursor"`
}

type CreateTicket struct {
	Title   string        `json:"title"`
	Body    string        `json:"body"`
	Contact TicketContact `json:"contact"`
}

type UpdateTicket struct {
	// Optional status to update
	Status UpdateTicketStatus `json:"status,omitempty"`
}

type UpdateTicketStatus string

const (
	UpdateTicketStatusClosed UpdateTicketStatus = "closed"
)
