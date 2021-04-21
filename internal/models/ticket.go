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
	TicketSeverityLow      TicketSeverity = 1
	TicketSeverityMedium   TicketSeverity = 2
	TicketSeverityHigh     TicketSeverity = 3
	TicketSeverityVeryHigh TicketSeverity = 4
)

type TicketStatus string

const (
	TicketTypePending    TicketStatus = "pending"
	TicketTypeInProgress TicketStatus = "in_progress"
	TicketStatusResolved TicketStatus = "resolved"
)

type TicketPriority int

const (
	TicketPriorityLow      TicketPriority = 1
	TicketPriorityMedium   TicketPriority = 2
	TicketPriorityHigh     TicketPriority = 3
	TicketPriorityVeryHigh TicketPriority = 4
)

// Ticket represents a ticket
type Ticket struct {
	TicketID    int64          `json:"ticketID" db:"id"`
	Title       string         `json:"title" db:"title"`
	Description string         `json:"description" db:"ticket_description"`
	Type        TicketType     `json:"type" db:"ticket_type"`
	Severity    TicketSeverity `json:"severity" db:"severity"`
	Priority    TicketPriority `json:"priority" db:"ticket_priority"`
	Status      TicketStatus   `json:"status" db:"ticket_status"`
	CreatorID   int64          `json:"creatorID" db:"creator_id"`
	OwnerID     *int64         `json:"ownerID,omitempty" db:"owner_id"`
	CreatedAt   *time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time     `json:"updatedAt" db:"updated_at"`
	ResolvedAt  *time.Time     `json:"resolvedAt,omitempty" db:"resolved_at"`
}

func IsValidTicketType(ticketType TicketType) bool {
	return ValidTicketTypes[ticketType]
}
