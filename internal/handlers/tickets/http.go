package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/syned13/ticket-support-back/internal/models"
	ticketsService "github.com/syned13/ticket-support-back/internal/service/tickets"
	"github.com/syned13/ticket-support-back/pkg/httputils"
)

var (
	// ErrMissingContentType missing content type
	ErrMissingContentType = httputils.NewBadRequestError("missing content type")
	// ErrInvalidBody invalid request body
	ErrInvalidBody = httputils.NewBadRequestError("invalid request body")
	// ErrInvalidTokenSigningMethod invalid token signing method
	ErrInvalidTokenSigningMethod = errors.New("invalid token signing method")
	// ErrInvalidID invalid pagination id start
	ErrInvalidID = httputils.NewBadRequestError("invalid pagination start id")
)

type claims struct {
	UserType string `json:"userType"`
	Sub      string `json:"sub"`
}

type HTTPHandler interface {
	HandleCreateTicket(ctx context.Context) http.HandlerFunc
	HandleGetTickets(ctx context.Context) http.HandlerFunc
	HandleGetTicket(ctx context.Context) http.HandlerFunc
	HandleGetChanges(ctx context.Context) http.HandlerFunc
	HandlePreflightRequest() http.HandlerFunc
}

type httpHandler struct {
	service ticketsService.Service
}

func SetupRoutes(ctx context.Context, service ticketsService.Service, router *mux.Router) {
	handler := httpHandler{service: service}

	router.HandleFunc("/tickets", authMiddleWare(handler.HandleCreateTicket(ctx))).Methods(http.MethodPost)
	router.HandleFunc("/tickets", authMiddleWare(handler.HandleGetTickets(ctx))).Methods(http.MethodGet)
	router.HandleFunc("/tickets", handler.HandlePreflightRequest()).Methods(http.MethodOptions)

	router.HandleFunc("/tickets/{id}", authMiddleWare(handler.HandleGetTicket(ctx))).Methods(http.MethodGet)
	router.HandleFunc("/tickets/{id}", authMiddleWare(handler.HandleUpdateTicket(ctx))).Methods(http.MethodPatch)
	router.HandleFunc("/tickets/{id}", handler.HandlePreflightRequest()).Methods(http.MethodGet)

	router.HandleFunc("/changes", authMiddleWare(handler.HandleGetChanges(ctx))).Methods(http.MethodGet)
	router.HandleFunc("/changes", handler.HandlePreflightRequest()).Methods(http.MethodOptions)
}

func (h httpHandler) HandlePreflightRequest() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		setupPreflightResponse(&rw, r)
	}
}

func setupPreflightResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (h httpHandler) HandleCreateTicket(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := validateContentType(*r)
		if err != nil {
			httputils.RespondWithError(rw, err)
			return
		}

		ticket := models.Ticket{}
		err = json.NewDecoder(r.Body).Decode(&ticket)
		if err != nil {
			httputils.RespondWithError(rw, ErrInvalidBody)
			return
		}

		creatorIDStr := r.Header.Get("sub")
		if creatorIDStr == "" {
			httputils.RespondWithError(rw, errors.New("missing creator id"))
			return
		}

		var creatorID int64
		if id, err := strconv.ParseInt(creatorIDStr, 10, 64); err == nil {
			creatorID = id
		} else {
			httputils.RespondWithError(rw, errors.New("invalid creator id"))
			return
		}

		ticket.CreatorID = creatorID

		createdTicket, err := h.service.CreateTicket(ctx, ticket)
		if err != nil {
			httputils.RespondWithError(rw, err)
			return
		}

		httputils.RespondJSON(rw, http.StatusCreated, createdTicket)
	}
}

func (h httpHandler) HandleGetTickets(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		lastIDStr := r.URL.Query().Get("after_id")

		var lastID int64 = 0

		if lastIDStr != "" {
			if id, err := strconv.ParseInt(lastIDStr, 10, 64); err == nil {
				lastID = id
			} else {
				httputils.RespondWithError(rw, ErrInvalidID)
				return
			}
		}

		userIDStr := r.Header.Get("sub")
		var userID int64

		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = id
		} else {
			fmt.Println("parsing_sub_failed: " + err.Error())
			httputils.RespondWithError(rw, errors.New("invalid user id"))
			return
		}

		userType := mux.Vars(r)["userType"]

		response, err := h.service.GetTickets(ctx, userID, models.UserType(userType), lastID)
		if err != nil {
			fmt.Println("getting_tickets_failed: " + err.Error())
			httputils.RespondWithError(rw, err)
			return
		}

		httputils.RespondJSON(rw, http.StatusOK, response)
	}
}

func (h httpHandler) HandleGetTicket(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// TODO: take the sub from headers and verify the one requesting the ticket is either the creatoe or an admin
		ticketIDStr := vars["id"]
		if ticketIDStr == "" {
			httputils.RespondWithError(rw, errors.New("missing ticket id"))
			return
		}

		var ticketID int64
		if id, err := strconv.ParseInt(ticketIDStr, 10, 64); err == nil {
			ticketID = id
		} else {
			httputils.RespondWithError(rw, errors.New("invalid ticket id"))
			return
		}

		ticket, err := h.service.GetTicket(ctx, ticketID)
		if err != nil {
			fmt.Println("getting_ticket_failed: " + err.Error())
			httputils.RespondWithError(rw, err)
			return
		}

		httputils.RespondJSON(rw, http.StatusOK, ticket)
	}
}

func (h httpHandler) HandleUpdateTicket(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := validateContentType(*r)
		if err != nil {
			httputils.RespondWithError(rw, err)
			return
		}

		vars := mux.Vars(r)

		// TODO: take the sub from headers and verify the one requesting the ticket is either the creatoe or an admin
		ticketIDStr := vars["id"]
		if ticketIDStr == "" {
			httputils.RespondWithError(rw, errors.New("missing ticket id"))
			return
		}

		var ticketID int64
		if id, err := strconv.ParseInt(ticketIDStr, 10, 64); err == nil {
			ticketID = id
		} else {
			httputils.RespondWithError(rw, errors.New("invalid ticket id"))
			return
		}

		patchRequest := httputils.PatchRequest{}
		err = json.NewDecoder(r.Body).Decode(&patchRequest)
		if err != nil {
			httputils.RespondWithError(rw, ErrInvalidBody)
			return
		}

		updatedTicket, err := h.service.UpdateTicket(ctx, patchRequest, ticketID)
		if err != nil {
			httputils.RespondWithError(rw, err)
			return
		}

		httputils.RespondJSON(rw, http.StatusOK, updatedTicket)
	}
}

func (h httpHandler) HandleGetChanges(ctx context.Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("sub")
		var userID int64

		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = id
		} else {
			fmt.Println("parsing_sub_failed: " + err.Error())
			httputils.RespondWithError(rw, errors.New("invalid user id"))
			return
		}

		changes, err := h.service.GetTicketChanges(ctx, userID)
		if err != nil {
			fmt.Println("getting_changes_failed: " + err.Error())
			httputils.RespondWithError(rw, err)
			return
		}

		httputils.RespondJSON(rw, http.StatusOK, changes)
	}
}

func authMiddleWare(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		token, err := getToken(*r)
		if err != nil {
			fmt.Println(err.Error())
			httputils.RespondWithError(rw, err)
			return
		}

		authToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("aqui tambien")
				return nil, ErrInvalidTokenSigningMethod
			}

			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil {
			fmt.Println("here")
			httputils.RespondWithError(rw, httputils.ForbiddenError)
			return
		}

		if !authToken.Valid {
			fmt.Println("hier")
			httputils.RespondWithError(rw, httputils.ForbiddenError)
			return
		}

		claims, err := getTokenClaims(authToken)
		if err != nil {
			fmt.Println("aqui")
			httputils.RespondWithError(rw, httputils.ForbiddenError)
			return
		}

		r.Header.Set("sub", claims.Sub)
		r.Header.Set("userType", claims.UserType)

		handler.ServeHTTP(rw, r)
	}
}

func getToken(r http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", httputils.UnauthorizedError
	}

	splittedToken := strings.Split(authHeader, "Bearer ")
	if len(splittedToken) < 2 {
		return "", httputils.UnauthorizedError
	}

	return splittedToken[1], nil
}

func getTokenClaims(token *jwt.Token) (claims, error) {
	claimsBytes, err := json.Marshal(token.Claims)
	if err != nil {
		fmt.Println(err.Error())
		return claims{}, err
	}

	tokenClaims := claims{}
	err = json.Unmarshal(claimsBytes, &tokenClaims)
	if err != nil {
		fmt.Println("yes yes")
		fmt.Println(err.Error())
		return claims{}, err
	}

	return tokenClaims, nil
}

func validateContentType(r http.Request) error {
	if r.Header.Get("Content-Type") == "" {
		return ErrMissingContentType
	}

	return nil
}
