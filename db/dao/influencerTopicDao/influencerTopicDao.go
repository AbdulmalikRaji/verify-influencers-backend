package influencerTopicDao

import (
	"github.com/abdulmalikraji/verify-influencers-backend/db/connection"
	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"gorm.io/gorm"
)

type DataAccess interface {
	// Postgres Data Access Object Methods
	FindAll() ([]models.InfluencerTopic, error)
	FindById(id int) (models.InfluencerTopic, error)
	FindAllByInfluencerId(influencerId int) ([]models.InfluencerTopic, error)
	Insert(item models.InfluencerTopic) (models.InfluencerTopic, error)
	Update(item models.InfluencerTopic) error
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

func (d dataAccess) FindAll() ([]models.InfluencerTopic, error) {

	var influencerTopics []models.InfluencerTopic
	result := d.db.Table(models.InfluencerTopic{}.TableName()).Where("del_flg = ?", false).Find(&influencerTopics)
	if result.Error != nil {
		return []models.InfluencerTopic{}, result.Error
	}
	return influencerTopics, nil
}

func (d dataAccess) FindById(id int) (models.InfluencerTopic, error) {

	var influencerTopic models.InfluencerTopic
	result := d.db.Table(models.InfluencerTopic{}.TableName()).Where("id = ? AND del_flg = ?", id, false).First(&influencerTopic)
	if result.Error != nil {
		return models.InfluencerTopic{}, result.Error
	}
	return influencerTopic, nil
}

func (d dataAccess) FindAllByInfluencerId(influencerId int) ([]models.InfluencerTopic, error) {

	var influencerTopics []models.InfluencerTopic
	result := d.db.Table(models.InfluencerTopic{}.TableName()).Where("influencer_id = ? AND del_flg = ?", influencerId, false).Find(&influencerTopics)
	if result.Error != nil {
		return []models.InfluencerTopic{}, result.Error
	}
	return influencerTopics, nil
}

func (d dataAccess) Insert(item models.InfluencerTopic) (models.InfluencerTopic, error) {

	result := d.db.Table(item.TableName()).Create(&item)

	if result.Error != nil {
		return models.InfluencerTopic{}, result.Error
	}

	return item, nil
}

func (d dataAccess) Update(item models.InfluencerTopic) error {

	result := d.db.Table(item.TableName()).Where("id = ? ", item.ID).Updates(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) SoftDelete(id int) error {

	var item models.InfluencerTopic

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Update("del_flg", true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) Delete(id int) error {

	var item models.InfluencerTopic

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Delete(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
