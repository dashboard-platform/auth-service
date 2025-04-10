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

type argonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

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
