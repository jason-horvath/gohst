package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Params defines the parameters for Argon2id hashing
type Argon2Params struct {
    Memory      uint32
    Iterations  uint32
    Parallelism uint8
    SaltLength  uint32
    KeyLength   uint32
}

// DefaultArgon2Params returns recommended parameters for Argon2id
func DefaultArgon2Params() *Argon2Params {
    return &Argon2Params{
        Memory:      64 * 1024, // 64MB
        Iterations:  4,         // Time cost
        Parallelism: 2,         // Threads
        SaltLength:  16,        // 16 bytes
        KeyLength:   32,        // 32 bytes
    }
}

// HashPassword hashes a password using Argon2id with default parameters
func HashPassword(password string) (string, error) {
    return HashPasswordWithParams(password, DefaultArgon2Params())
}

// HashPasswordWithParams hashes a password using Argon2id with custom parameters
func HashPasswordWithParams(password string, params *Argon2Params) (string, error) {
    if password == "" {
        return "", errors.New("password cannot be empty")
    }

    // Generate a random salt
    salt := make([]byte, params.SaltLength)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    // Hash the password
    hash := argon2.IDKey(
        []byte(password),
        salt,
        params.Iterations,
        params.Memory,
        params.Parallelism,
        params.KeyLength,
    )

    // Encode as PHC format: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    encodedHash := fmt.Sprintf(
        "$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
        params.Memory,
        params.Iterations,
        params.Parallelism,
        b64Salt,
        b64Hash,
    )

    return encodedHash, nil
}

// CheckPassword verifies if a password matches an Argon2id hash
func CheckPassword(password, encodedHash string) (bool, error) {
    if password == "" || encodedHash == "" {
        return false, errors.New("password or hash is empty")
    }

    // Parse the hash string
    vals := strings.Split(encodedHash, "$")
    if len(vals) != 6 {
        return false, errors.New("invalid hash format")
    }

    if vals[1] != "argon2id" {
        return false, fmt.Errorf("incompatible hash algorithm: %s", vals[1])
    }

    var version int
    _, err := fmt.Sscanf(vals[2], "v=%d", &version)
    if err != nil {
        return false, errors.New("invalid hash version")
    }
    if version != 19 {
        return false, fmt.Errorf("incompatible hash version: %d", version)
    }

    // Parse parameters
    params := &Argon2Params{}
    _, err = fmt.Sscanf(
        vals[3],
        "m=%d,t=%d,p=%d",
        &params.Memory,
        &params.Iterations,
        &params.Parallelism,
    )
    if err != nil {
        return false, errors.New("invalid hash parameters")
    }

    // Decode salt and hash
    salt, err := base64.RawStdEncoding.DecodeString(vals[4])
    if err != nil {
        return false, errors.New("invalid salt encoding")
    }

    decodedHash, err := base64.RawStdEncoding.DecodeString(vals[5])
    if err != nil {
        return false, errors.New("invalid hash encoding")
    }
    params.KeyLength = uint32(len(decodedHash))

    // Compute hash with same parameters and salt
    comparisonHash := argon2.IDKey(
        []byte(password),
        salt,
        params.Iterations,
        params.Memory,
        params.Parallelism,
        params.KeyLength,
    )

    // Constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
