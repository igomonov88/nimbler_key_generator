package keygen_test

import (
	"context"
	"testing"

	"github.com/igomonov88/nimbler_key_generator/internal/keygen"
	"github.com/igomonov88/nimbler_key_generator/internal/tests"
)

func TestGenerate(t *testing.T) {
	t.Log("Given the need to test key generate functionality:")
	{
		_, err := keygen.Generate(context.Background())
		if err != nil {
			t.Fatalf("\t%s\tShould be able to generate key: %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to generate key.", tests.Success)
	}
}
