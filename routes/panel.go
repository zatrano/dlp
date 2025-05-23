package routes

import (
	handlers "davet.link/handlers/panel"
	"davet.link/middlewares"
	"davet.link/models"

	"github.com/gofiber/fiber/v2"
)

func registerPanelRoutes(app *fiber.App) {
	panelGroup := app.Group("/panel")
	panelGroup.Use(
		middlewares.AuthMiddleware,
		middlewares.StatusMiddleware,
		middlewares.TypeMiddleware(models.Panel),
		middlewares.VerifiedMiddleware,
	)
	panelHomeHandler := handlers.NewPanelHomeHandler()
	panelGroup.Get("/home", panelHomeHandler.HomePage)

	cardHandler := handlers.NewPanelCardHandler()
	panelGroup.Get("/card", cardHandler.ShowCard)
	panelGroup.Get("/card/create", cardHandler.ShowCreateCard)
	panelGroup.Post("/card/create", cardHandler.CreateCard)
	panelGroup.Get("/card/update", cardHandler.ShowUpdateCard)
	panelGroup.Post("/card/update", cardHandler.UpdateCard)

	// Card Bank routes
	panelGroup.Get("/card/banks", cardHandler.ListBanks)
	panelGroup.Post("/card/banks/add", cardHandler.AddBank)
	panelGroup.Post("/card/banks/update/:id", cardHandler.UpdateBank)
	panelGroup.Delete("/card/banks/delete/:id", cardHandler.DeleteBank)

	// Card Social Media routes
	panelGroup.Get("/card/social-media", cardHandler.ListSocialMedia)
	panelGroup.Post("/card/social-media/add", cardHandler.AddSocialMedia)
	panelGroup.Post("/card/social-media/update/:id", cardHandler.UpdateSocialMedia)
	panelGroup.Delete("/card/social-media/delete/:id", cardHandler.DeleteSocialMedia)

	invitationHandler := handlers.NewPanelInvitationHandler()
	panelGroup.Get("/invitations", invitationHandler.ListPanelInvitations)
	panelGroup.Get("/invitations/create", invitationHandler.ShowCreatePanelInvitation)
	panelGroup.Post("/invitations/create", invitationHandler.CreatePanelInvitation)
	panelGroup.Get("/invitations/update/:id", invitationHandler.ShowUpdatePanelInvitation)
	panelGroup.Post("/invitations/update/:id", invitationHandler.UpdatePanelInvitation)
	panelGroup.Delete("/invitations/delete/:id", invitationHandler.DeletePanelInvitation)

	// Invitation Detail routes
	panelGroup.Get("/invitations/:id/details", invitationHandler.ShowInvitationDetails)
	panelGroup.Post("/invitations/:id/details", invitationHandler.UpdateInvitationDetails)

	// Invitation Participant routes
	panelGroup.Get("/invitations/:id/participants", invitationHandler.ListInvitationParticipants)
	panelGroup.Post("/invitations/:id/participants/add", invitationHandler.AddInvitationParticipant)
	panelGroup.Delete("/invitations/:id/participants/:participant_id", invitationHandler.DeleteInvitationParticipant)
}
