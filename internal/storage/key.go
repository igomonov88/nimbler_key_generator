package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ErrNotFound is used when a key is not found in keys table so in that case we will need
// to generate key and put it to database.
var ErrNotFound = errors.New("Key not found")

// AddGeneratedKeys add a slice of generated keys to keys table.
func AddGeneratedKeys(ctx context.Context, db *sqlx.DB, keys []string) error {
	ctx, span := trace.StartSpan(ctx, "internal.storage.AddGeneratedKeys")
	defer span.End()

	const query = `INSERT INTO keys (key) values ($1)`
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "begin transaction in add generated keys method")
	}
	for i := range keys {
		_, err := tx.Exec(query, keys[i])
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return errors.Wrapf(err, "failed to rollback")
			}
			return errors.Wrapf(err, "inserting key: %s", keys[i])
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "failed to commit")
	}

	return nil
}

// GetKey returns generated key from keys table. This function guarantee that key, which
// it returns is not used by anyone.
func GetKey(ctx context.Context, db *sqlx.DB) (string, error) {
	ctx, span := trace.StartSpan(ctx, "internal.storage.GetKey")
	defer span.End()

	const (
		getQuery    = `SELECT * FROM keys where is_used = false limit 1`
		updateQuery = `UPDATE keys SET is_used = true where key = $1`
	)
	var key Key
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return key.Value, errors.Wrapf(err, "begin transaction")
	}

	row := tx.QueryRowContext(ctx, getQuery)
	if err := row.Scan(&key.ID, &key.Value, &key.IsUsed); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", errors.Wrapf(err, "failed to roll back transaction")
		}
		return key.Value, errors.Wrapf(err, "selecting key from keys table")
	}

	if _, err := tx.ExecContext(ctx, updateQuery, key.Value); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", errors.Wrapf(err, "failed to roll back transaction")
		}
		return key.Value, errors.Wrapf(err, "updating is_used parameter in keys table ")
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", errors.Wrapf(err, "failed to roll back transaction")
		}
		return key.Value, errors.Wrapf(err, "commit result of tra ")
	}

	return key.Value, nil
}

// AddKey used only then all of the generate kyes are used and we go out of generate keys
// per day capacity. In that keys we generate key by keygen package, add it to keys table
// and if key is successfully inserted to keys table, we return generated key.
// it returns is not used by anyone.
func AddKey(ctx context.Context, db *sqlx.DB, key string) error {
	ctx, span := trace.StartSpan(ctx, "internal.storage.AddKey")
	defer span.End()

	const query = `INSERT INTO keys (key, is_used) VALUES ($1, $2)`

	if _, err := db.ExecContext(ctx, query, key, false); err != nil {
		return errors.Wrapf(err, "inserting generated key to keys")
	}

	return nil
}

// ReuseKey used only then such key time ti leave is expired and to not generate new key
// we just reuse it. Be careful of using that functionality because users can get to
// wrong url path instead of pasted link.
func ReuseKeys(ctx context.Context, db *sqlx.DB, keys []string) error {
	ctx, span := trace.StartSpan(ctx, "internal.storage.ReuseKey")
	defer span.End()

	const query = `UPDATE keys SET is_used = false where key = $1`

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "begin transaction in reuse generated keys method")
	}
	for i := range keys {
		_, err := tx.Exec(query, keys[i])
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return errors.Wrapf(err, "failed to rollback")
			}
			return errors.Wrapf(err, "inserting key: %s", keys[i])
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "failed to commit")
	}

	return nil
}
