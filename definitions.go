package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

//var connectionString = "host=/tmp port=5433 database=gorm-sample2"
var connectionString = os.Getenv("CONNECTION_STRING")

type BaseModel struct {
	ID        string `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Singer struct {
	BaseModel
	FirstName sql.NullString
	LastName  string
	FullName  string `gorm:"->;type:GENERATED ALWAYS AS (coalesce(concat(first_name,' '::varchar,last_name))) STORED;default:(-);"`
	Active    bool
	Albums    []Album
}

type Album struct {
	BaseModel
	Title           string
	MarketingBudget decimal.NullDecimal
	ReleaseDate     datatypes.Date
	CoverPicture    []byte
	SingerId        string
	Singer          Singer
	Tracks          []Track `gorm:"foreignKey:ID"`
}

type Track struct {
	BaseModel
	TrackNumber int64 `gorm:"primaryKey;autoIncrement:false"`
	Title       string
	SampleRate  float64
	Album       Album `gorm:"foreignKey:ID"`
}

type Venue struct {
	BaseModel
	Name        string
	Description string
}

type Concert struct {
	BaseModel
	Name      string
	Venue     Venue
	VenueId   string
	Singer    Singer
	SingerId  string
	StartTime time.Time
	EndTime   time.Time
}
