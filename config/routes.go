package config

import (
	"github.com/abdulmalikraji/verify-influencers-backend/db/connection"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/claimDao"
	"github.com/abdulmalikraji/verify-influencers-backend/db/dao/influencerDao"
	"github.com/abdulmalikraji/verify-influencers-backend/handlers"
	"github.com/abdulmalikraji/verify-influencers-backend/services"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App, client connection.Client) {

	//repositories
	claimDataAccess := claimDao.New(client)
	influencerDataAccess := influencerDao.New(client) 

	//services
	claimService := services.NewClaimService(claimDataAccess, influencerDataAccess)

	//Handlers
	claimHandler := handlers.NewClaimController(claimService)

	// Routes
	api := app.Group("/api/v1")
	api.Get("/claims", claimHandler.GetInfluencerClaims)

}
