package dblib

import (
	"github.com/AksAman/gophercises/phone/utils"
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

type NoRecordFoundError struct {
}

func (e *NoRecordFoundError) Error() string {
	return "No record found"
}

func init() {
	utils.InitializeLogger("")

	logger = utils.Logger
}

type IPhoneDB[T interface{}] interface {
	Close() error
	Seed([]string) error
	Migrate() error
	InsertPhone(string) (id int, err error)    // C
	All() (phones []T, err error)              // R (all)
	Get(int) (phone *T, err error)             // R (by id)
	FindPhone(string) (phone *T, err error)    // Search (by number)
	FindPhones(string) (phones []T, err error) // Search (by number)
	UpdatePhone(phone *T) error                // U
	DeletePhone(id int) error                  // D
}
