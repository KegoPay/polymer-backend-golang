package cryptography

type Hasher interface {
	HashString(data string, salt []byte) ([]byte, error)
	VerifyData(hash string, data string) bool
}
