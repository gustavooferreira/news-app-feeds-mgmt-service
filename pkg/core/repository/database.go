package repository

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	conn *gorm.DB
}

func NewDatabase(host string, port int, username string, password string, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	dbconn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// dbconn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	// create session
	dbconn = dbconn.Session(&gorm.Session{})
	dbconn = dbconn.Debug()

	// TODO: Setup logger for gorm here

	db := Database{conn: dbconn}

	return &db, nil
}

func (db *Database) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (db *Database) HealthCheck() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) FindAllFeedRecords(provider string, category string, enabled bool) ([]Feed, error) {
	var feedResults []Feed
	chain := db.conn.Joins("Provider").Joins("Category")

	if provider != "" {
		chain = chain.Where("`Provider`.`name` = ?", provider)
	}

	if category != "" {
		chain = chain.Where("`Category`.`name` = ?", category)
	}

	chain = chain.Where(&Feed{Enabled: &enabled})
	result := chain.Find(&feedResults)
	return feedResults, result.Error
}

func (db *Database) InsertFeedRecord(url string, provider string, category string, enabled bool) error {
	// Add Provider if it doesn't exist
	var providerRecord Provider
	result := db.conn.Where(Provider{Name: provider}).FirstOrCreate(&providerRecord)
	if result.Error != nil {
		return result.Error
	}

	// Add Category if it doesn't exist
	var categoryRecord Category
	result = db.conn.Where(Category{Name: category}).FirstOrCreate(&categoryRecord)
	if result.Error != nil {
		return result.Error
	}

	feedRecord := Feed{
		URL:        url,
		ProviderID: providerRecord.ID,
		CategoryID: categoryRecord.ID,
		Enabled:    &enabled,
	}

	result = db.conn.Create(&feedRecord)
	return result.Error
}

func (db *Database) UpdateFeedState(url string, enabled bool) error {
	// First check record exists
	var feedRecord Feed
	result := db.conn.Where(&Feed{URL: url}).Take(&feedRecord)
	if result.Error != nil {
		return result.Error
	}

	result = db.conn.Model(&Feed{URL: url}).Update("enabled", enabled)
	return result.Error
}

func (db *Database) DeleteFeedState(url string) error {
	// First check record exists
	var feedRecord Feed
	result := db.conn.Where(&Feed{URL: url}).Take(&feedRecord)
	if result.Error != nil {
		return result.Error
	}

	result = db.conn.Where(&Feed{URL: url}).Delete(&feedRecord)
	return result.Error
}
