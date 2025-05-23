package requests

type CreateCardRequest struct {
	Name      string `form:"name"`
	Title     string `form:"title"`
	Photo     string `form:"photo"`
	Telephone string `form:"telephone"`
	Email     string `form:"email"`
	Location  string `form:"location"`
	Website   string `form:"website"`
}

type UpdateCardRequest struct {
	Name      string `form:"name"`
	Title     string `form:"title"`
	Photo     string `form:"photo"`
	Telephone string `form:"telephone"`
	Email     string `form:"email"`
	Location  string `form:"location"`
	Website   string `form:"website"`
}

type UpdateCardStatusRequest struct {
	IsActive bool `form:"is_active"`
}

type AddCardBankRequest struct {
	BankID uint   `form:"bank_id" validate:"required"`
	IBAN   string `form:"iban" validate:"required,min=26,max=26"`
}

type UpdateCardBankRequest struct {
	IBAN string `form:"iban" validate:"required,min=26,max=26"`
}

type AddCardSocialMediaRequest struct {
	SocialMediaID uint   `form:"social_media_id" validate:"required"`
	URL           string `form:"url" validate:"required,url"`
}

type UpdateCardSocialMediaRequest struct {
	URL string `form:"url" validate:"required,url"`
}

func ValidateCreateCardRequest(c *fiber.Ctx) error {
	var req CreateCardRequest
	errorMessages := map[string]string{
		"Email_email": "Geçerli bir e-posta adresi giriniz",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/cards/create"); err != nil {
		return err
	}

	c.Locals("createCardRequest", req)
	return c.Next()
}

func ValidateUpdateCardRequest(c *fiber.Ctx) error {
	var req UpdateCardRequest
	errorMessages := map[string]string{
		"Email_email": "Geçerli bir e-posta adresi giriniz",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/cards"); err != nil {
		return err
	}

	c.Locals("updateCardRequest", req)
	return c.Next()
}

func ValidateAddCardBankRequest(c *fiber.Ctx) error {
	var req AddCardBankRequest
	errorMessages := map[string]string{
		"BankID_required": "Banka seçimi zorunludur",
		"IBAN_required":   "IBAN zorunludur",
		"IBAN_min":        "IBAN 26 karakter olmalıdır",
		"IBAN_max":        "IBAN 26 karakter olmalıdır",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/cards"); err != nil {
		return err
	}

	c.Locals("addCardBankRequest", req)
	return c.Next()
}

func ValidateUpdateCardBankRequest(c *fiber.Ctx) error {
	var req UpdateCardBankRequest
	errorMessages := map[string]string{
		"IBAN_required": "IBAN zorunludur",
		"IBAN_min":      "IBAN 26 karakter olmalıdır",
		"IBAN_max":      "IBAN 26 karakter olmalıdır",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/cards"); err != nil {
		return err
	}

	c.Locals("updateCardBankRequest", req)
	return c.Next()
}

func ValidateAddCardSocialMediaRequest(c *fiber.Ctx) error {
	var req AddCardSocialMediaRequest
	errorMessages := map[string]string{
		"SocialMediaID_required": "Sosyal medya seçimi zorunludur",
		"URL_required":          "URL zorunludur",
		"URL_url":              "Geçerli bir URL giriniz",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/cards"); err != nil {
		return err
	}

	c.Locals("addCardSocialMediaRequest", req)
	return c.Next()
}

func ValidateUpdateCardSocialMediaRequest(c *fiber.Ctx) error {
	var req UpdateCardSocialMediaRequest
	errorMessages := map[string]string{
		"URL_required": "URL zorunludur",
		"URL_url":     "Geçerli bir URL giriniz",
	}

	if err := validateRequest(c, &req, errorMessages, "/panel/cards"); err != nil {
		return err
	}

	c.Locals("updateCardSocialMediaRequest", req)
	return c.Next()
}
