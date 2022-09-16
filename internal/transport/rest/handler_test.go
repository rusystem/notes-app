package rest

import (
	"github.com/rusystem/notes-app/internal/service"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(&service.Service{})

	require.IsType(t, &Handler{}, h)
}

func TestHandler_InitRoutes(t *testing.T) {
	h := NewHandler(&service.Service{})

	ts := httptest.NewServer(h.InitRoutes())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
