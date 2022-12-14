package main

import (
	"fmt"

	"github.com/AksAman/gophercises/phone/dblib"
	"github.com/AksAman/gophercises/phone/models"
	"github.com/AksAman/gophercises/phone/normalizer"
	"github.com/AksAman/gophercises/phone/utils"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	utils.InitializeLogger("")
	logger = utils.Logger
}

func must(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

var seedData = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

func main() {
	RunRawDB()
	RunSqlxDB()
	RunGormDB()

}

func RunRawDB() {

	// region Initialize DB

	utils.Title("Initialize RawDB")
	db, err := dblib.InitRawDB(true)
	must(err)

	phoneDB := db.(*dblib.RawDB)
	defer func() {
		utils.Title("Closing DB")
		err := phoneDB.Close()
		must(err)
	}()
	// endregion

	var phone *models.PhoneRaw
	// region Seed DB
	err = phoneDB.Seed(seedData)
	must(err)
	// endregion

	// region Get By ID
	id := 2
	phone, err = phoneDB.Get(id)
	must(err)
	logger.Infof("Phone for id %d: %v\n", id, phone.String())
	// endregion

	// region Search
	searchNumber := func(phoneNumberToFind string) {
		phone, err := phoneDB.FindPhone(phoneNumberToFind)
		if _, ok := err.(*dblib.NoRecordFoundError); ok {
			logger.Warnf("No record found for %s", phoneNumberToFind)
		} else {
			must(err)
		}
		if phone != nil {
			logger.Infof("Found phone: %v", phone.String())
		}
	}
	testNumbers := []string{
		"1234567890",
		"Not a phone number",
	}
	for _, testNumber := range testNumbers {
		searchNumber(testNumber)
	}

	// endregion

	// region All
	allPhones, err := phoneDB.All()
	must(err)
	for _, p := range allPhones {
		logger.Infof("phone: %#v", p.String())
	}
	// endregion

	// normalize and update phone numbers
	utils.Title("Normalize and update phone numbers")
	for _, p := range allPhones {
		normalizedNumber := normalizer.NormalizePhoneNumber(p.Number)
		if p.Number == normalizedNumber {
			logger.Infof("Phone number %s is already normalized", p.Number)
			continue
		}

		logger.Infof("normalizing %v to %s", p.String(), normalizedNumber)
		existingPhones, err := phoneDB.FindPhones(normalizedNumber)
		must(err)
		if len(existingPhones) > 0 {
			logger.Warnf("%d Phone numbers already exists with id %d and number %s", len(existingPhones), p.ID, normalizedNumber)
			for _, existingPhone := range existingPhones {
				err := phoneDB.DeletePhone(existingPhone.ID)
				if err != nil {
					logger.Errorf("Error deleting phone: %#v", err)
					continue
				}
				logger.Warnf("Deleted phone: %v", existingPhone.String())
			}
		}

		p.Number = normalizedNumber
		err = phoneDB.UpdatePhone(&p)
		must(err)
		logger.Infof("Updated phone: %v\n", p.String())
	}

	// region All
	allPhones, err = phoneDB.All()
	must(err)
	for _, p := range allPhones {
		logger.Infof("phone: %v", p.String())
	}
	// endregion
}

func RunSqlxDB() {

	// region Initialize DB
	utils.Title("Initialize SqlxDB")
	db, err := dblib.InitSqlxDB(true)
	must(err)

	phoneDB := db.(*dblib.SqlxDB)
	defer func() {
		utils.Title("Closing DB")
		err := phoneDB.Close()
		must(err)
	}()
	// endregion

	// region Seed DB
	err = phoneDB.Seed(seedData)
	must(err)
	// endregion

	var phone *models.PhoneSqlx

	// region Get By ID
	id := 2
	phone, err = phoneDB.Get(id)
	must(err)
	logger.Infof("Phone for id %d: %v\n", id, phone.String())
	// endregion

	// region Search
	searchNumber := func(phoneNumberToFind string) {
		phone, err := phoneDB.FindPhone(phoneNumberToFind)
		if _, ok := err.(*dblib.NoRecordFoundError); ok {
			logger.Warnf("No record found for %s", phoneNumberToFind)
		} else {
			must(err)
		}
		if phone != nil {
			logger.Infof("Found phone: %v", phone.String())
		}
	}
	testNumbers := []string{
		"1234567890",
		"Not a phone number",
	}
	for _, testNumber := range testNumbers {
		searchNumber(testNumber)
	}

	// endregion

	// region All
	allPhones, err := phoneDB.All()
	must(err)
	for _, p := range allPhones {
		logger.Infof("phone: %v", p.String())
	}
	// endregion

	// normalize and update phone numbers
	utils.Title("Normalize and update phone numbers")
	for _, p := range allPhones {
		normalizedNumber := normalizer.NormalizePhoneNumber(p.Number)
		if p.Number == normalizedNumber {
			logger.Infof("Phone number %s is already normalized", p.Number)
			continue
		}

		logger.Infof("normalizing %v to %s", p.String(), normalizedNumber)
		existingPhones, err := phoneDB.FindPhones(normalizedNumber)
		must(err)
		if len(existingPhones) > 0 {
			logger.Warnf("%d Phone numbers already exists with id %d and number %s", len(existingPhones), p.ID, normalizedNumber)
			for _, existingPhone := range existingPhones {
				err := phoneDB.DeletePhone(existingPhone.ID)
				if err != nil {
					logger.Errorf("Error deleting phone: %#v", err)
					continue
				}
				logger.Warnf("Deleted phone: %v", existingPhone.String())
			}
		}

		p.Number = normalizedNumber
		err = phoneDB.UpdatePhone(&p)
		must(err)
		logger.Infof("Updated phone: %#v\n", p.String())
	}

	// region All
	allPhones, err = phoneDB.All()
	must(err)
	for _, p := range allPhones {
		logger.Infof("phone: %v", p.String())
	}
	// endregion
}

func RunGormDB() {
	var phone *models.PhoneGorm

	// region Initialize DB
	utils.Title("Initialize GormDB")
	db, err := dblib.InitGormDB(true)
	must(err)
	fmt.Printf("db: %v\n", db)

	phoneDB := db.(*dblib.GormDB)
	defer func() {
		utils.Title("Closing DB")
		err := phoneDB.Close()
		must(err)
	}()
	// endregion

	// region Seed DB
	err = phoneDB.Seed(seedData)
	must(err)
	// endregion

	// region Get By ID
	id := 2
	phone, err = phoneDB.Get(id)
	must(err)
	logger.Infof("Phone for id %d: %v\n", id, phone.String())
	// endregion

	// region Search
	searchNumber := func(phoneNumberToFind string) {
		phone, err := phoneDB.FindPhone(phoneNumberToFind)
		if _, ok := err.(*dblib.NoRecordFoundError); ok {
			logger.Warnf("No record found for %s", phoneNumberToFind)
		} else {
			must(err)
		}
		if phone != nil {
			logger.Infof("Found phone: %v", phone.String())
		}
	}
	testNumbers := []string{
		"1234567890",
		"Not a phone number",
	}
	for _, testNumber := range testNumbers {
		searchNumber(testNumber)
	}

	// endregion

	// region All
	allPhones, err := phoneDB.All()
	must(err)
	for _, p := range allPhones {
		logger.Infof("phone: %v", p.String())
	}
	// endregion

	// normalize and update phone numbers
	utils.Title("Normalize and update phone numbers")
	for _, p := range allPhones {
		normalizedNumber := normalizer.NormalizePhoneNumber(p.Number)
		if p.Number == normalizedNumber {
			logger.Infof("Phone number %s is already normalized", p.Number)
			continue
		}

		logger.Infof("normalizing %#v to %s", p.String(), normalizedNumber)
		existingPhones, err := phoneDB.FindPhones(normalizedNumber)
		must(err)
		if len(existingPhones) > 0 {
			logger.Warnf("%d Phone numbers already exists with id %d and number %s", len(existingPhones), p.ID, normalizedNumber)
			for _, existingPhone := range existingPhones {
				err := phoneDB.DeletePhone(int(existingPhone.ID))
				if err != nil {
					logger.Errorf("Error deleting phone: %#v", err)
					continue
				}
				logger.Warnf("Deleted phone: %#v", existingPhone.String())
			}
		}

		p.Number = normalizedNumber
		err = phoneDB.UpdatePhone(&p)
		must(err)
		logger.Infof("Updated phone: %v\n", p.String())
	}

	// region All
	allPhones, err = phoneDB.All()
	must(err)
	for _, p := range allPhones {
		logger.Infof("phone: %v", p.String())
	}
	// endregion

}
