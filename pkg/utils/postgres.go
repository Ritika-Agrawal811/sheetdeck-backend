package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// PgText converts a string to pgtype.Text, handling empty strings as null values.
func PgText(text string) pgtype.Text {
	// If the input string is empty, return a null pgtype.Text
	if text == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{Valid: true, String: text}
}

// StringToUUID converts a string to pgtype.UUID, returning an error if the string is not a valid UUID.
func StringToUUID(id string) (pgtype.UUID, error) {
	uuidBytes, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, err
	}

	var pgtypeUUID pgtype.UUID

	copy(pgtypeUUID.Bytes[:], uuidBytes[:])
	pgtypeUUID.Valid = true

	return pgtypeUUID, nil
}
