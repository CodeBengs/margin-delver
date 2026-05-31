package lib

import (
	"path/filepath"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func RunMigrations(cfg *AppConfig, log *BaseLog, db *gorm.DB) error {
	if !cfg.DBRunMigrations {
		log.SugarLog().Info("DB migrations skipped (DB_RUN_MIGRATIONS=false)")
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.SugarLog().Errorf("failed to get sql.DB: %v", err)
		return err
	}

	migrationDir := filepath.Join("migrations")

	if err := goose.SetDialect("mysql"); err != nil {
		log.SugarLog().Errorf("failed to set goose dialect: %v", err)
		return err
	}

	log.SugarLog().Infof("running goose migrations from %s", migrationDir)
	if err := goose.Up(sqlDB, migrationDir); err != nil {
		log.SugarLog().Errorf("goose migration failed: %v", err)
		return err
	}

	log.SugarLog().Info("goose migrations completed successfully")
	return nil
}
