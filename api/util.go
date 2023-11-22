package api

import (
	"fmt"
	"github.com/asparkoffire/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthenticatedUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return user, nil
}
