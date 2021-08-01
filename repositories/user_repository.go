package repositories

import "github.com/wildanpurnomo/abw-rematch/models"

func (p *Repository) CreateNewUser(user *models.User) error {
	return p.db.Create(user).Error
}
