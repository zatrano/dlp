package handlers

import (
	"net/http"
	"strings"

	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CardHandler struct {
	cardService services.ICardService
}

func NewCardHandler() *CardHandler {
	return &CardHandler{cardService: services.NewCardService()}
}

func (h *CardHandler) ListCards(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Kart listesi: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
		params = queryparams.DefaultListParams()
	}

	if params.Page <= 0 {
		params.Page = queryparams.DefaultPage
	}
	if params.PerPage <= 0 || params.PerPage > queryparams.MaxPerPage {
		params.PerPage = queryparams.DefaultPerPage
	}
	if params.SortBy == "" {
		params.SortBy = queryparams.DefaultSortBy
	}
	if params.OrderBy == "" {
		params.OrderBy = queryparams.DefaultOrderBy
	}

	paginatedResult, dbErr := h.cardService.GetAllCards(params)

	renderData := fiber.Map{
		"Title":  "Kartlar",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		logconfig.Log.Error("Kart listesi DB Hatası", zap.Error(dbErr))
		renderData[renderer.FlashErrorKeyView] = "Kartlar getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.Card{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/cards/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *CardHandler) ShowCreateCard(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Kart Ekle",
	})
}

func (h *CardHandler) CreateCard(c *fiber.Ctx) error {
	var req struct {
		Name      string `form:"name"`
		Title     string `form:"title"`
		Photo     string `form:"photo"`
		Telephone string `form:"telephone"`
		Email     string `form:"email"`
		Location  string `form:"location"`
		Website   string `form:"website"`
		IsActive  string `form:"is_active"`
		UserID    uint   `form:"user_id"`
	}
	_ = c.BodyParser(&req)

	if req.Name == "" || req.UserID == 0 {
		return renderCardFormError("Yeni Kart Ekle", req, "Ad ve Kullanıcı alanları zorunludur.", c)
	}

	card := &models.Card{
		Name:      req.Name,
		Title:     req.Title,
		Photo:     req.Photo,
		Telephone: req.Telephone,
		Email:     req.Email,
		Location:  req.Location,
		Website:   req.Website,
		IsActive:  req.IsActive == "true",
		UserID:    req.UserID,
	}

	if err := h.cardService.CreateCard(c.UserContext(), card); err != nil {
		return renderCardFormError("Yeni Kart Ekle", req, "Kart oluşturulamadı: "+err.Error(), c)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/cards", fiber.StatusFound)
}

func (h *CardHandler) ShowUpdateCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	card, err := h.cardService.GetCardByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart bulunamadı.")
		return c.Redirect("/dashboard/cards", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
		"Title": "Kart Düzenle",
		"Card":  card,
	})
}

func (h *CardHandler) UpdateCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	cardID := uint(id)

	var req struct {
		Name      string `form:"name"`
		Title     string `form:"title"`
		Photo     string `form:"photo"`
		Telephone string `form:"telephone"`
		Email     string `form:"email"`
		Location  string `form:"location"`
		Website   string `form:"website"`
		IsActive  string `form:"is_active"`
	}
	_ = c.BodyParser(&req)

	if req.Name == "" {
		card, _ := h.cardService.GetCardByID(cardID)
		return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
			"Title":                    "Kart Düzenle",
			renderer.FlashErrorKeyView: "Ad alanı zorunludur.",
			renderer.FormDataKey:       req,
			"Card":                     card,
		}, http.StatusBadRequest)
	}

	cardData := &models.Card{
		Name:      req.Name,
		Title:     req.Title,
		Photo:     req.Photo,
		Telephone: req.Telephone,
		Email:     req.Email,
		Location:  req.Location,
		Website:   req.Website,
		IsActive:  req.IsActive == "true",
	}

	if err := h.cardService.UpdateCard(c.UserContext(), cardID, cardData); err != nil {
		card, _ := h.cardService.GetCardByID(cardID)
		return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
			"Title":                    "Kart Düzenle",
			renderer.FlashErrorKeyView: "Güncelleme hatası: " + err.Error(),
			renderer.FormDataKey:       req,
			"Card":                     card,
		}, http.StatusInternalServerError)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla güncellendi.")
	return c.Redirect("/dashboard/cards", fiber.StatusFound)
}

func (h *CardHandler) DeleteCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	cardID := uint(id)

	if err := h.cardService.DeleteCard(c.UserContext(), cardID); err != nil {
		errMsg := "Kart silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/cards", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Kart başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla silindi.")
	return c.Redirect("/dashboard/cards", fiber.StatusFound)
}

func renderCardFormError(title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
