package uuid

import (
	"github.com/google/uuid"

	auth "github.com/medods-technical-assessment"
)

type UUIDService struct{}

func NewUUIDService() *UUIDService {
	return &UUIDService{}
}

func (u *UUIDService) New() auth.UUID {
	return uuid.New()
}
func (u *UUIDService) Parse(s string) (auth.UUID, error) {
	return uuid.Parse(s)
}
