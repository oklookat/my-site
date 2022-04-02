package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"servus/apps/elven/model"
)

// create tables in database from SQL file.
func migrate(values []string) error {

	if values == nil || len(values) < 1 {
		return errors.New("empty .sql path")
	}

	var sqlPath = path.Clean(values[0])
	script, err := ioutil.ReadFile(sqlPath)

	if err != nil {
		return fmt.Errorf("migration failed. Read SQL file error: %w", err)
	}

	if _, err = model.StringAdapter.Exec(string(script)); err != nil {
		return fmt.Errorf("migration failed. Failed to execute SQL file: %w", err)
	}

	return nil
}

// delete tables from database.
func rollback() error {
	_, err := model.StringAdapter.Exec(`
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;
	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO public;
	`)

	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	return nil
}
