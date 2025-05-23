package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"davet.link/configs/logconfig"
	"davet.link/configs/sessionconfig"
	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PanelInvitationHandler struct {
	invitationService services.IInvitationService
}

func NewPanelInvitationHandler() *PanelInvitationHandler {
	return &PanelInvitationHandler{invitationService: services.NewInvitationService()}
}

func (h *PanelInvitationHandler) ListPanelInvitations(c *fiber.Ctx) error {
	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil {
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	invitations, err := h.invitationService.GetInvitationsByUserID(userID)
	if err != nil {
		logconfig.Log.Error("Davetiye listesi alınamadı", zap.Error(err))
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiyeler alınamadı.")
		return renderer.Render(c, "panel/invitations/list", "layouts/panel", fiber.Map{
			"Title":       "Davetiyelerim",
			"Invitations": []models.Invitation{},
		}, http.StatusOK)
	}

	return renderer.Render(c, "panel/invitations/list", "layouts/panel", fiber.Map{
		"Title":       "Davetiyelerim",
		"Invitations": invitations,
	}, http.StatusOK)
}

func (h *PanelInvitationHandler) ShowCreatePanelInvitation(c *fiber.Ctx) error {
	return renderer.Render(c, "panel/invitations/create", "layouts/panel", fiber.Map{
		"Title": "Yeni Davetiye Oluştur",
	})
}

func (h *PanelInvitationHandler) CreatePanelInvitation(c *fiber.Ctx) error {
	var req struct {
		Title       string `form:"title"`
		CategoryID  uint   `form:"category_id"`
		Template    string `form:"template"`
		Type        string `form:"type"`
		Image       string `form:"image"`
		Description string `form:"description"`
		Venue       string `form:"venue"`
		Address     string `form:"address"`
	}
	if err := c.BodyParser(&req); err != nil {
		logconfig.Log.Error("Davetiye formu parse edilemedi", zap.Error(err))
		return renderPanelInvitationFormError("Yeni Davetiye Oluştur", req, "Form bilgileri alınamadı.", c)
	}

	if req.Title == "" || req.CategoryID == 0 || req.Template == "" {
		return renderPanelInvitationFormError("Yeni Davetiye Oluştur", req, "Başlık, Kategori ve Şablon alanları zorunludur.", c)
	}

	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil {
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	invitation := &models.Invitation{
		Title:       req.Title,
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Template:    req.Template,
		Type:        req.Type,
		Image:       req.Image,
		Description: req.Description,
		Venue:       req.Venue,
		Address:     req.Address,
	}

	if err := h.invitationService.CreateInvitation(c.UserContext(), invitation); err != nil {
		logconfig.Log.Error("Davetiye oluşturulamadı", zap.Error(err))
		return renderPanelInvitationFormError("Yeni Davetiye Oluştur", req, "Davetiye oluşturulamadı: "+err.Error(), c)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla oluşturuldu.")
	return c.Redirect("/panel/invitations", fiber.StatusFound)
}

func (h *PanelInvitationHandler) ShowUpdatePanelInvitation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	invitationID := uint(id)

	invitation, err := h.invitationService.GetInvitationByID(invitationID)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/panel/invitations", fiber.StatusSeeOther)
	}

	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil {
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	if invitation.UserID != userID {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Bu davetiyeyi düzenleme yetkiniz yok.")
		return c.Redirect("/panel/invitations", fiber.StatusSeeOther)
	}

	return renderer.Render(c, "panel/invitations/update", "layouts/panel", fiber.Map{
		"Title":      "Davetiyeyi Düzenle",
		"Invitation": invitation,
	})
}

func (h *PanelInvitationHandler) UpdatePanelInvitation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	invitationID := uint(id)

	invitation, err := h.invitationService.GetInvitationByID(invitationID)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/panel/invitations", fiber.StatusSeeOther)
	}

	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil {
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	if invitation.UserID != userID {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Bu davetiyeyi düzenleme yetkiniz yok.")
		return c.Redirect("/panel/invitations", fiber.StatusSeeOther)
	}

	var req struct {
		Title       string `form:"title"`
		Template    string `form:"template"`
		Type        string `form:"type"`
		Image       string `form:"image"`
		Description string `form:"description"`
		Venue       string `form:"venue"`
		Address     string `form:"address"`
	}
	if err := c.BodyParser(&req); err != nil {
		logconfig.Log.Error("Davetiye formu parse edilemedi", zap.Error(err))
		return renderer.Render(c, "panel/invitations/update", "layouts/panel", fiber.Map{
			"Title":                    "Davetiyeyi Düzenle",
			renderer.FlashErrorKeyView: "Form bilgileri alınamadı.",
			"Invitation":               invitation,
		}, http.StatusBadRequest)
	}

	if req.Title == "" || req.Template == "" {
		return renderer.Render(c, "panel/invitations/update", "layouts/panel", fiber.Map{
			"Title":                    "Davetiyeyi Düzenle",
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
		logconfig.Log.Error("Davetiye güncellenemedi", zap.Error(err))
		return renderer.Render(c, "panel/invitations/update", "layouts/panel", fiber.Map{
			"Title":                    "Davetiyeyi Düzenle",
			renderer.FlashErrorKeyView: "Güncelleme hatası: " + err.Error(),
			renderer.FormDataKey:       req,
			"Invitation":               invitation,
		}, http.StatusInternalServerError)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla güncellendi.")
	return c.Redirect("/panel/invitations", fiber.StatusFound)
}

func (h *PanelInvitationHandler) DeletePanelInvitation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	invitationID := uint(id)

	invitation, err := h.invitationService.GetInvitationByID(invitationID)
	if err != nil {
		errMsg := "Davetiye bulunamadı."
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/panel/invitations", fiber.StatusSeeOther)
	}

	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil {
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	if invitation.UserID != userID {
		errMsg := "Bu davetiyeyi silme yetkiniz yok."
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/panel/invitations", fiber.StatusSeeOther)
	}

	if err := h.invitationService.DeleteInvitation(c.UserContext(), invitationID); err != nil {
		errMsg := "Davetiye silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/panel/invitations", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Davetiye başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla silindi.")
	return c.Redirect("/panel/invitations", fiber.StatusFound)
}

func renderPanelInvitationFormError(title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, "panel/invitations/create", "layouts/panel", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
