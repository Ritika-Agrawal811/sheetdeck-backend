package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

/**
 * Converts a string to pgtype.Text. Handles empty strings as null values.
 * @param string
 * @returns pgtype.Text
 */
func PgText(text string) pgtype.Text {
	if text == "" {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{Valid: true, String: text}
}

/**
 * Converts int64 to pgtype.Int8
 */
func PgInt8(number int64) pgtype.Int8 {
	return pgtype.Int8{Int64: number, Valid: true}
}

/**
 * Converts a string to pgtype.UUID
 * Handles invalid UUIDs
 */
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
