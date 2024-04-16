package models

import (
	"database/sql"
	"emailparser/internals/entities"
	"errors"
	"log"
)

type UserReferenceDataModel struct {
	DB *sql.DB
}

// GetUserFromReference retrieves and returns a User using the reference provided in reference
func (ur *UserReferenceDataModel) GetUserFromReference(reference string) (*entities.User, error) {
	log.Printf("attempting to retrieve user with reference: %s", reference)

	stmt := `SELECT u.id, u.name, u.email, u.phoneNumber 
	FROM user u 
	INNER JOIN user_references ur ON ur.user_id = u.id 
	WHERE ur.id = ?`
	row := ur.DB.QueryRow(stmt, reference)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}
	log.Printf("found user with ref matching: %s\n", reference)

	return user, nil
}
