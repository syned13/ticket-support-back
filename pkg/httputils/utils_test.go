package httputils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBadRequestError(t *testing.T) {
	c := require.New(t)

	message := "the message"

	errorResponse := NewBadRequestError(message)
	c.Equal(http.StatusBadRequest, errorResponse.Code)
	c.Contains(errorResponse.Message, message)
}

func TestNewNotFoundError(t *testing.T) {
	c := require.New(t)

	resourceName := "address"
	errorResponse := NewNotFoundError(resourceName)
	c.Equal(http.StatusNotFound, errorResponse.Code)
	c.Contains(errorResponse.Message, "not found")
	c.Contains(errorResponse.Message, resourceName)
}

func TestRespondJSON(t *testing.T) {
	c := require.New(t)

	w := httptest.NewRecorder()

	RespondJSON(w, http.StatusOK, map[string]string{"message": "hello"})

	body, err := ioutil.ReadAll(w.Body)
	c.Nil(err)
	c.Contains(string(body), `"message":"hello"`)
}

func TestRespondWithError(t *testing.T) {
	c := require.New(t)

	w := httptest.NewRecorder()

	RespondWithError(w, errors.New("some error"))
	c.Equal(http.StatusInternalServerError, w.Code)
}
