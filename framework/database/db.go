package database

import (
	"encoder/domain"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DbType        string
	DsnTest       string
	DbTestType    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	instance := NewDb()

	instance.Env = "test"
	instance.DbTestType = "sqlite3"
	instance.DsnTest = ":memory:"
	instance.AutoMigrateDb = true
	instance.Debug = true
	instance.Env = "test"

	connection, err := instance.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (db *Database) Connect() (*gorm.DB, error) {
	var err error

	if db.Env == "test" {
		db.Db, err = gorm.Open(db.DbTestType, db.DsnTest)
	} else {
		db.Db, err = gorm.Open(db.DbType, db.Dsn)
	}

	if err != nil {
		return nil, err
	}

	if db.Debug {
		db.Db.LogMode(true)
	}

	if db.AutoMigrateDb {
		db.Db.AutoMigrate(&domain.Video{}, &domain.Job{})

		db.Db.Model(domain.Job{}).AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return db.Db, nil
}
