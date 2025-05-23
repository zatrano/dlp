package requests

type CreateInvitationRequest struct {
	CategoryID    uint      `form:"category_id" validate:"required"`
	Template      string    `form:"template" validate:"required"`
	Type         string    `form:"type" validate:"required"`
	Title        string    `form:"title"`
	Image        string    `form:"image"`
	Description  string    `form:"description"`
	Venue        string    `form:"venue"`
	Address      string    `form:"address"`
	Location     string    `form:"location"`
	Link         string    `form:"link"`
	Telephone    string    `form:"telephone"`
	Note         string    `form:"note"`
	Date         string    `form:"date" validate:"required"`
	Time         string    `form:"time" validate:"required"`
	IsParticipant bool     `form:"is_participant"`
}

type UpdateInvitationRequest struct {
	CategoryID    uint      `form:"category_id" validate:"required"`
	Template      string    `form:"template" validate:"required"`
	Type         string    `form:"type" validate:"required"`
	Title        string    `form:"title"`
	Image        string    `form:"image"`
	Description  string    `form:"description"`
	Venue        string    `form:"venue"`
	Address      string    `form:"address"`
	Location     string    `form:"location"`
	Link         string    `form:"link"`
	Telephone    string    `form:"telephone"`
	Note         string    `form:"note"`
	Date         string    `form:"date" validate:"required"`
	Time         string    `form:"time" validate:"required"`
	IsParticipant bool     `form:"is_participant"`
}

type UpdateInvitationStatusRequest struct {
	IsConfirmed bool `form:"is_confirmed"`
}

func ValidateCreateInvitationRequest(c *fiber.Ctx) error {
	var req CreateInvitationRequest
	errorMessages := map[string]string{
		"CategoryID_required": "Kategori seçimi zorunludur",
		"Template_required":   "Şablon seçimi zorunludur",
		"Type_required":      "Tip seçimi zorunludur",
		"Date_required":      "Tarih seçimi zorunludur",
		"Time_required":      "Saat seçimi zorunludur",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/invitations/create"); err != nil {
		return err
	}

	c.Locals("createInvitationRequest", req)
	return c.Next()
}

func ValidateUpdateInvitationRequest(c *fiber.Ctx) error {
	var req UpdateInvitationRequest
	errorMessages := map[string]string{
		"CategoryID_required": "Kategori seçimi zorunludur",
		"Template_required":   "Şablon seçimi zorunludur",
		"Type_required":      "Tip seçimi zorunludur",
		"Date_required":      "Tarih seçimi zorunludur",
		"Time_required":      "Saat seçimi zorunludur",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/invitations"); err != nil {
		return err
	}

	c.Locals("updateInvitationRequest", req)
	return c.Next()
}

func ValidateUpdateInvitationStatusRequest(c *fiber.Ctx) error {
	var req UpdateInvitationStatusRequest
	errorMessages := map[string]string{}

	if err := validateRequest(c, &req, errorMessages, "/dashboard/invitations"); err != nil {
		return err
	}

	c.Locals("updateInvitationStatusRequest", req)
	return c.Next()
}
