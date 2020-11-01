package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Using constants in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Create keys table.",
		Script: `
	CREATE TABLE IF NOT EXISTS keys (
  		id SERIAL PRIMARY KEY, 
  		key VARCHAR(6) UNIQUE
	);`,
	},
	{
		Version:     2,
		Description: "Create used keys table.",
		Script: `
	CREATE TABLE IF NOT EXISTS used_keys (
		id SERIAL PRIMARY KEY,
		key VARCHAR(6) UNIQUE
	);`,
	},
	{
		Version:     3,
		Description: "Add index to used_keys table.",
		Script:      `CREATE INDEX idx_key ON used_keys(key);`,
	},
	{
		Version:     4,
		Description: "Add is_used flag to keys table.",
		Script:      `ALTER TABLE keys ADD COLUMN is_used bool not null default false;`,
	},
	{
		Version:     5,
		Description: "Remove used_keys table.",
		Script:      `DROP TABLE IF EXISTS used_keys;`,
	},
}

