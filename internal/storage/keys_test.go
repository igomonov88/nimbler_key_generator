package storage_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/igomonov88/nimbler_key_generator/internal/keygen"
	"github.com/igomonov88/nimbler_key_generator/internal/storage"
	"github.com/igomonov88/nimbler_key_generator/internal/tests"
)

func TestGetKeys(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()
	t.Log("Given the need to test get keys from storage functionality:")

	keys := make([]string, 0, 5)

	for i := 0; i < cap(keys); i++ {
		key, err := keygen.Generate(context.Background())
		if err != nil {
			t.Fatalf("\t%s\tShould be able to generate key: %s", tests.Failed, err)
		}
		keys = append(keys, key)
	}

	// Operate with keys
	{
		err := storage.AddGeneratedKeys(context.Background(), db, keys)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to add new keys to storage: %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to add new keys to storage.", tests.Success)

		for i := 0; i < len(keys); i++ {
			_, err := storage.GetKey(context.Background(), db)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get generated key from storage: %s", tests.Failed, err)
			}
		}
		t.Logf("\t%s\tShould be able to get generated key from storage.", tests.Success)

		if err := storage.ReuseKeys(context.Background(), db, keys); err != nil {
			t.Fatalf("\t%s\tShould be able to reuse keys: %s", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to reuse keys.", tests.Success)

		key, err := keygen.Generate(context.Background())
		if err != nil {
			t.Fatalf("\t%s\tShould be able to generate key: %s", tests.Failed, err)
		}

		if err := storage.AddKey(context.Background(), db, key); err != nil {
			t.Fatalf("\t%s\tShould be able to add generated key to storage: %s", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to add generated key to storage.", tests.Success)

	}
}

func TestGetKeysNegativeScenarios(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()
	t.Log("Given the need to test get keys from storage functionality when keys are not exist:")

	_, err := storage.GetKey(context.Background(), db)
	if err == nil && err != sql.ErrNoRows {
		t.Fatalf("\t%s\tShould not be able to get generated key from storage: %s", tests.Failed, err)
	}
	t.Logf("\t%s\tShould not be able to get generated key from storage.", tests.Success)

	keys := make([]string, 0, 3)

	for i := 0; i < cap(keys); i++ {
		key, err := keygen.Generate(context.Background())
		if err != nil {
			t.Fatalf("\t%s\tShould be able to generate key: %s", tests.Failed, err)
		}
		keys = append(keys, key)
	}

	if err := storage.ReuseKeys(context.Background(), db, keys); err != nil {
		t.Fatalf("\t%s\tShould be able to reuse keys: %s", tests.Failed, err)
	}
	t.Logf("\t%s\tShould be able to reuse keys.", tests.Success)

}
