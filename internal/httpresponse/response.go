package httpresponse

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"err,omitempty"`
}

func WriteJSON(w http.ResponseWriter, code int, data interface{}, errMSG string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	res := Response{
		Data: data,
		Err:  errMSG,
	}
	json.NewEncoder(w).Encode(res)
}
