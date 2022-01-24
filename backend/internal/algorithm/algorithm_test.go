package algorithm

import (
	"backend/internal/structure"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlgorithm_RunSimulation(t *testing.T) {
	t.Run("should run simulation", func(t *testing.T) {

		sendData := structure.Body{ID: "id", Model: "model", Image: "image"}
		jsonSendData, err := json.Marshal(sendData)

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}
		testServer := httptest.NewServer(http.HandlerFunc(handler))
		defer testServer.Close()
		alg := NewAlgorithm("id", testServer.URL)
		status, err := alg.RunSimulation("demo", jsonSendData)
		assert.NoError(t, err)
		assert.Equal(t, status, http.StatusOK)
	})
}
