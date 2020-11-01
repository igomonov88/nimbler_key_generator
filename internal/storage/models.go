package storage

type Key struct {
	ID     int    `db:"id"`
	Value  string `db:"key"`
	IsUsed bool   `db:"is_used"`
}
