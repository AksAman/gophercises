package models

import (
	"fmt"

	"gorm.io/gorm"
)

type IPhone interface {
	GetNumber() string
	SetNumber(string)

	GetID() int
}

type Phone struct {
	Number string
}

type PhoneRaw struct {
	ID int
	Phone
}

// implement IPhone
func (p *PhoneRaw) GetNumber() string {
	return p.Number
}

func (p *PhoneRaw) SetNumber(number string) {
	p.Number = number
}

func (p *PhoneRaw) GetID() int {
	return p.ID
}

func (p *PhoneRaw) String() string {
	return fmt.Sprintf("PhoneRaw: %d %s", p.ID, p.Number)
}

type PhoneSqlx struct {
	ID int
	Phone
}

// Implement IPhone

func (p *PhoneSqlx) GetNumber() string {
	return p.Number
}

func (p *PhoneSqlx) SetNumber(number string) {
	p.Number = number
}

func (p *PhoneSqlx) GetID() int {
	return p.ID
}

func (p *PhoneSqlx) String() string {
	return fmt.Sprintf("PhoneSqlx: %d %s", p.ID, p.Number)
}

type PhoneGorm struct {
	gorm.Model
	Phone
}

func (p *PhoneGorm) String() string {
	return fmt.Sprintf("PhoneGorm: %d %s", p.ID, p.Number)
}
