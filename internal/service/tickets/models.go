package service

import "github.com/syned13/ticket-support-back/internal/models"

type GetTicketsResponse struct {
	Tickets []models.Ticket `json:"tickets"`
	Last    int64           `json:"last"`
	Total   int             `json:"total"`
}
