package influencerDao

import (
	"github.com/abdulmalikraji/verify-influencers-backend/db/connection"
	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"gorm.io/gorm"
)

type DataAccess interface {
	// Postgres Data Access Object Methods
	FindAll() ([]models.Influencer, error)
	FindById(id int) (models.Influencer, error)
	FindByUsername(username string) (models.Influencer, error)		
	Insert(item models.Influencer) (models.Influencer, error)
	Update(item models.Influencer) error
	SoftDelete(id int) error
	Delete(id int) error
}

type dataAccess struct {
	db *gorm.DB
}

func New(client connection.Client) DataAccess {
	return dataAccess{
		db: client.PostgresConnection,
	}
}

func (d dataAccess) FindAll() ([]models.Influencer, error) {

	var influencers []models.Influencer
	result := d.db.Table(models.Influencer{}.TableName()).Where("del_flg = ?", false).Find(&influencers)
	if result.Error != nil {
		return []models.Influencer{}, result.Error
	}
	return influencers, nil
}

func (d dataAccess) FindById(id int) (models.Influencer, error) {

	var influencer models.Influencer
	result := d.db.Table(models.Influencer{}.TableName()).Where("influencer_id = ? AND del_flg = ?", id, false).First(&influencer)
	if result.Error != nil {
		return models.Influencer{}, result.Error
	}
	return influencer, nil
}

func (d dataAccess) FindByUsername(username string) (models.Influencer, error) {

	var influencer models.Influencer
	result := d.db.Table(models.Influencer{}.TableName()).Where("username = ? AND del_flg = ?", username, false).First(&influencer)
	if result.Error != nil {
		return models.Influencer{}, result.Error
	}
	return influencer, nil
}

func (d dataAccess) Insert(item models.Influencer) (models.Influencer, error) {

	result := d.db.Table(item.TableName()).Create(&item)

	if result.Error != nil {
		return models.Influencer{}, result.Error
	}

	return item, nil
}

func (d dataAccess) Update(item models.Influencer) error {

	result := d.db.Table(item.TableName()).Where("influencer_id = ? ", item.ID).Updates(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) SoftDelete(id int) error {

	var item models.Influencer

	result := d.db.Table(item.TableName()).Where("influencer_id = ? ", id).Update("del_flg", true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) Delete(id int) error {

	var item models.Influencer

	result := d.db.Table(item.TableName()).Where("influencer_id = ? ", id).Delete(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
