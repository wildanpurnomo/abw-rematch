package repositories

import "github.com/wildanpurnomo/abw-rematch/models"

func (p *Repository) CreateNewRedirection(content *models.Content, redirection models.Redirection) error {
	return p.db.Model(&content).Association("Redirections").Append(redirection).Error
}
