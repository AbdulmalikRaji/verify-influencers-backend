package topicDao

import (
	"github.com/abdulmalikraji/verify-influencers-backend/db/connection"
	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"gorm.io/gorm"
)

type DataAccess interface {
	// Postgres Data Access Object Methods
	FindAll() ([]models.Topic, error)
	FindById(id int) (models.Topic, error)
	Insert(item models.Topic) (models.Topic, error)
	Update(item models.Topic) error
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

func (d dataAccess) FindAll() ([]models.Topic, error) {

	var topics []models.Topic
	result := d.db.Table(models.Topic{}.TableName()).Where("del_flg = ?", false).Find(&topics)
	if result.Error != nil {
		return []models.Topic{}, result.Error
	}
	return topics, nil
}

func (d dataAccess) FindById(id int) (models.Topic, error) {

	var topic models.Topic
	result := d.db.Table(models.Topic{}.TableName()).Where("id = ? AND del_flg = ?", id, false).First(&topic)
	if result.Error != nil {
		return models.Topic{}, result.Error
	}
	return topic, nil
}

func (d dataAccess) Insert(item models.Topic) (models.Topic, error) {

	result := d.db.Table(item.TableName()).Create(&item)

	if result.Error != nil {
		return models.Topic{}, result.Error
	}

	return item, nil
}

func (d dataAccess) Update(item models.Topic) error {

	result := d.db.Table(item.TableName()).Where("id = ? ", item.ID).Updates(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) SoftDelete(id int) error {

	var item models.Topic

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Update("del_flg", true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d dataAccess) Delete(id int) error {

	var item models.Topic

	result := d.db.Table(item.TableName()).Where("id = ? ", id).Delete(&item)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
