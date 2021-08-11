package repositories

import "github.com/wildanpurnomo/abw-rematch/models"

func (p *Repository) CreateNewRedirection(redirection *models.Redirection) error {
	return p.db.Create(&redirection).Error
}

func (p *Repository) GetNewRedirection(oldSlug string) string {
	var redirection models.Redirection
	if err := p.db.Where("old = ?", oldSlug).First(&redirection).Error; err != nil {
		return ""
	} else {
		return redirection.New
	}
}
