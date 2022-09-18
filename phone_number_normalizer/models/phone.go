package models

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
