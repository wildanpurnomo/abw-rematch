package repositories

import "github.com/wildanpurnomo/abw-rematch/models"

func (p *Repository) CreateNewRedirection(redirection *models.Redirection) error {
	return p.db.Create(&redirection).Error
}
