package dblib

import (
	"database/sql"
	"fmt"

	"github.com/AksAman/gophercises/phone/models"
	"github.com/AksAman/gophercises/phone/utils"
)

type RawDB struct {
	db *sql.DB
}

// InsertPhone "C"RUD
func (r *RawDB) InsertPhone(number string) (id int, err error) {
	id = -1
	statement := `INSERT INTO phone_numbers(number) VALUES($1) RETURNING id`

	row := r.db.QueryRow(statement, number)
	err = row.Scan(&id)
	if err != nil {
		return
	}
	return
}

// All : C"R"UD
func (r *RawDB) All() (phones []*models.PhoneRaw, err error) {
	utils.Title("Getting all phone numbers")
	statement := `SELECT id, number FROM phone_numbers ORDER BY number ASC`

	rows, err := r.db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var phone *models.PhoneRaw
		if err = rows.Scan(&phone.ID, &phone.Number); err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phones, err
}

func (r *RawDB) Get(id int) (phone *models.PhoneRaw, err error) {
	utils.Titlef("Getting phone number with id %d", id)
	phone = &models.PhoneRaw{}
	statement := `SELECT id, number FROM phone_numbers WHERE id=$1`
	err = r.db.QueryRow(statement, id).Scan(&phone.ID, &phone.Number)
	return phone, err
}

// FindPhone Search (by number)
func (r *RawDB) FindPhone(number string) (p *models.PhoneRaw, err error) {
	utils.Titlef("Searching for phone number %q", number)
	p = &models.PhoneRaw{}
	statement := `SELECT id, number FROM phone_numbers WHERE number=$1`

	row := r.db.QueryRow(statement, number)
	err = row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NoRecordFoundError{}
		}
		return nil, err
	}
	return p, err
}

// FindPhones : Filter by number
func (r *RawDB) FindPhones(number string) (phones []models.PhoneRaw, err error) {
	// utils.Titlef("FindPhones all phone numbers with number=%s", number)
	statement := `SELECT id, number FROM phone_numbers WHERE number=$1`

	rows, err := r.db.Query(statement, number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var phone models.PhoneRaw
		if err = rows.Scan(&phone.ID, &phone.Number); err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phones, err
}

// UpdatePhone CR"U"D
func (r *RawDB) UpdatePhone(p *models.PhoneRaw) error {
	// utils.Titlef("Updating phone to %#v", p)
	statement := `UPDATE phone_numbers SET number=$2 WHERE id=$1`
	_, err := r.db.Exec(statement, p.ID, p.Number)
	return err
}

// DeletePhone CRU"D"
func (r *RawDB) DeletePhone(id int) error {
	// utils.Titlef("Deleting phone number with id %d", id)
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := r.db.Exec(statement, id)
	return err
}

func (r *RawDB) Close() error {

	logger.Debug("Closing rawdb")
	return r.db.Close()
}

func (r *RawDB) Seed(data []string) error {
	utils.Title("Seeding RawDB with sample data")
	for _, number := range data {
		id, err := r.InsertPhone(number)
		if err != nil {
			return err
		}
		logger.Infof("Successfully inserted phone number %s with id %d", number, id)
	}
	return nil
}

func (r *RawDB) Migrate() error {
	utils.Title("Migrating RawDB tables")

	_, err := r.db.Exec(schema)
	if err != nil {
		return err
	}
	logger.Info("Successfully migrated rawdb")
	return nil
}

func InitRawDB(reset bool) (IPhoneDB[models.PhoneRaw], error) {
	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	if reset {
		err = ResetRawDB("postgres", config.GetPGConnectionString(), config.DBDatabaseNameRaw)
		if err != nil {
			return nil, err
		}
	}

	connectionString := config.GetDBConnectionString(utils.RawDB)

	rawDB, err := openRawDB("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = rawDB.Migrate()
	if err != nil {
		err = rawDB.Close()
		return nil, err
	}

	return rawDB, err
}

func openRawDB(driverName, dataSourceName string) (IPhoneDB[models.PhoneRaw], error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	logger.Info("Successfully connected to rawdb")

	return &RawDB{db: db}, nil
}

func ResetRawDB(driverName, connectionString, dbName string) error {
	utils.Title("ResetRawDB")
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	if err := disconnectAllUsersFromRawDB(db, dbName); err != nil {
		return err
	}

	if err := resetRawDB(db, dbName); err != nil {
		return err
	}
	return nil
}

func disconnectAllUsersFromRawDB(db *sql.DB, dbName string) error {
	statement := fmt.Sprintf(
		`SELECT pg_terminate_backend(pg_stat_activity.pid)
	FROM pg_stat_activity
	WHERE pg_stat_activity.datname = '%s' 
		AND pid <> pg_backend_pid();`,
		dbName,
	)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	logger.Debug("Successfully disconnected all users from rawdb")
	return nil
}

func resetRawDB(db *sql.DB, dbName string) error {
	utils.Title("resetRawDB")
	statement := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return createRawDB(db, dbName)
}

func createRawDB(db *sql.DB, dbName string) error {
	utils.Title("creating rawdb")
	statement := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}
