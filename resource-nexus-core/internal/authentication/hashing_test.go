package authentication

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tbauriedel/resource-nexus-core/internal/config"
)

func TestHashPasswordString(t *testing.T) {
	p := config.HashingParams{
		Iterations:   1,
		MemoryCost:   16 * 1024,
		ThreadsCount: 1,
		KeyLength:    32,
		SaltLength:   16,
		Salt:         []byte{},
	}

	hash1 := HashPasswordString("dummy", p)

	if !strings.HasPrefix(hash1, "$argon2id$v=19$m=16384,t=1,p=1") {
		t.Fatal("wrong hash format")
	}

	hash2 := HashPasswordString("dummy", p)

	if hash1 == hash2 {
		t.Log(hash1, hash2)
		t.Fatal("hashes should be different because of the random salt!")
	}
}

func TestExportParamsFromHash(t *testing.T) {
	p := config.HashingParams{
		Iterations:   3,
		MemoryCost:   65536,
		ThreadsCount: 1,
		KeyLength:    32,
		SaltLength:   16,
		Salt:         []byte("foobar"),
	}

	// pass: foobar
	// salt: foobar
	encoded := "$argon2id$v=19$m=65536,t=3,p=1$+kn21LcRetAkE7zObeS3xA$FzdfjLWlAiJbHLE+Rjm2hBMUmMb3TdmWQ7AMTtryYfk"

	params, err := ExportParamsFromHash(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(params, p) {
		t.Fatal("exported params should not be equal to the original ones.\nexported: ", params, "\noriginal: ", p)
	}
}
