package models

import "errors"

// ErrNoRecord is a generic error indicating that a given model/record could not be found
var ErrNoRecord = errors.New("models: no matching record found")