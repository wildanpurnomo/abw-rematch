package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/wildanpurnomo/abw-rematch/models"
)

func (p *Repository) GetContentByUserId(contents *[]models.Content, userId uint) error {
	return p.db.Where("user_id = ?", userId).Find(&contents).Error
}

func (p *Repository) GetContentByUserIdAndTitle(content *models.Content, userId uint, title string) *gorm.DB {
	return p.db.Where("user_id = ? AND title = ?", userId, title).First(&content)
}

func (p *Repository) CreateNewContent(user *models.User, content models.Content) error {
	return p.db.Model(&user).Association("Contents").Append(&content).Error
}
