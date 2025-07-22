package httpresponse

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"err,omitempty"`
}

func MarshalResponse(data interface{}, errMsg string) ([]byte, error) {
	res, err := json.Marshal(Response{Data: data, Err: errMsg})
	if err != nil {
		return nil, err
	}
	return res, nil
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
