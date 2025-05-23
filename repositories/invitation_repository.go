package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"

	"gorm.io/gorm"
)

type IInvitationRepository interface {
	GetAllInvitations(params queryparams.ListParams) ([]models.Invitation, int64, error)
	GetInvitationByID(id uint) (*models.Invitation, error)
	GetInvitationByKey(key string) (*models.Invitation, error)
	GetInvitationsByUserID(userID uint) ([]models.Invitation, error)
	CreateInvitation(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitation(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	DeleteInvitation(ctx context.Context, id uint) error
	GetInvitationCount() (int64, error)

	// Detail related methods
	CreateDetail(ctx context.Context, detail *models.InvitationDetail) error
	UpdateDetail(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	
	// Participant related methods
	AddParticipant(ctx context.Context, participant *models.InvitationParticipant) error
	UpdateParticipant(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	DeleteParticipant(ctx context.Context, id uint) error
	GetParticipantCount(invitationID uint) (int64, error)
}

type InvitationRepository struct {
	base IBaseRepository[models.Invitation]
	db   *gorm.DB
}

func NewInvitationRepository() IInvitationRepository {
	base := NewBaseRepository[models.Invitation](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "title", "created_at", "date", "is_confirmed"})
	base.SetPreloads("Category", "InvitationDetail", "Participants")

	return &InvitationRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *InvitationRepository) GetAllInvitations(params queryparams.ListParams) ([]models.Invitation, int64, error) {
	return r.base.GetAll(params)
}

func (r *InvitationRepository) GetInvitationByID(id uint) (*models.Invitation, error) {
	return r.base.GetByID(id)
}

func (r *InvitationRepository) GetInvitationByKey(key string) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.Preload("Category").
		Preload("InvitationDetail").
		Preload("Participants").
		Where("invitation_key = ?", key).
		First(&invitation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return &invitation, err
}

func (r *InvitationRepository) GetInvitationsByUserID(userID uint) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Find(&invitations).Error
	return invitations, err
}

func (r *InvitationRepository) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
	return r.base.CreateWithRelations(ctx, invitation)
}

func (r *InvitationRepository) UpdateInvitation(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.base.Update(ctx, id, data, updatedBy)
}

func (r *InvitationRepository) DeleteInvitation(ctx context.Context, id uint) error {
	return r.base.DeleteWithRelations(ctx, id)
}

func (r *InvitationRepository) GetInvitationCount() (int64, error) {
	return r.base.GetCount()
}

// Detail related methods
func (r *InvitationRepository) CreateDetail(ctx context.Context, detail *models.InvitationDetail) error {
	return r.db.Create(detail).Error
}

func (r *InvitationRepository) UpdateDetail(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.db.Model(&models.InvitationDetail{}).Where("id = ?", id).Updates(data).Error
}

// Participant related methods
func (r *InvitationRepository) AddParticipant(ctx context.Context, participant *models.InvitationParticipant) error {
	return r.db.Create(participant).Error
}

func (r *InvitationRepository) UpdateParticipant(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.db.Model(&models.InvitationParticipant{}).Where("id = ?", id).Updates(data).Error
}

func (r *InvitationRepository) DeleteParticipant(ctx context.Context, id uint) error {
	return r.db.Delete(&models.InvitationParticipant{}, id).Error
}

func (r *InvitationRepository) GetParticipantCount(invitationID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.InvitationParticipant{}).
		Where("invitation_id = ?", invitationID).
		Count(&count).Error
	return count, err
}

var _ IInvitationRepository = (*InvitationRepository)(nil)
var _ IBaseRepository[models.Invitation] = (*BaseRepository[models.Invitation])(nil)
