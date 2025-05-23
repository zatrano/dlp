package requests

type CreateInvitationParticipantRequest struct {
	Title       string `form:"title" validate:"required"`
	PhoneNumber string `form:"phone_number" validate:"required"`
	GuestCount  int    `form:"guest_count" validate:"required,min=1"`
}

type UpdateInvitationParticipantRequest struct {
	Title       string `form:"title" validate:"required"`
	PhoneNumber string `form:"phone_number" validate:"required"`
	GuestCount  int    `form:"guest_count" validate:"required,min=1"`
}

func ValidateCreateInvitationParticipantRequest(c *fiber.Ctx) error {
	var req CreateInvitationParticipantRequest
	errorMessages := map[string]string{
		"Title_required":       "İsim zorunludur",
		"PhoneNumber_required": "Telefon numarası zorunludur",
		"GuestCount_required": "Misafir sayısı zorunludur",
		"GuestCount_min":      "Misafir sayısı en az 1 olmalıdır",
	}

	if err := validateRequest(c, &req, errorMessages, "/invitations"); err != nil {
		return err
	}

	c.Locals("createInvitationParticipantRequest", req)
	return c.Next()
}

func ValidateUpdateInvitationParticipantRequest(c *fiber.Ctx) error {
	var req UpdateInvitationParticipantRequest
	errorMessages := map[string]string{
		"Title_required":       "İsim zorunludur",
		"PhoneNumber_required": "Telefon numarası zorunludur",
		"GuestCount_required": "Misafir sayısı zorunludur",
		"GuestCount_min":      "Misafir sayısı en az 1 olmalıdır",
	}

	if err := validateRequest(c, &req, errorMessages, "/invitations"); err != nil {
		return err
	}

	c.Locals("updateInvitationParticipantRequest", req)
	return c.Next()
}
