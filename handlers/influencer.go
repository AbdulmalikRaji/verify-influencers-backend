package handlers

import (
	"log"

	"github.com/abdulmalikraji/verify-influencers-backend/dto"
	"github.com/abdulmalikraji/verify-influencers-backend/services"
	"github.com/gofiber/fiber/v2"
)

type InfluencerController interface {
	GetInfluencer(c *fiber.Ctx) error
	GetAllInfluencers(c *fiber.Ctx) error
}

type influencerController struct {
	service services.InfluencerService
}

func NewInfluencerController(service services.InfluencerService) InfluencerController {
	return influencerController{
		service: service,
	}
}

func (s influencerController) GetInfluencer(c *fiber.Ctx) error {
	var influencer dto.GetInfluencerRequest
	err := c.QueryParser(&influencer)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}
	if influencer.ID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	response, status, err := s.service.GetInfluencer(c, influencer)

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
		"message": "influencers found successfully",
		"data":    response,
	})

	return err
}

func (s influencerController) GetAllInfluencers(c *fiber.Ctx) error {
	
	response, status, err := s.service.GetAllInfluencers(c)

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
		"message": "influencers found successfully",
		"data":    response,
	})

	return err
}
