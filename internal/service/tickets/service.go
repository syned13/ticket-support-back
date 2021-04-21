package service

import (
	"context"
	"errors"

	"github.com/syned13/ticket-support-back/internal/models"
	ticketsRepository "github.com/syned13/ticket-support-back/internal/repositories/tickets"
	usersRepository "github.com/syned13/ticket-support-back/internal/repositories/users"
	"github.com/syned13/ticket-support-back/pkg/httputils"
)

var (
	ErrMissingTitle       = httputils.NewBadRequestError("missing title")
	ErrMissingDescription = httputils.NewBadRequestError("missing description")
	ErrMissingType        = httputils.NewBadRequestError("missing type")
	ErrMissingSeverity    = httputils.NewBadRequestError("missing severity")
	ErrMissingStatus      = httputils.NewBadRequestError("missing status")
	ErrMissingPriority    = httputils.NewBadRequestError("missing priority")
	ErrInvalidTicketType  = httputils.NewBadRequestError("invalid type")

	ErrMissingCreatorID = errors.New("missing priority")
)

type service struct {
	ticketsRepo ticketsRepository.Repository
	usersRepo   usersRepository.Repository
}

func New(ticketsRepo ticketsRepository.Repository, usersRepo usersRepository.Repository) Service {
	return service{
		ticketsRepo: ticketsRepo,
		usersRepo:   usersRepo,
	}
}

func (s service) CreateTicket(ctx context.Context, ticket models.Ticket) (models.Ticket, error) {
	err := validateCreateTicketParams(ticket)
	if err != nil {
		return models.Ticket{}, err
	}

	ticket.Status = models.TicketTypePending

	createdTicket, err := s.ticketsRepo.SaveTicket(ctx, ticket)
	if err != nil {
		return models.Ticket{}, err
	}

	return createdTicket, nil
}

func validateCreateTicketParams(ticket models.Ticket) error {
	if ticket.Title == "" {
		return ErrMissingTitle
	}

	if ticket.Description == "" {
		return ErrMissingDescription
	}

	if ticket.Type == "" {
		return ErrMissingType
	}

	if models.IsValidTicketType(ticket.Type) {
		return ErrInvalidTicketType
	}

	// TODO: add validations for these numbers
	if ticket.Severity == 0 {
		return ErrMissingSeverity
	}

	if ticket.Priority == 0 {
		return ErrMissingPriority
	}

	if ticket.CreatorID == 0 {
		return ErrMissingCreatorID
	}

	return nil
}

func (s service) GetTickets(ctx context.Context, userID int64, userType models.UserType, lastID int64) (GetTicketsResponse, error) {
	var tickets []models.Ticket
	var err error
	var last int64

	if userType == models.UserTypeAdmin {
		tickets, last, err = s.ticketsRepo.GetTickets(ctx, lastID)
	} else {
		tickets, last, err = s.ticketsRepo.GetTicketsByCreator(ctx, userID, lastID)
	}

	if errors.Is(err, ticketsRepository.ErrNotFound) {
		return GetTicketsResponse{Last: 0, Total: 0, Tickets: []models.Ticket{}}, nil
	}

	if err != nil {
		return GetTicketsResponse{}, err
	}

	return GetTicketsResponse{Tickets: tickets, Last: last, Total: len(tickets)}, nil
}

func (s service) GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error) {
	ticket, err := s.ticketsRepo.GetTicket(ctx, ticketID)
	if errors.Is(err, ticketsRepository.ErrNotFound) {
		return models.Ticket{}, httputils.NewNotFoundError("ticket")
	}

	if err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}
