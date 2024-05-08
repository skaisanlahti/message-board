package app

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/skaisanlahti/message-board/internal/assert"
)

func SessionId() string {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	assert.Ok(err, "failed to create session id random bytes")
	return base64.URLEncoding.EncodeToString(bytes)
}
