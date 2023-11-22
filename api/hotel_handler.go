package api

import (
	"github.com/asparkoffire/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	Store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		Store: store,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return ErrInvalidID()
	}
	filter := db.Map{"_id": oid}
	hotel, err := h.Store.Hotel.GetHotels(c.Context(), filter, nil)
	if len(hotel) == 0 {
		return ErrNotFound()
	}
	if err != nil {
		return ErrInvalidID()
	}
	return c.JSON(hotel)
}

type ResourceResp struct {
	Results int `json:"results"`
	Page    int `json:"page"`
	Data    any `json:"data"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams

	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{
		"rating": params.Rating,
	}
	hotels, err := h.Store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrInvalidID()
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(params.Page),
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return ErrInvalidID()
	}
	filter := db.Map{"hotelID": oid}
	rooms, err := h.Store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrNotFound()
	}
	return c.JSON(rooms)
}
