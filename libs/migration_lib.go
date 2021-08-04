package libs

import (
	"fmt"

	"github.com/wildanpurnomo/abw-rematch/models"
)

func main() {
	// used for operations with existing table that are not handled by AutoMigrate
	// e.g. add foreign key, add index, delete column, etc.
	// cd libs -> go run migration_lib.go
	// edit as much as you like
	db, err := models.ConnectDatabase()
	if err != nil {
		fmt.Printf("error connecting database: %v", err)
		return
	}
	db.Model(&models.Content{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	fmt.Println("done migrating")
}
