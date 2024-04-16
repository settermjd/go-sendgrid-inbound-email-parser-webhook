package entities

type Attachment struct {
	ID, NoteID            int64
	File                  []byte
	ContentType, Filename string
}
