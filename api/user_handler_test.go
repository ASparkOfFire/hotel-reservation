package api

import (
	"bytes"
	"encoding/json"
	"github.com/asparkoffire/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "James",
		LastName:  "Madison",
		Email:     "james@example.com",
		Password:  "1234567890",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	_ = json.NewDecoder(resp.Body).Decode(&user)

	if len(user.Id) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the JSON response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName: %s, got: %s", user.FirstName, params.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName: %s, got: %s", user.LastName, params.FirstName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email: %s, got: %s", user.Email, params.Email)
	}
}
