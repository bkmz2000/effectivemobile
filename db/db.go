package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SongDB struct {
	*gorm.DB
}

func Init() *SongDB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	url := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	var db *gorm.DB
	var err error
	retries := 5
	for i := 0; i < retries; i++ {
		log.Println("Connecting to database...")
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
		if err == nil {
			log.Println("Database connection established successfully.")
			break
		}
		log.Printf("Error connecting to db: %v. Retrying in 2 seconds...\n", err)
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if err != nil {
		log.Fatal("Failed to connect to the database after multiple attempts: ", err)
	}

	if err := db.AutoMigrate(&Song{}); err != nil {
		log.Fatal("Error during database migration: ", err)
	}
	log.Println("Database migration completed successfully.")

	return &SongDB{db}
}

func (db *SongDB) Filter(group, name string, dateAfter, dateBefore string) ([]Song, error) {
	var afterTime, beforeTime time.Time
	var err error

	log.Printf("Filtering songs with group: %s, name: %s, dateAfter: %s, dateBefore: %s\n", group, name, dateAfter, dateBefore)

	if dateAfter != "" {
		afterTime, err = time.Parse("2006-01-02", dateAfter)
		if err != nil {
			log.Println("Error parsing dateAfter:", err)
			return nil, err
		}
	}

	if dateBefore != "" {
		beforeTime, err = time.Parse("2006-01-02", dateBefore)
		if err != nil {
			log.Println("Error parsing dateBefore:", err)
			return nil, err
		}
	}

	query := db.Model(&Song{})

	if group != "" {
		query = query.Where(`"group" = ?`, group)
		log.Printf("Filtering by group: %s\n", group)
	}
	if name != "" {
		query = query.Where(`"name" = ?`, name)
		log.Printf("Filtering by name: %s\n", name)
	}
	if !afterTime.IsZero() {
		query = query.Where("release_date >= ?", afterTime)
		log.Printf("Filtering by release date after: %s\n", afterTime)
	}
	if !beforeTime.IsZero() {
		query = query.Where("release_date <= ?", beforeTime)
		log.Printf("Filtering by release date before: %s\n", beforeTime)
	}

	var songs []Song
	err = query.Find(&songs).Error
	if err != nil {
		log.Println("Error retrieving songs:", err)
		return nil, err
	}

	log.Printf("Retrieved %d songs\n", len(songs))
	return songs, nil
}

func (db *SongDB) GetSong(group, name string) *Song {
	log.Printf("Getting song with group: %s and name: %s\n", group, name)
	ret := &Song{}
	db.Where(`"group" = ? AND "name" = ?`, group, name).First(&ret)

	if ret.Id == 0 {
		log.Printf("No song found with group: %s and name: %s\n", group, name)
	} else {
		log.Printf("Found song: %+v\n", ret)
	}

	return ret
}

func (db *SongDB) Remove(id uint) error {
	log.Printf("Removing song with ID: %d\n", id)
	return db.Delete(&Song{}, id).Error
}

func (db *SongDB) Update(song *Song) error {
	log.Printf("Updating song: %+v\n", song)
	return db.Save(song).Error
}

func (db *SongDB) Add(song *Song) error {
	log.Printf("Adding new song: %+v\n", song)
	return db.Create(song).Error
}
