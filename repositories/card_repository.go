package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"

	"gorm.io/gorm"
)

type ICardRepository interface {
	GetAllCards(params queryparams.ListParams) ([]models.Card, int64, error)
	GetCardByID(id uint) (*models.Card, error)
	GetCardByUserID(userID uint) (*models.Card, error)
	CreateCard(ctx context.Context, card *models.Card) error
	UpdateCard(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	DeleteCard(ctx context.Context, id uint) error
	GetCardCount() (int64, error)
	
	// Bank related methods
	AddBank(ctx context.Context, cardBank *models.CardBank) error
	UpdateBank(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	DeleteBank(ctx context.Context, id uint) error
	
	// Social Media related methods
	AddSocialMedia(ctx context.Context, cardSocialMedia *models.CardSocialMedia) error
	UpdateSocialMedia(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	DeleteSocialMedia(ctx context.Context, id uint) error
}

type CardRepository struct {
	base IBaseRepository[models.Card]
	db   *gorm.DB
}

func NewCardRepository() ICardRepository {
	base := NewBaseRepository[models.Card](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "name", "created_at", "is_active"})
	base.SetPreloads("CardBanks", "CardBanks.Bank", "CardSocialMedia", "CardSocialMedia.SocialMedia")

	return &CardRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *CardRepository) GetAllCards(params queryparams.ListParams) ([]models.Card, int64, error) {
	return r.base.GetAll(params)
}

func (r *CardRepository) GetCardByID(id uint) (*models.Card, error) {
	return r.base.GetByID(id)
}

func (r *CardRepository) GetCardByUserID(userID uint) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("CardBanks.Bank").
		Preload("CardSocialMedia.SocialMedia").
		Where("user_id = ?", userID).
		First(&card).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	return &card, err
}

func (r *CardRepository) CreateCard(ctx context.Context, card *models.Card) error {
	return r.base.CreateWithRelations(ctx, card)
}

func (r *CardRepository) UpdateCard(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.base.Update(ctx, id, data, updatedBy)
}

func (r *CardRepository) DeleteCard(ctx context.Context, id uint) error {
	return r.base.DeleteWithRelations(ctx, id)
}

func (r *CardRepository) GetCardCount() (int64, error) {
	return r.base.GetCount()
}

// Bank related methods
func (r *CardRepository) AddBank(ctx context.Context, cardBank *models.CardBank) error {
	return r.db.Create(cardBank).Error
}

func (r *CardRepository) UpdateBank(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.db.Model(&models.CardBank{}).Where("id = ?", id).Updates(data).Error
}

func (r *CardRepository) DeleteBank(ctx context.Context, id uint) error {
	return r.db.Delete(&models.CardBank{}, id).Error
}

// Social Media related methods
func (r *CardRepository) AddSocialMedia(ctx context.Context, cardSocialMedia *models.CardSocialMedia) error {
	return r.db.Create(cardSocialMedia).Error
}

func (r *CardRepository) UpdateSocialMedia(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.db.Model(&models.CardSocialMedia{}).Where("id = ?", id).Updates(data).Error
}

func (r *CardRepository) DeleteSocialMedia(ctx context.Context, id uint) error {
	return r.db.Delete(&models.CardSocialMedia{}, id).Error
}

var _ ICardRepository = (*CardRepository)(nil)
var _ IBaseRepository[models.Card] = (*BaseRepository[models.Card])(nil)
