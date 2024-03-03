package requester

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequester(t *testing.T) {
	t.Run("should return a new requester", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}))

		request := NewRequester(0)
		response, err := request.Get(fmt.Sprintf("%s/ping", server.URL))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		resp, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, "pong", string(resp))

	})

}
