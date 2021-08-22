package messages

import "crypto/rand"

func GetSalt() []byte {
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	return salt
}
