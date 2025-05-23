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

type InvitationHandler struct {
	invitationService services.IInvitationService
}

func NewInvitationHandler() *InvitationHandler {
	return &InvitationHandler{invitationService: services.NewInvitationService()}
}

func (h *InvitationHandler) ListInvitations(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Davetiye listesi: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
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

	paginatedResult, dbErr := h.invitationService.GetAllInvitations(params)

	renderData := fiber.Map{
		"Title":  "Davetiyeler",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		logconfig.Log.Error("Davetiye listesi DB Hatası", zap.Error(dbErr))
		renderData[renderer.FlashErrorKeyView] = "Davetiyeler getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.Invitation{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/invitations/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *InvitationHandler) ShowCreateInvitation(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Davetiye Ekle",
	})
}

func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	var req struct {
		Title         string `form:"title"`
		InvitationKey string `form:"invitation_key"`
		UserID        uint   `form:"user_id"`
		CategoryID    uint   `form:"category_id"`
		Template      string `form:"template"`
		Type          string `form:"type"`
		Image         string `form:"image"`
		Description   string `form:"description"`
		Venue         string `form:"venue"`
		Address       string `form:"address"`
	}
	_ = c.BodyParser(&req)

	if req.Title == "" || req.UserID == 0 || req.CategoryID == 0 || req.Template == "" {
		return renderInvitationFormError("Yeni Davetiye Ekle", req, "Başlık, Kullanıcı, Kategori ve Şablon alanları zorunludur.", c)
	}

	invitation := &models.Invitation{
		Title:         req.Title,
		InvitationKey: req.InvitationKey,
		UserID:        req.UserID,
		CategoryID:    req.CategoryID,
		Template:      req.Template,
		Type:          req.Type,
		Image:         req.Image,
		Description:   req.Description,
		Venue:         req.Venue,
		Address:       req.Address,
	}

	if err := h.invitationService.CreateInvitation(c.UserContext(), invitation); err != nil {
		return renderInvitationFormError("Yeni Davetiye Ekle", req, "Davetiye oluşturulamadı: "+err.Error(), c)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/invitations", fiber.StatusFound)
}

func (h *InvitationHandler) ShowUpdateInvitation(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	invitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/dashboard/invitations", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
		"Title":      "Davetiye Düzenle",
		"Invitation": invitation,
	})
}

func (h *InvitationHandler) UpdateInvitation(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	invitationID := uint(id)

	var req struct {
		Title       string `form:"title"`
		Template    string `form:"template"`
		Type        string `form:"type"`
		Image       string `form:"image"`
		Description string `form:"description"`
		Venue       string `form:"venue"`
		Address     string `form:"address"`
	}
	_ = c.BodyParser(&req)

	if req.Title == "" || req.Template == "" {
		invitation, _ := h.invitationService.GetInvitationByID(invitationID)
		return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
			"Title":                    "Davetiye Düzenle",
			renderer.FlashErrorKeyView: "Başlık ve Şablon alanları zorunludur.",
			renderer.FormDataKey:       req,
			"Invitation":               invitation,
		}, http.StatusBadRequest)
	}

	invitationData := &models.Invitation{
		Title:       req.Title,
		Template:    req.Template,
		Type:        req.Type,
		Image:       req.Image,
		Description: req.Description,
		Venue:       req.Venue,
		Address:     req.Address,
	}

	if err := h.invitationService.UpdateInvitation(c.UserContext(), invitationID, invitationData); err != nil {
		invitation, _ := h.invitationService.GetInvitationByID(invitationID)
		return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
			"Title":                    "Davetiye Düzenle",
			renderer.FlashErrorKeyView: "Güncelleme hatası: " + err.Error(),
			renderer.FormDataKey:       req,
			"Invitation":               invitation,
		}, http.StatusInternalServerError)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla güncellendi.")
	return c.Redirect("/dashboard/invitations", fiber.StatusFound)
}

func (h *InvitationHandler) DeleteInvitation(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	invitationID := uint(id)

	if err := h.invitationService.DeleteInvitation(c.UserContext(), invitationID); err != nil {
		errMsg := "Davetiye silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/invitations", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Davetiye başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla silindi.")
	return c.Redirect("/dashboard/invitations", fiber.StatusFound)
}

func renderInvitationFormError(title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
