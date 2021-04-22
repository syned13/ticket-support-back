package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/syned13/ticket-support-back/internal/models"
	repository "github.com/syned13/ticket-support-back/internal/repositories/tickets"
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
	ErrInvalidOwnerID     = httputils.NewBadRequestError("invalid owner id")
	ErrInvalidStatus      = httputils.NewBadRequestError("invalid status")

	// ErrMissingPatchOperation missing patch operation
	ErrMissingPatchOperation = httputils.NewBadRequestError("missing patch operation")
	// ErrMissingPatchPath missing patch path
	ErrMissingPatchPath = httputils.NewBadRequestError("missing patch path")
	// ErrMissingPatchValue missing patch value
	ErrMissingPatchValue = httputils.NewBadRequestError("missing patch value")

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
	// TODO: only the own creator or an admin can get a ticket
	ticket, err := s.ticketsRepo.GetTicket(ctx, ticketID)
	if errors.Is(err, ticketsRepository.ErrNotFound) {
		return models.Ticket{}, httputils.NewNotFoundError("ticket")
	}

	if err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (s service) UpdateTicket(ctx context.Context, request httputils.PatchRequest, ticketID int64) (models.Ticket, error) {
	ticketChange := models.TicketChange{}
	updatedStatus := false

	ticket, err := s.ticketsRepo.GetTicket(ctx, ticketID)
	if errors.Is(err, ticketsRepository.ErrNotFound) {
		return models.Ticket{}, httputils.NewNotFoundError("ticket")
	}

	if err != nil {
		return models.Ticket{}, err
	}

	ticketChange.TicketID = ticket.TicketID
	ticketChange.CreatorID = ticket.CreatorID

	// 	TODO: move this process to a new function to separate concerns
	for _, op := range request {
		if op.Op == "" {
			return models.Ticket{}, ErrMissingPatchOperation
		}

		if op.Path == "" {
			return models.Ticket{}, ErrMissingPatchPath
		}

		if op.Value == "" {
			return models.Ticket{}, ErrMissingPatchValue
		}

		if op.Op != "update" { // TODO: remove maginc string
			return models.Ticket{}, httputils.NewBadRequestError("invalid patch operation: " + op.Op)
		}

		switch op.Path {
		case "ownerID":
			id, ok := op.Value.(int64)
			if !ok {
				return models.Ticket{}, ErrInvalidOwnerID
			}

			ticket.OwnerID = &id
		case "status":
			status, ok := op.Value.(string)
			if !ok {
				return models.Ticket{}, ErrInvalidStatus
			}

			updatedStatus = true
			// TODO: validate status, only valid statuses and only final statuses
			ticket.Status = models.TicketStatus(status)
			ticketChange.To = models.TicketStatus(status)
		}
	}

	updatedTicket, err := s.ticketsRepo.UpdateTicket(ctx, ticket)
	if errors.Is(err, ticketsRepository.ErrNotFound) {
		return models.Ticket{}, httputils.NewNotFoundError("ticket")
	}

	if errors.Is(err, repository.ErrNothingToUpdate) {
		return models.Ticket{}, httputils.NewBadRequestError("nothing to update")
	}

	if updatedStatus {
		fmt.Println(ticketChange)
		err = s.ticketsRepo.SaveTicketChange(ctx, ticketChange)
		if err != nil {
			fmt.Println("could not add change to change log")
		}
	}

	return updatedTicket, nil
}

func (s service) GetTicketChanges(ctx context.Context, creatorID int64) ([]models.TicketChange, error) {
	return s.ticketsRepo.GetTicketChanges(ctx, creatorID)
}
