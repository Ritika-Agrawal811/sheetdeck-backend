package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func PgText(text string) pgtype.Text {
	// If the input string is empty, return a null pgtype.Text
	if text == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{Valid: true, String: text}
}
