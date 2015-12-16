package main

import (
	"database/sql"
	"errors"

	"github.com/avabot/ava/shared/datatypes"
)

var (
	ErrMissingFlexID = errors.New("missing flexid")
)

func saveStructuredInput(m *dt.Msg, rid uint64, pkg, route string) (uint64,
	error) {
	q := `
		INSERT INTO inputs (
			userid,
			flexid,
			flexidtype,
			sentence,
			sentenceannotated,
			commands,
			objects,
			actors,
			times,
			places,
			responseid,
			package,
			route
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`
	in := m.Input
	si := in.StructuredInput
	row := db.QueryRowx(q, in.UserID, in.FlexID, in.FlexIDType, in.Sentence,
		in.SentenceAnnotated, si.Commands, si.Objects, si.Actors,
		si.Times, si.Places, rid, pkg, route)
	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func saveTrainingSentence(in *dt.Input) (int, error) {
	q := `INSERT INTO trainings (sentence) VALUES ($1) RETURNING id`
	var id int
	if err := db.QueryRowx(q, in.Sentence).Scan(&id); err != nil {
		return 0, err
	}
	q = `UPDATE inputs SET trainingid=$1 WHERE id=$2`
	_, err := db.Exec(q, id, in.ID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func updateTraining(trainID int, hitID string, maxAssignments uint) error {
	q := `UPDATE trainings SET foreignid=$1, maxassignments=$2 WHERE id=$3`
	_, err := db.Exec(q, hitID, maxAssignments, trainID)
	if err != nil {
		return err
	}
	return nil
}

func getUser(in *dt.Input) (*dt.User, error) {
	if in.UserID == 0 {
		q := `SELECT userid
		      FROM userflexids
		      WHERE flexid=$1 AND flexidtype=2
		      ORDER BY createdat DESC`
		err := db.Get(&in.UserID, q, in.FlexID)
		if err == sql.ErrNoRows {
			return nil, dt.ErrMissingUser
		} else if err != nil {
			return nil, err
		}
	} else if len(in.FlexID) == 0 {
		return nil, ErrMissingFlexID
	}
	q := `SELECT id, name, email, lastauthenticated, stripecustomerid
	      FROM users
	      WHERE id=$1`
	u := dt.User{}
	if err := db.Get(&u, q, in.UserID); err != nil {
		return nil, err
	}
	return &u, nil
}

func getInputAnnotation(id int) (string, error) {
	var annotation string
	q := `SELECT sentenceannotated FROM inputs WHERE trainingid=$1`
	if err := db.Get(&annotation, q, id); err != nil {
		return "", err
	}
	return annotation, nil
}

func getLastInputFromUser(u *dt.User) (*dt.StructuredInput,
	error) {
	return &(dt.StructuredInput{}), nil
}
