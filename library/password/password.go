package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/skaisanlahti/message-board/utilities/assert"
	"golang.org/x/crypto/argon2"
)

type PasswordOptions struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

func Hash(password string, options PasswordOptions) string {
	salt := RandomBytes(options.SaltLen)
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		options.Time,
		options.Memory,
		options.Threads,
		options.KeyLen,
	)

	hashedPassword := Encode(hash, salt, options)
	return hashedPassword
}

func Verify(hashedPassword, candidate string) bool {
	hash, salt, options, ok := Decode(hashedPassword)
	if !ok {
		return false
	}

	candidateHash := argon2.IDKey(
		[]byte(candidate),
		salt,
		options.Time,
		options.Memory,
		options.Threads,
		options.KeyLen,
	)

	ok = subtle.ConstantTimeCompare(hash, candidateHash) == 1
	return ok
}

func Encode(hash, salt []byte, options PasswordOptions) string {
	return fmt.Sprintf(
		"%s$%s$%d$%d$%d$%d",
		base64.StdEncoding.EncodeToString(hash),
		base64.StdEncoding.EncodeToString(salt),
		options.Time,
		options.Memory,
		options.Threads,
		options.KeyLen,
	)
}

func Decode(hashedPassword string) ([]byte, []byte, PasswordOptions, bool) {
	var hash, salt []byte
	var options PasswordOptions
	var err error

	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return hash, salt, options, false
	}

	hash, err = base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return hash, salt, options, false
	}

	salt, err = base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return hash, salt, options, false
	}

	time, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return hash, salt, options, false
	}

	memory, err := strconv.ParseUint(parts[3], 10, 32)
	if err != nil {
		return hash, salt, options, false
	}

	threads, err := strconv.ParseUint(parts[4], 10, 8)
	if err != nil {
		return hash, salt, options, false
	}

	keyLen, err := strconv.ParseUint(parts[5], 10, 32)
	if err != nil {
		return hash, salt, options, false
	}

	options = PasswordOptions{
		Time:    uint32(time),
		Memory:  uint32(memory),
		Threads: uint8(threads),
		KeyLen:  uint32(keyLen),
		SaltLen: uint32(len(salt)),
	}

	return hash, salt, options, true
}

func CompareOptions(o1, o2 PasswordOptions) bool {
	if o1.Time != o2.Time {
		return false
	}

	if o1.Memory != o2.Memory {
		return false
	}

	if o1.Threads != o2.Threads {
		return false
	}

	if o1.KeyLen != o2.KeyLen {
		return false
	}

	if o1.SaltLen != o2.SaltLen {
		return false
	}

	return true
}

func RandomBytes(length uint32) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	assert.Ok(err, "failed to create random bytes")
	return bytes
}
