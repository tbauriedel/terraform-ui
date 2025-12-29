package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"golang.org/x/crypto/argon2"
)

const (
	// regex pattern for argon2id hash
	// matches as plain text
	// $argon2id$v=19$m=65536,t=3,p=1$+kn21LcRetAkE7zObeS3xA$FzdfjLWlAiJbHLE+Rjm2hBMUmMb3TdmWQ7AMTtryYfk
	// [1] -> 65536
	// [2] -> 3
	// [3] -> 1
	// [4] -> +kn21LcRetAkE7zObeS3xA
	// [5] -> FzdfjLWlAiJbHLE+Rjm2hBMUmMb3TdmWQ7AMTtryYfk .
	phcRegexPattern = `^\$argon2id\$v=19\$m=(\d+),t=(\d+),p=(\d+)\$([^$]+)\$([^$]+)$`
)

// HashPasswordString hashes a password using argon2id.
// Salt and cost parameters are based on the recommendations from https://github.com/P-H-C/phc-winner-argon2
//
// Returns the hash as a byte slice.
func HashPasswordString(password string, params config.HashingParams) string {
	// generate salt if not provided
	if len(params.Salt) == 0 {
		params.Salt = make([]byte, params.SaltLength)
		// Read random bytes from the crypto/rand package and fills the slice.
		// never generates an error. panics in case of failure
		_, _ = rand.Read(params.Salt)
	}

	// generate hash
	hash := argon2.IDKey([]byte(password),
		params.Salt,
		params.Iterations,
		params.MemoryCost,
		params.ThreadsCount,
		params.KeyLength,
	)

	// encode salt and hash to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(params.Salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// format to phc format: $argon2id$v=19$m=MEM,t=ITER,p=PAR$salt$hash
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		params.MemoryCost, params.Iterations, params.ThreadsCount, b64Salt, b64Hash)

	return encoded
}

// ExportParamsFromHash extracts the parameters from a phc encoded hash.
// Returns the parameters as a HashingParams struct.
//
// Uses a fixed regex pattern to extract the parameters.
func ExportParamsFromHash(encoded string) (config.HashingParams, error) {
	var p config.HashingParams

	// regex to extract parameters from encoded string
	r := regexp.MustCompile(phcRegexPattern)

	// encoded string matches regex pattern
	m := r.FindStringSubmatch(encoded)
	if m == nil {
		return p, fmt.Errorf("invalid argon2id phc format provided. pattern does not match")
	}

	// memory cost
	v, err := strconv.ParseUint(m[1], 10, 32)
	if err != nil {
		return p, fmt.Errorf("cant extract memory cost from hash. invalid format")
	}

	p.MemoryCost = uint32(v)

	// iterations
	v, err = strconv.ParseUint(m[2], 10, 32)
	if err != nil {
		return p, fmt.Errorf("cant extract iterations from hash. invalid format")
	}

	p.Iterations = uint32(v)

	// threads / parallelism
	v, err = strconv.ParseUint(m[3], 10, 8)
	if err != nil {
		return p, fmt.Errorf("cant extract threads count from hash. invalid format")
	}

	p.ThreadsCount = uint8(v)

	// encode base 64 salt
	p.Salt, err = base64.RawStdEncoding.DecodeString(m[4])
	if err != nil {
		return p, fmt.Errorf("invalid salt provided. cant decode")
	}

	// encode base 64 hash
	hash, err := base64.RawStdEncoding.DecodeString(m[5])
	if err != nil {
		return p, fmt.Errorf("invalid hash provided. cant decode")
	}

	// check if salt and hash are not too long.
	// If so, the conversion from int (len) to uin32 would panic.
	if len(p.Salt) > math.MaxUint32 || len(hash) > math.MaxUint32 {
		return p, fmt.Errorf("invalid hash provided. salt or hash too long")
	}

	p.SaltLength = uint32(len(p.Salt)) //nolint:gosec
	p.KeyLength = uint32(len(hash))    //nolint:gosec

	return p, nil
}
