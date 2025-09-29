package tests

import (
	"WebSportwareShop/internal/handlers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOfProductsHandler(t *testing.T) {

	t.Run("ValidListOfProducts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		res := httptest.NewRecorder()

		handlers.ListOfProductsHandle(res, req)

		if res.Code != http.StatusOK {
			t.Fatalf("Server give invalid status %d", res.Code)
		}
		var expBody map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&expBody); err != nil {
			t.Fatalf("Unexpected Body structure %v", res.Body)
		}

		fmt.Println("HandlerOF List Product successfully working")
	})
}

func TestLoginHandler(t *testing.T) {
	t.Run("Valid Login Handler", func(t *testing.T) {
		loginReq := handlers.LoginReq{Email: "firstuser@gmail.com", Password: "1234"}
		res := httptest.NewRecorder()
		loginReqBytes, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginReqBytes))
		req.Header.Set("Content-Type", "application/json")
		handlers.LoginHandle(res, req)

		if res.Code != http.StatusOK {
			t.Fatalf("Server in trouble")
		}

		var expBody map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&expBody); err != nil {
			t.Fatalf("Invalid output %s", res.Body)
		}

		fmt.Println("Handler of login handler successfully working")
	})
}

// Unit and Integration testings

// functional testing is about that it works and give right results it give me right answer
// non-functional testing is about that it works the right way like does it work properly calculate properly

/*Smoke Testing

Goal: Quick check if the system basically runs.

Example: "Does the app start and return 200 on /health endpoint?"*/
