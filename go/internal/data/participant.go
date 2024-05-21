package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"nathejk.dk/internal/validator"
)

type FilterType string

const (
	FilterTypeRadio    FilterType = "radio"
	FilterTypeCheckbox FilterType = "checkbox"
	FilterTypeText     FilterType = "freetext"
)

type FilterOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Filter struct {
	Slug        string         `json:"slug"`
	Label       string         `json:"label"`
	Placeholder string         `json:"placeholder"`
	Type        FilterType     `json:"type"`
	Options     []FilterOption `json:"options"`
}

type Participant struct {
	ID        int64     `json:"id"`
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Version   uint64    `json:"version"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Team      string    `json:"team"`
	Days      []string  `json:"days"`
	Transport string    `json:"transport"`
	SeatCount string    `json:"seatCount"`
	Info      string    `json:"info"`
	Video     string    `json:"video"`
	Paid      bool      `json:"paid"`
	Diet      string    `json:"diet"`
	Tshirt    string    `json:"tshirt"`
}

func (p *Participant) Validate(v validator.Validator) {
	v.Check(p.Name != "", "name", "must be provided")
	v.Check(len(p.Name) <= 500, "name", "must not be more than 500 bytes long")
}

type ParticipantModel struct {
	DB *sql.DB
}

func (m ParticipantModel) Insert(p *Participant) error {
	p.UUID = uuid.New().String()
	query := `
		INSERT INTO participant (name, address, email, phone, team, days, transport, seatCount, info, video, uuid, diet, tshirt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, version`
	args := []any{p.Name, p.Address, p.Email, p.Phone, p.Team, pq.Array(p.Days), p.Transport, p.SeatCount, p.Info, p.Video, p.UUID, p.Diet, p.Tshirt}
	return m.DB.QueryRow(query, args...).Scan(&p.ID, &p.CreatedAt, &p.Version)
}

func (m ParticipantModel) GetByPayment(p *Payment) (*Participant, error) {
	if p == nil {
		return nil, ErrRecordNotFound
	}
	return m.Get(p.Message)
}

func (m ParticipantModel) Get(id string) (*Participant, error) {
	if id == "" {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, uuid, created_at, version, name, address, email, phone, team, days, transport, seatCount, info, video, diet, tshirt, (SELECT COALESCE(SUM(amount), 0) FROM payment WHERE message = uuid) >= 50 FROM participant
		WHERE uuid = $1`

	var p Participant
	err := m.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.UUID,
		&p.CreatedAt,
		&p.Version,
		&p.Name,
		&p.Address,
		&p.Email,
		&p.Phone,
		&p.Team,
		pq.Array(&p.Days),
		&p.Transport,
		&p.SeatCount,
		&p.Info,
		&p.Video,
		&p.Diet,
		&p.Tshirt,
		&p.Paid,
	)
	// Handle any errors. If there was no matching participant found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &p, nil
}

// Add a placeholder method for updating a specific record in the participants table.
func (m ParticipantModel) Update(p *Participant) error {
	query := `
		UPDATE participant
		SET
			name = $3,
			address = $4,
			email = $5,
			phone = $6,
			team = $7,
			days = $8,
			transport = $9,
			seatCount = $10,
			info = $11,
			video = $12,
			diet = $13,
			tshirt = $14,
			version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING version`

	args := []any{
		p.ID,
		p.Version,
		p.Name,
		p.Address,
		p.Email,
		p.Phone,
		p.Team,
		pq.Array(p.Days),
		p.Transport,
		p.SeatCount,
		p.Info,
		p.Video,
		p.Diet,
		p.Tshirt,
	}
	// Execute the SQL query. If no matching row could be found, we know the record
	// version has changed (or the record has been deleted) and we return ErrEditConflict error.
	err := m.DB.QueryRow(query, args...).Scan(&p.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m ParticipantModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM participant WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m ParticipantModel) GetAll(name string, genres []string, filters Filters) ([]*Participant, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, version, name, address, email, phone, team, days, transport, seatCount, info, video, CASE WHEN diet IS NULL THEN '' ELSE diet END, CASE WHEN tshirt IS NULL THEN '' ELSE tshirt END, (SELECT COALESCE(SUM(amount), 0) FROM payment WHERE message = uuid) >= 50
		FROM participant
		WHERE (LOWER(name) = LOWER($1) OR $1 = '') AND date_part('year', created_at) = date_part('year', now())
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	args := []any{name, filters.limit(), filters.offset()}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	participants := []*Participant{}
	for rows.Next() {
		var p Participant
		err := rows.Scan(
			&totalRecords,
			&p.ID,
			&p.CreatedAt,
			&p.Version,
			&p.Name,
			&p.Address,
			&p.Email,
			&p.Phone,
			&p.Team,
			pq.Array(&p.Days),
			&p.Transport,
			&p.SeatCount,
			&p.Info,
			&p.Video,
			&p.Diet,
			&p.Tshirt,
			&p.Paid,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		participants = append(participants, &p)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return participants, metadata, nil
}
