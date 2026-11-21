package database

import (
	"context"
	"errors"
	"fff-hedgedoc/internal/logging"
)

func GetAllNotes(ctx context.Context) ([]string, error) {
	if !DatabaseInitialized {
		return nil, errors.New("database has not been initalized")
	}

	rows, err := db.Query(ctx, `
		SELECT "title" FROM "Notes"
		ORDER BY "updatedAt" DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	noteTitles := make([]string, 0)

	for rows.Next() {
		var title string

		if err := rows.Scan(&title); err != nil {
			return nil, err
		}

		noteTitles = append(noteTitles, title)
		logging.Logger.Infof("Valid Note: %s", title)
	}

	logging.Logger.Infof("notes list size: %d", len(noteTitles))
	return noteTitles, nil
}
