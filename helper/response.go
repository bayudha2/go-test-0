package helper

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithMultiError(w http.ResponseWriter, code int, messages []string) {
	type errorStruct struct {
		Error string `json:"error"`
	}

	errPayload := []errorStruct{}
	for _, val := range messages {
		errPayload = append(errPayload, errorStruct{
			Error: val,
		})
	}

	RespondWithJSON(w, code, map[string]interface{}{"errors": errPayload})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
