package entities

type Note struct {
	ID      int64
	User    *User
	Details string
}