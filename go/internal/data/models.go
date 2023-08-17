package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Participants interface {
		Insert(*Participant) error
		Get(id string) (*Participant, error)
		GetByPayment(p *Payment) (*Participant, error)
		Update(*Participant) error
		Delete(id int64) error
		GetAll(string, []string, Filters) ([]*Participant, Metadata, error)
	}
	Payments interface {
		Insert(*Payment) error
		Latest() time.Time
	}
	Permissions interface {
		AddForUser(int64, ...string) error
		GetAllForUser(int64) (Permissions, error)
	}
	Tokens interface {
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		Insert(token *Token) error
		DeleteAllForUser(scope string, userID int64) error
	}
	Users interface {
		Insert(*User) error
		GetByEmail(string) (*User, error)
		Update(*User) error
		GetForToken(string, string) (*User, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Participants: ParticipantModel{DB: db},
		Payments:     PaymentModel{DB: db},
		Permissions:  PermissionModel{DB: db},
		Tokens:       TokenModel{DB: db},
		Users:        UserModel{DB: db},
	}
}
