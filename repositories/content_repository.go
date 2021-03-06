package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/wildanpurnomo/abw-rematch/models"
)

func (p *Repository) GetContentById(content *models.Content, contentId string) error {
	return p.db.Where("id = ?", contentId).First(&content).Error
}

func (p *Repository) GetContentBySlug(content *models.Content, slug string) error {
	return p.db.Where("slug = ?", slug).First(&content).Error
}

func (p *Repository) GetContentByUserId(contents *[]models.Content, userId string, limit int, offset int) error {
	return p.db.Where("user_id = ?", userId).Limit(limit).Offset(offset).Find(&contents).Error
}

func (p *Repository) GetContentInUserIds(contents *[]models.Content, userIds []string) error {
	return p.db.Where("user_id IN (?)", userIds).Find(&contents).Error
}

func (p *Repository) GetContentByUserIdAndContentId(content *models.Content, userId string, contentId string) error {
	return p.db.Where("user_id = ? AND id = ?", userId, contentId).First(&content).Error
}

func (p *Repository) GetContentByUserIdAndTitle(content *models.Content, userId string, title string) *gorm.DB {
	return p.db.Where("user_id = ? AND title = ?", userId, title).First(&content)
}

func (p *Repository) CreateNewContent(content *models.Content) error {
	return p.db.Create(&content).Error
}

func (p *Repository) UpdateContent(content *models.Content, update models.Content) error {
	return p.db.Model(&content).Updates(update).Error
}

func (p *Repository) DeleteContent(content *models.Content) error {
	return p.db.Model(&content).Delete(&content).Error
}
