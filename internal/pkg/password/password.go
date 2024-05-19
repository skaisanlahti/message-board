package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Options struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

var DefaultOptions = Options{
	Time:    5,
	Memory:  1024 * 7,
	Threads: 1,
	SaltLen: 32,
	KeyLen:  64,
}

type Hasher struct {
	options Options
}

func NewHasher(options Options) *Hasher {
	return &Hasher{
		options: options,
	}
}

func (hasher *Hasher) Hash(password string) string {
	options := hasher.options
	salt := randomBytes(options.SaltLen)
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		options.Time,
		options.Memory,
		options.Threads,
		options.KeyLen,
	)

	return encode(hash, salt, options)
}

func (hasher *Hasher) Verify(hashedPassword, candidate string) bool {
	hash, salt, options, ok := decode(hashedPassword)
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

	return subtle.ConstantTimeCompare(hash, candidateHash) == 1
}

func (hasher *Hasher) CompareOptions(hashedPassword string) bool {
	_, _, options, ok := decode(hashedPassword)
	if !ok {
		return false
	}

	return equals(options, hasher.options)
}

func encode(hash, salt []byte, options Options) string {
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

func decode(hashedPassword string) ([]byte, []byte, Options, bool) {
	var hash, salt []byte
	var options Options
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

	options = Options{
		Time:    uint32(time),
		Memory:  uint32(memory),
		Threads: uint8(threads),
		KeyLen:  uint32(keyLen),
		SaltLen: uint32(len(salt)),
	}

	return hash, salt, options, true
}

func equals(o1, o2 Options) bool {
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

func randomBytes(length uint32) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return bytes
}
