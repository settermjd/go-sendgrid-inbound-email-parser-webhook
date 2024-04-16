package models

import (
	"bytes"
	"database/sql"
	"emailparser/internals/entities"
	"errors"
	"fmt"
	"io"

	"github.com/DusanKasan/parsemail"
)

// AttachmentDataModel centralises CRUD operations for attachments
type AttachmentDataModel struct {
	DB *sql.DB
}

// Create inserts a new attachment record into the attachment table
func (a *AttachmentDataModel) Create(note *entities.Note, attachment parsemail.Attachment) (int64, error) {
	var data bytes.Buffer
	_, err := io.Copy(&data, attachment.Data)
	if err != nil {
		return 0, fmt.Errorf("could not read attachment data")
	}

	stmt := `INSERT INTO attachment (note_id, content_type, filename, file) VALUES(?, ?, ?, ?)`
	result, err := a.DB.Exec(stmt, note.ID, attachment.ContentType, attachment.Filename, data.String())
	if err != nil {
		return 0, err
	}

	attachmentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int64(attachmentID), nil
}

// Get retrieves and returns a single attachment from the database
func (a *AttachmentDataModel) Get(attachmentID int64) (*entities.Attachment, error) {
	stmt := `SELECT content_type, file, filename FROM attachment WHERE id = ?`
	row := a.DB.QueryRow(stmt, attachmentID)
	data := &entities.Attachment{}
	err := row.Scan(&data.ContentType, &data.File, &data.Filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}

	return data, nil
}
