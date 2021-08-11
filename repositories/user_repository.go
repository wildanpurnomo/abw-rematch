package repositories

import "github.com/wildanpurnomo/abw-rematch/models"

func (p *Repository) CreateNewUser(user *models.User) error {
	return p.db.Create(&user).Error
}

func (p *Repository) FetchUserById(user *models.User, id string) error {
	return p.db.Where("id = ?", id).First(&user).Error
}

func (p *Repository) FetchUserByUsername(user *models.User, username string) error {
	return p.db.Where("username = ?", username).First(&user).Error
}

func (p *Repository) UpdateUser(user *models.User, update models.User) error {
	return p.db.Model(&user).Updates(update).Error
}
