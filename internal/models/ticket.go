package models

import "time"

type TicketType string

const (
	// TicketTypeSupport support
	TicketTypeSupport TicketType = "support"
	// TicketTypeSugestion suggestion
	TicketTypeSugestion TicketType = "sugestion"
	// TicketTypeAsistance asistance
	TicketTypeAsistance TicketType = "asistance"
)

var (
	ValidTicketTypes = map[TicketType]bool{
		TicketTypeSupport:   true,
		TicketTypeSugestion: true,
		TicketTypeAsistance: true,
	}
)

type TicketSeverity int

const (
	TicketSeverityLow      TicketSeverity = 0
	TicketSeverityMedium   TicketSeverity = 1
	TicketSeverityHigh     TicketSeverity = 2
	TicketSeverityVeryHigh TicketSeverity = 3
)

type TicketStatus string

const (
	TicketTypePending    TicketStatus = "pending"
	TicketTypeInProgress TicketStatus = "in_progress"
	TicketStatusResolved TicketStatus = "resolved"
)

type TicketPriority int

const (
	TicketPriorityLow      TicketPriority = 0
	TicketPriorityMedium   TicketPriority = 1
	TicketPriorityHigh     TicketPriority = 2
	TicketPriorityVeryHigh TicketPriority = 3
)

// Ticket represents a ticket
type Ticket struct {
	TicketID    int64          `json:"ticketID"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        TicketType     `json:"ticketType"`
	Severity    TicketSeverity `json:"severity"`
	Priority    TicketPriority `json:"priority"`
	Status      TicketStatus   `json:"status"`
	CreatorID   int64          `json:"creatorID"`
	OwnerID     int64          `json:"ownerID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	ResolvedAt  time.Time      `json:"resolvedAt"`
}
