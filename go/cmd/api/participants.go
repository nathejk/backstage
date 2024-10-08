package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"nathejk.dk/internal/data"
	"nathejk.dk/internal/validator"
)

func (app *application) showParticipantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.ReadUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	participant, err := app.models.Participants.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	/*
		participant.TargetGroup.Filters = []data.Filter{
			data.Filter{Slug: "gender", Label: "Køn", Type: data.FilterTypeRadio, Options: []data.FilterOption{
				data.FilterOption{Label: "Mand", Value: "M"},
				data.FilterOption{Label: "Kvinde", Value: "F"},
				data.FilterOption{Label: "Ukendt", Value: "X"},
			}},
		}
	*/
	err = app.writeJSON(w, http.StatusOK, envelope{"participant": participant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createParticipantHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string   `json:"name"`
		Address   string   `json:"address"`
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		Team      string   `json:"team"`
		Days      []string `json:"days"`
		Transport string   `json:"transport"`
		SeatCount string   `json:"seatCount"`
		Diet      string   `json:"diet"`
		Tshirt    string   `json:"tshirt"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	participant := &data.Participant{
		Name:      input.Name,
		Address:   input.Address,
		Email:     input.Email,
		Phone:     input.Phone,
		Team:      input.Team,
		Days:      input.Days,
		Transport: input.Transport,
		SeatCount: input.SeatCount,
		Diet:      input.Diet,
		Tshirt:    input.Tshirt,
	}

	v := validator.New()

	if participant.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Participants.Insert(participant); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/participants/%d", participant.ID))

	err := app.writeJSON(w, http.StatusCreated, envelope{"participant": participant}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) requestPaymentParticipantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.ReadUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	/*
		participant, err := app.models.Participants.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}*/

	var input struct {
		Phone string `json:"phone"`
	}
	if err = app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	participant, err := app.models.Participants.Get(id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	amount := 50
	if participant.Tshirt != "ingen" {
		log.Println("t-shirt")
		//amount += 175
	}
	link := fmt.Sprintf("https://www.mobilepay.dk/erhverv/betalingslink/betalingslink-svar?phone=775771&amount=%d&lock=1&comment=%s", amount, id)
	err = app.sms.Send(input.Phone, link)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"status": "ok"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateParticipantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.ReadUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	participant, err := app.models.Participants.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If the request contains a X-Expected-Version header, verify that the
	// version in the database matches the expected version specified in the header.
	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(participant.Version), 32) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Name      *string  `json:"name"`
		Address   *string  `json:"address"`
		Email     *string  `json:"email"`
		Phone     *string  `json:"phone"`
		Team      *string  `json:"team"`
		Days      []string `json:"days"`
		Transport *string  `json:"transport"`
		SeatCount *string  `json:"seatCount"`
		Info      *string  `json:"info"`
		Video     *string  `json:"video"`
	}
	if err = app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		participant.Name = *input.Name
	}
	if input.Address != nil {
		participant.Address = *input.Address
	}
	if input.Email != nil {
		participant.Email = *input.Email
	}
	if input.Phone != nil {
		participant.Phone = *input.Phone
	}
	if input.Team != nil {
		participant.Team = *input.Team
	}
	if input.Days != nil {
		participant.Days = input.Days
	}
	if input.Transport != nil {
		participant.Transport = *input.Transport
	}
	if input.SeatCount != nil {
		participant.SeatCount = *input.SeatCount
	}
	if input.Info != nil {
		participant.Info = *input.Info
	}
	if input.Video != nil {
		participant.Video = *input.Video
	}
	v := validator.New()
	if participant.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Participants.Update(participant)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"participant": participant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteParticipantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Participants.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "participant successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if input.Filters.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	participants, metadata, err := app.models.Participants.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "participants": participants}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
