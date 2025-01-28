package handlers

import (
	"log"

	"github.com/abdulmalikraji/verify-influencers-backend/dto"
	"github.com/abdulmalikraji/verify-influencers-backend/services"
	"github.com/gofiber/fiber/v2"
)

type ClaimController interface {
	GetInfluencerClaims(c *fiber.Ctx) error
}

type claimController struct {
	service services.ClaimService
}

func NewClaimController(service services.ClaimService) ClaimController {
	return claimController{
		service: service,
	}
}

func (s claimController) GetInfluencerClaims(c *fiber.Ctx) error {
	var claim dto.GetInfluencerClaimsRequest
	err := c.QueryParser(&claim)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	log.Println(claim)

	response, status, err := s.service.GetInfluencerClaims(c, claim)

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
