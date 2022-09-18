package dblib

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/AksAman/gophercises/phone/models"
	"github.com/AksAman/gophercises/phone/utils"
)

type SqlxDB struct {
	db *sqlx.DB
}

func (s *SqlxDB) Close() error {
	logger.Debug("Closing sqlx db")
	return s.db.Close()
}

func (s *SqlxDB) Seed(data []string) error {
	utils.Title("Seeding RawDB with sample data")
	for _, number := range data {
		id, err := s.InsertPhone(number)
		if err != nil {
			return err
		}
		logger.Infof("Successfully inserted phone number %s with id %d", number, id)
	}
	return nil
}

func (s *SqlxDB) Migrate() error {
	utils.Title("Migrating RawDB tables")

	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}
	logger.Info("Successfully migrated rawdb")
	return nil
}

func (s *SqlxDB) InsertPhone(number string) (id int, err error) {
	id = -1
	statement := `INSERT INTO phone_numbers(number) VALUES($1) RETURNING id`

	row := s.db.QueryRow(statement, number)
	err = row.Scan(&id)
	if err != nil {
		return
	}
	return
}

func (s *SqlxDB) All() (phones []models.PhoneSqlx, err error) {
	utils.Title("Getting all phone numbers")
	statement := `SELECT id, number FROM phone_numbers ORDER BY number ASC`

	rows, err := s.db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var phone models.PhoneSqlx
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

func (s *SqlxDB) Get(id int) (phone *models.PhoneSqlx, err error) {
	utils.Titlef("Getting phone number with id %d", id)
	phone = &models.PhoneSqlx{}
	statement := `SELECT id, number FROM phone_numbers WHERE id=$1`
	err = s.db.QueryRow(statement, id).Scan(&phone.ID, &phone.Number)
	return phone, err
}

func (s *SqlxDB) FindPhone(number string) (p *models.PhoneSqlx, err error) {
	utils.Titlef("Searching for phone number %q", number)
	p = &models.PhoneSqlx{}
	statement := `SELECT id, number FROM phone_numbers WHERE number=$1`

	row := s.db.QueryRow(statement, number)
	err = row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NoRecordFoundError{}
		}
		return nil, err
	}
	return p, err
}

func (s *SqlxDB) FindPhones(number string) (phones []models.PhoneSqlx, err error) {
	// utils.Titlef("FindPhones all phone numbers with number=%s", number)
	statement := `SELECT id, number FROM phone_numbers WHERE number=$1`

	rows, err := s.db.Query(statement, number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var phone models.PhoneSqlx
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
func (s SqlxDB) UpdatePhone(phone *models.PhoneSqlx) error {
	statement := `UPDATE phone_numbers SET number=$2 WHERE id=$1`
	_, err := s.db.Exec(statement, phone.ID, phone.Number)
	return err
}

func (s SqlxDB) DeletePhone(id int) error {
	// utils.Titlef("Deleting phone number with id %d", id)
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := s.db.Exec(statement, id)
	return err
}

func InitSqlxDB(reset bool) (IPhoneDB[models.PhoneSqlx], error) {
	config, err := utils.LoadConfig()
	if err != nil {
		return nil, err
	}

	if reset {
		err = ResetSqlxDB("postgres", config.GetPGConnectionString(), config.DBDatabaseNameSqlx)
		if err != nil {
			return nil, err
		}
	}

	connectionString := config.GetDBConnectionString(utils.SqlxDB)

	db, err := openSqlxDB("postgres", connectionString)
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

func openSqlxDB(driverName, dataSourceName string) (IPhoneDB[models.PhoneSqlx], error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	logger.Info("Successfully connected to rawdb")

	return &SqlxDB{db: db}, nil
}

func ResetSqlxDB(driverName, connectionString, dbName string) error {
	utils.Title("ResetRawDB")
	db, err := sqlx.Connect(driverName, connectionString)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	if err := disconnectAllUsersFromSqlxDB(db, dbName); err != nil {
		return err
	}

	if err := resetSqlxDB(db, dbName); err != nil {
		return err
	}
	return nil
}

func disconnectAllUsersFromSqlxDB(db *sqlx.DB, dbName string) error {
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

func resetSqlxDB(db *sqlx.DB, dbName string) error {
	utils.Title("resetSqlxDB")
	statement := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return createSqlxDB(db, dbName)
}

func createSqlxDB(db *sqlx.DB, dbName string) error {
	utils.Title("creating rawdb")
	statement := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}
