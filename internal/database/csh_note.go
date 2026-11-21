package database

import "context"

func AddNewNote() {

}

func InitCSHNotes(ctx context.Context) error {
	_, err := db.Exec(ctx, `
		DO $$ BEGIN
			CREATE TYPE note_access_modes AS ENUM ('Whitelist', 'Blacklist');
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS csh_notes (
			note_id 	UUID PRIMARY KEY,
			note_owner 	UUID NOT NULL,
			access_mode note_access_modes NOT NULL DEFAULT 'Blacklist',
			created_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
