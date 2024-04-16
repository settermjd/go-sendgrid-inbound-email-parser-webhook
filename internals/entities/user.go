package entities

// User is a simplistic representation of a user within the application
// Currently, the user only has four properties: ID, name, email, and phone number.
type User struct {
	ID                       uint32
	Name, Email, PhoneNumber string
}
