package uuid

import (
	"github.com/google/uuid"
)

type UUIDService struct{}

func NewUUIDService() *UUIDService {
	return &UUIDService{}
}

func (u *UUIDService) New() [16]byte {
	return uuid.New()
}
func (u *UUIDService) Parse(s string) ([16]byte, error) {
	return uuid.Parse(s)
}
