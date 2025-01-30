package claimDao

import (
	"fmt"
	"strings"

	"github.com/abdulmalikraji/verify-influencers-backend/db/connection"
	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"gorm.io/gorm"
)

type DataAccess interface {
	// Postgres Data Access Object Methods
	FindAll() ([]models.Claim, error)
	FindAllByInfluencerId(influencerId string) ([]models.Claim, error)
	FindById(id int) (models.Claim, error)
	Insert(item models.Claim) error
	Update(item models.Claim) error
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

func (d dataAccess) FindAll() ([]models.Claim, error) {

	var claims []models.Claim
	result := d.db.Table(models.Claim{}.TableName()).Where("del_flg = ?", false).Find(&claims)
	if result.Error != nil {
		return []models.Claim{}, result.Error
	}
	return claims, nil
}

func (d dataAccess) FindAllByInfluencerId(influencerId string) ([]models.Claim, error) {

	var claims []models.Claim
	result := d.db.Table(models.Claim{}.TableName()).Where("influencer_id = ? AND del_flg = ?", influencerId, false).Find(&claims)
	if result.Error != nil {
		return []models.Claim{}, result.Error
	}

	return claims, nil
}

func (d dataAccess) FindById(id int) (models.Claim, error) {

	var claim models.Claim
	result := d.db.Table(models.Claim{}.TableName()).Where("id = ? AND del_flg = ?", id, false).First(&claim)
	if result.Error != nil {
		return models.Claim{}, result.Error
	}
	return claim, nil
}

// func (d dataAccess) Insert(item models.Claim) error {

// 	result := d.db.Table(item.TableName()).Create(&item)

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

func (d dataAccess) Insert(item models.Claim) error {
	// Normalize the claim text (case-insensitive comparison)
	normalizedClaim := strings.ToLower(strings.TrimSpace(item.ParsedClaim))

	// Check if a claim with the same text and influencer_id already exists
	var existingClaim models.Claim
	err := d.db.Table(item.TableName()).
		Where("LOWER(parsed_claim) = ? AND influencer_id = ?", normalizedClaim, item.InfluencerID).
		First(&existingClaim).Error

	if err == nil {
		// Claim already exists for the same influencer, skip insertion
		fmt.Println("Claim already exists for ifluencer")
		return nil
	} else if err != gorm.ErrRecordNotFound {
		// Return any other database error
		return err
	}

	// Insert the new claim if it doesn't exist
	result := d.db.Table(item.TableName()).Create(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) Update(item models.Claim) error {

	result := d.db.Table(item.TableName()).Where("id = ? ", item.ID).Updates(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) SoftDelete(id int) error {

	var item models.Claim

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Update("del_flg", true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) Delete(id int) error {

	var item models.Claim

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Delete(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
