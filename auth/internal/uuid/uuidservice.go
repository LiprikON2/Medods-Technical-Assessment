package uuid

import (
	googleuuid "github.com/google/uuid"

	auth "github.com/medods-technical-assessment"
)

type UUIDService struct{}

func NewUUIDService() *UUIDService {
	return &UUIDService{}
}

func (u *UUIDService) New() auth.UUID {
	return googleuuid.New()
}
func (u *UUIDService) Parse(s string) (auth.UUID, error) {
	return googleuuid.Parse(s)
}
func (u *UUIDService) FromBytes(b []byte) (uuid auth.UUID, err error) {
	return googleuuid.FromBytes(b)
}
