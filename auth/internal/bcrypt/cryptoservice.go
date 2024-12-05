package bcrypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type CryptoService struct{}

func NewCryptoService() *CryptoService {
	return &CryptoService{}
}

func (c *CryptoService) HashPassword(password string) string {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
func (c *CryptoService) ComparePasswords(hpass string, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}
