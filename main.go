package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var connString = os.Getenv("CONNECTION_STRING")

func main() {
	mode := flag.String("mode", "", "")
	firstName := flag.String("firstname", "", "")
	lastName := flag.String("lastname", "", "")
	albumname := flag.String("albumname", "", "")
	singerid := flag.String("singerid", "", "")

	flag.Parse()

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		// DisableNestedTransaction will turn off the use of Savepoints if gorm
		// detects a nested transaction. Cloud Spanner does not support Savepoints,
		// so it is recommended to set this configuration option to true.
		DisableNestedTransaction: true,
		Logger:                   logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		fmt.Printf("Failed to open gorm connection: %v\n", err)
		return
	}

	switch *mode {
	case "createsinger":
		if *firstName != "" && *lastName != "" {
			_, err := CreateSinger(db, *firstName, *lastName)
			if err != nil {
				log.Println(err)
				return
			}
		}
	case "createalbum":
		if *singerid != "" && *albumname != "" {
			_, err := CreateAlbum(db, *singerid, *albumname, 3)
			if err != nil {
				log.Println(err)
				return
			}
		}
	case "listsingers":
		singers, err := ListSingers(db)
		if err != nil {
			log.Println(err)
			return
		}
		for n, singer := range singers {
			fmt.Println(n+1, singer.ID, singer.FullName, singer.Albums, singer.UpdatedAt)
		}
	default:
		help()
		return
	}
}

func CreateSinger(db *gorm.DB, firstName, lastName string) (string, error) {
	singer := Singer{
		BaseModel: BaseModel{ID: uuid.NewString()},
		FirstName: sql.NullString{String: firstName, Valid: true},
		LastName:  lastName,
	}
	log.Printf("%+v\n", singer)
	res := db.Create(&singer)
	log.Printf("%+v\n", res)
	if singer.FullName != firstName+" "+lastName {
		return "", fmt.Errorf("unexpected full name for singer: %v", singer.FullName)
	}

	return singer.ID, res.Error
}

func ListSingers(db *gorm.DB) ([]*Singer, error) {
	var singers []*Singer
	res := db.Find(&singers)
	if res.Error != nil {
		log.Println(res.Error)
		return []*Singer{}, res.Error
	}
	return singers, nil
}

func CreateAlbum(db *gorm.DB, singerId, albumTitle string, numTracks int) (string, error) {
	albumId := uuid.NewString()
	// We cannot include the Tracks that we want to create in the definition here, as gorm would then try to
	// use an UPSERT to save-or-update the album that we are creating. Instead, we need to create the album first,
	// and then create the tracks.
	res := db.Create(&Album{
		BaseModel: BaseModel{ID: albumId},
		Title:     albumTitle,
		// MarketingBudget: decimal.NullDecimal{Decimal: decimal.NewFromFloat(randFloat64(0, 10000000))},
		ReleaseDate: datatypes.Date(time.Now()),
		SingerId:    singerId,
		// CoverPicture:    randBytes(randInt(5000, 15000)),
	})
	if res.Error != nil {
		return albumId, res.Error
	}
	tracks := make([]*Track, numTracks)
	for n := 0; n < numTracks; n++ {
		randTrackTitle := fmt.Sprintf("track-%d", n)
		tracks[n] = &Track{BaseModel: BaseModel{ID: albumId}, TrackNumber: int64(n + 1), Title: randTrackTitle}
	}

	// Note: The batch size is deliberately kept small here in order to prevent the statement from getting too big and
	// exceeding the maximum number of parameters in a prepared statement. PGAdapter can currently handle at most 50
	// parameters in a prepared statement.
	res = db.CreateInBatches(tracks, 8)
	log.Printf("%+v\n", tracks)

	return albumId, res.Error
}

func help() {
	fmt.Println("Nothing to do")
}
