package models

import (
	"database/sql"
	"emailparser/internals/entities"
)

// NoteDataModel centralises CRUD operations for notes
type NoteDataModel struct {
	DB *sql.DB
}

// Create inserts a new record into the note table linking it to the user
// provided in user and returns an instantiated Note with the persisted
// information and the new note's id.
func (m *NoteDataModel) Create(details string, user *entities.User) (*entities.Note, error) {
	stmt := `INSERT INTO note (details, user_id) VALUES(?, ?)`
	result, err := m.DB.Exec(stmt, details, user.ID)
	if err != nil {
		return nil, err
	}

	noteID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entities.Note{ID: noteID, User: user, Details: details}, nil
}
