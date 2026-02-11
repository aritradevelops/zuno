package repository

import "github.com/google/uuid"

type Actor struct {
	UID   uuid.UUID
	Scope string
}
