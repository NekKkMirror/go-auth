package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ParseUUIDParam safely parses and validates UUID params from request context.
func ParseUUIDParam(c *fiber.Ctx, param string) (uuid.UUID, error) {
	raw := c.Query(param)
	if raw == "" {
		return uuid.Nil, errors.New("missing or empty UUID parameter: " + param)
	}

	parsedUUID, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID format for parameter: " + param)
	}

	return parsedUUID, nil
}
