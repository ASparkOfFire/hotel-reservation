package api

import (
	"encoding/json"
	"fmt"
	"github.com/asparkoffire/hotel-reservation/db/fixtures"
	"github.com/asparkoffire/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "Fake", "User", false)
		user           = fixtures.AddUser(db.Store, "John", "Doe", false)
		hotel          = fixtures.AddHotel(db.Store, "BarHotel", "New York", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 45.0, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.Id, room.ID, from, till)
		bookingHandler = NewBookingHandler(db.Store)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(db.User))
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code, got: %d", resp.StatusCode)
	}
	var bookingResponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}
	have := bookingResponse
	if have.ID != booking.ID {
		t.Fatalf("expected: %s, got: %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected: %s, got: %s", booking.UserID, have.UserID)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 code, got: %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		adminUser      = fixtures.AddUser(db.Store, "Elliot", "Alderson", true)
		user           = fixtures.AddUser(db.Store, "John", "Doe", false)
		hotel          = fixtures.AddHotel(db.Store, "BarHotel", "New York", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 45.0, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.Id, room.ID, from, till)
		bookingHandler = NewBookingHandler(db.Store)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
	)
	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response: %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected: %s, got: %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected: %s, got: %s", booking.UserID, have.UserID)
	}

	//test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected a status code %d got, %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
