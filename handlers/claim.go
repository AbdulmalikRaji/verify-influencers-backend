package handlers

import (
	"log"
	"time"

	"github.com/abdulmalikraji/verify-influencers-backend/dto"
	"github.com/abdulmalikraji/verify-influencers-backend/services"
	"github.com/gofiber/fiber/v2"
)

type ClaimController interface {
	FindInfluencerClaims(c *fiber.Ctx) error
}

type claimController struct {
	service services.ClaimService
}

func NewClaimController(service services.ClaimService) ClaimController {
	return claimController{
		service: service,
	}
}

func (s claimController) FindInfluencerClaims(c *fiber.Ctx) error {
	var claim dto.FindInfluencerClaimsRequest
	err := c.QueryParser(&claim)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Manually parse start_date and end_date
	if c.Query("start_date") != "" {
		claim.StartDate, err = time.Parse("2006-01-02", c.Query("start_date"))
		if err != nil {
			log.Println("Invalid start_date format:", err)
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"message": "Invalid start_date format. Use YYYY-MM-DD.",
				"data":    nil,
			})
		}
	}

	if c.Query("end_date") != "" {
		claim.EndDate, err = time.Parse("2006-01-02", c.Query("end_date"))
		if err != nil {
			log.Println("Invalid end_date format:", err)
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"message": "Invalid end_date format. Use YYYY-MM-DD.",
				"data":    nil,
			})
		}
	}

	response, status, err := s.service.FindInfluencerClaims(c, claim)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"status":  status,
			"message": err.Error(),
			"data":    nil,
		})
	}

	err = c.JSON(&fiber.Map{
		"success": true,
		"status":  status,
		"message": "claims found successfully",
		"data":    response,
	})

	return err
}
