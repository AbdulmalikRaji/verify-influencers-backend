package claimVerificationDao

import (
	"github.com/abdulmalikraji/verify-influencers-backend/db/connection"
	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"gorm.io/gorm"
)

type DataAccess interface {
	// Postgres Data Access Object Methods
	FindAll() ([]models.ClaimVerification, error)
	FindByClaimId(claimId int) (models.ClaimVerification, error)
	FindById(id int) (models.ClaimVerification, error)
	Insert(item models.ClaimVerification) (models.ClaimVerification, error)
	Update(item models.ClaimVerification) error
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

func (d dataAccess) FindAll() ([]models.ClaimVerification, error) {

	var claimVerifications []models.ClaimVerification
	result := d.db.Table(models.ClaimVerification{}.TableName()).Where("del_flg = ?", false).Find(&claimVerifications)
	if result.Error != nil {
		return []models.ClaimVerification{}, result.Error
	}
	return claimVerifications, nil
}

func (d dataAccess) FindByClaimId(claimId int) (models.ClaimVerification, error) {

	var claimVerifications models.ClaimVerification
	result := d.db.Table(models.ClaimVerification{}.TableName()).Where("claim_id = ? AND del_flg = ?", claimId, false).First(&claimVerifications)
	if result.Error != nil {
		return models.ClaimVerification{}, result.Error
	}

	return claimVerifications, nil
}

func (d dataAccess) FindById(id int) (models.ClaimVerification, error) {

	var claimVerification models.ClaimVerification
	result := d.db.Table(models.ClaimVerification{}.TableName()).Where("id = ? AND del_flg = ?", id, false).First(&claimVerification)
	if result.Error != nil {
		return models.ClaimVerification{}, result.Error
	}
	return claimVerification, nil
}

// func (d dataAccess) Insert(item models.ClaimVerification) error {

// 	result := d.db.Table(item.TableName()).Create(&item)

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

func (d dataAccess) Insert(item models.ClaimVerification) (models.ClaimVerification, error) {

	result := d.db.Table(item.TableName()).Create(&item)

	if result.Error != nil {
		return models.ClaimVerification{}, result.Error
	}

	return item, nil
}

func (d dataAccess) Update(item models.ClaimVerification) error {

	result := d.db.Table(item.TableName()).Where("id = ? ", item.ID).Updates(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) SoftDelete(id int) error {

	var item models.ClaimVerification

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Update("del_flg", true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) Delete(id int) error {

	var item models.ClaimVerification

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Delete(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
