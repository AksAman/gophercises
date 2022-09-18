package dblib

import (
	"errors"
	"fmt"

	"github.com/AksAman/gophercises/phone/models"
	"github.com/AksAman/gophercises/phone/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func (g *GormDB) Close() error {
	logger.Debug("Closing gorm db")
	sqlDB, _ := g.db.DB()
	return sqlDB.Close()
}

func (g *GormDB) Seed(data []string) error {
	utils.Title("Seeding GormDB with sample data")
	for _, number := range data {
		id, err := g.InsertPhone(number)
		if err != nil {
			return err
		}
		logger.Infof("Successfully inserted phone number %s with id %d", number, id)
	}
	return nil
}

func (g *GormDB) Migrate() error {
	utils.Title("Migrating GormDB tables")

	g.db.AutoMigrate(&models.PhoneGorm{})

	logger.Info("Successfully migrated rawdb")
	return nil
}

func (g *GormDB) InsertPhone(number string) (id int, err error) {
	id = -1

	phone := models.PhoneGorm{Phone: models.Phone{Number: number}}

	result := g.db.Create(&phone)

	if result.Error != nil {
		return id, result.Error
	}

	return int(phone.ID), result.Error
}

func (g *GormDB) All() (phones []models.PhoneGorm, err error) {
	utils.Title("Getting all phone numbers")
	tx := g.db.Order("number asc").Find(&phones)
	if tx.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &NoRecordFoundError{}
		}
		return nil, tx.Error
	}
	return phones, err
}

func (g *GormDB) Get(id int) (phone *models.PhoneGorm, err error) {
	utils.Titlef("Getting phone number with id %d", id)
	phone = &models.PhoneGorm{}

	tx := g.db.First(&phone, id)

	if tx.Error != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &NoRecordFoundError{}
		}
		return nil, tx.Error
	}

	return phone, err
}

func (g *GormDB) FindPhone(number string) (phone *models.PhoneGorm, err error) {
	utils.Titlef("Searching for phone number %q", number)
	phone = &models.PhoneGorm{}

	err = g.db.Where(&models.PhoneGorm{Phone: models.Phone{Number: number}}).First(phone).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &NoRecordFoundError{}
		}
		return nil, err
	}

	return phone, err
}

func (g *GormDB) FindPhones(number string) (phones []models.PhoneGorm, err error) {
	// utils.Titlef("FindPhones all phone numbers with number=%s", number)
	utils.Titlef("Searching for phone number %q", number)

	tx := g.db.Where(&models.PhoneGorm{Phone: models.Phone{Number: number}}).Find(&phones)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return phones, err
}
func (g *GormDB) UpdatePhone(phone *models.PhoneGorm) error {
	tx := g.db.Model(phone).Updates(phone)
	return tx.Error
}

func (g *GormDB) DeletePhone(id int) error {
	tx := g.db.Unscoped().Delete(&models.PhoneGorm{}, id)
	return tx.Error
}

// INITIALIZATION
func InitGormDB(reset bool) (IPhoneDB[models.PhoneGorm], error) {
	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	if reset {
		err = ResetGormDB("postgres", config.GetPGConnectionString(), config.DBDatabaseNameGorm)
		if err != nil {
			return nil, err
		}
	}

	connectionString := config.GetDBConnectionString(utils.GormDB)

	db, err := openGormDB("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Migrate()
	if err != nil {
		err = db.Close()
		return nil, err
	}

	return db, err
}

func ResetGormDB(driverName, connectionString, dbName string) error {
	utils.Title("ResetGormDB")
	// db, err := sqlx.Connect(driverName, connectionString)
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DriverName: driverName,
				DSN:        connectionString,
			},
		),
	)
	if err != nil {
		return err
	}

	if err := disconnectAllUsersFromGormDB(db, dbName); err != nil {
		return err
	}

	if err := resetGormDB(db, dbName); err != nil {
		return err
	}
	return nil
}

func openGormDB(driverName, dataSourceName string) (IPhoneDB[models.PhoneGorm], error) {
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DriverName: driverName,
				DSN:        dataSourceName,
			},
		),
	)
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully connected to gormdb")

	return &GormDB{db: db}, nil
}

func disconnectAllUsersFromGormDB(db *gorm.DB, dbName string) error {
	statement := fmt.Sprintf(
		`SELECT pg_terminate_backend(pg_stat_activity.pid)
	FROM pg_stat_activity
	WHERE pg_stat_activity.datname = '%s' 
		AND pid <> pg_backend_pid();`,
		dbName,
	)
	err := db.Exec(statement).Error
	if err != nil {
		return err
	}
	logger.Debug("Successfully disconnected all users from rawdb")
	return nil
}

func resetGormDB(db *gorm.DB, dbName string) error {
	utils.Title("resetGormDB")
	statement := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	err := db.Exec(statement).Error
	if err != nil {
		return err
	}
	return createGormDB(db, dbName)
}

func createGormDB(db *gorm.DB, dbName string) error {
	utils.Title("creating rawdb")
	statement := fmt.Sprintf("CREATE DATABASE %s", dbName)
	err := db.Exec(statement).Error
	if err != nil {
		return err
	}
	return nil
}
