package repository

type Feed struct {
	URL        string `gorm:"primaryKey;type:varchar(250);not null"`
	Provider   Provider
	ProviderID uint64 `gorm:"not null"` // Foreign Key
	Category   Category
	CategoryID uint64 `gorm:"not null"` // Foreign Key
	Enabled    *bool  `gorm:"not null;default:false"`
}

type Provider struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement;not null"`
	Name string `gorm:"type:varchar(30);uniqueIndex;not null"`
}

type Category struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement;not null"`
	Name string `gorm:"type:varchar(30);uniqueIndex;not null"`
}
