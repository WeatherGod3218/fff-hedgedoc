package database

import "context"

func InitAccessTable(ctx context.Context) error {
	_, err := db.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS access (
				row_id 		UUID PRIMARY KEY,
				note_id		UUID REFERENCES csh_notes(note_id) ON DELETE CASCADE,
				created_by	TEXT NOT NULL,
				created_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)
		`)
	if err != nil {
		return err
	}

	return nil
}
