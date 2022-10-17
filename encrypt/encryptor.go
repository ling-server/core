package encrypt

// Encryptor encrypts or decrypts a strings
type Encryptor interface {
	// Encrypt encrypts plaintext
	Encrypt(string) (string, error)
	// Decrypt decrypts ciphertext
	Decrypt(string) (string, error)
}
