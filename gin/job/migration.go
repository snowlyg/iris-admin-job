package job

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func GetMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20211215120700_create_jobs_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Job{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("jobs")
		},
	}
}
