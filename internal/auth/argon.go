// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// argonParams defines the parameters for the Argon2id hashing algorithm.
// These parameters control the memory usage, number of iterations, degree of
// parallelism, and the lengths of the salt and key.
type argonParams struct {
	Memory      uint32 // Memory usage in kibibytes.
	Iterations  uint32 // Number of iterations to perform.
	Parallelism uint8  // Number of threads to use.
	SaltLength  uint32 // Length of the random salt in bytes.
	KeyLength   uint32 // Length of the generated key in bytes.
}

// hashPassword generates a hashed password using the Argon2id algorithm.
// It returns the encoded hash string or an error if the hashing process fails.
//
// The encoded hash includes the algorithm parameters, salt, and hashed key.
//
// Parameters:
//   - password: The plaintext password to hash.
//
// Returns:
//   - string: The encoded Argon2id hash.
//   - error: An error if the hashing process fails.
func hashPassword(password string) (string, error) {
	p := &argonParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory, p.Iterations, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hashed),
	)

	return encoded, nil
}

// verifyPassword verifies a plaintext password against an encoded Argon2id hash.
// It returns true if the password matches the hash, or false otherwise.
//
// Parameters:
//   - encoded: The encoded Argon2id hash to verify against.
//   - password: The plaintext password to verify.
//
// Returns:
//   - bool: True if the password matches the hash, false otherwise.
//   - error: An error if the verification process fails.
func verifyPassword(encoded, password string) (bool, error) {
	var (
		memory, iterations uint32
		parallelism        uint8
		saltB64, hashB64   string
		err                error
	)

	memory, iterations, parallelism, saltB64, hashB64, err = parseHash(encoded)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}
	hashed, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, err
	}

	computed := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hashed)))
	if subtle.ConstantTimeCompare(hashed, computed) == 1 {
		return true, nil
	}
	return false, nil
}

// parseHash parses an encoded Argon2id hash and extracts its parameters, salt,
// and hashed key.
//
// Parameters:
//   - encoded: The encoded Argon2id hash to parse.
//
// Returns:
//   - memory: The memory usage parameter.
//   - iterations: The number of iterations parameter.
//   - parallelism: The degree of parallelism parameter.
//   - saltB64: The base64-encoded salt.
//   - hashB64: The base64-encoded hashed key.
//   - error: An error if the hash format is invalid.
func parseHash(encoded string) (memory uint32, iterations uint32, parallelism uint8, saltB64, hashB64 string, err error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		err = errors.New("invalid hashPassword format")
		return
	}

	paramStr := parts[3]
	saltB64 = parts[4]
	hashB64 = parts[5]

	_, err = fmt.Sscanf(paramStr, "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	return
}
