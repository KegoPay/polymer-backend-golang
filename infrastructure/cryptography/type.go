package cryptography

type Hasher interface{
	HashString(data string) ([]byte, error)
	VerifyData(hash string, data string) bool
}