package chi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	auth "github.com/medods-technical-assessment"
	"github.com/medods-technical-assessment/internal/common"
)

const accessTokenExpireTime = 5 * time.Minute

type AuthController struct {
	service           auth.AuthService
	validationService auth.ValidationService
	cryptoService     auth.CryptoService
	uuidService       auth.UUIDService
	jwtService        auth.JWTService
}

func NewAuthController(service auth.AuthService, validationService auth.ValidationService, cryptoService auth.CryptoService, uuidService auth.UUIDService, jwtService auth.JWTService) *AuthController {
	return &AuthController{
		service:           service,
		validationService: validationService,
		cryptoService:     cryptoService,
		uuidService:       uuidService,
		jwtService:        jwtService,
	}
}

// ref: https://stackoverflow.com/a/68100270
type CtxUUIDKey struct{}

func (c *AuthController) GetUser(w http.ResponseWriter, r *http.Request) {

	userUUID, ok := r.Context().Value(CtxUUIDKey{}).(auth.UUID)
	if !ok {
		InternalErrorHandler(w, fmt.Errorf("failed to get UUID from context"))
		return
	}
	user, err := c.service.GetUser(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(user); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := c.service.GetUsers()

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(users); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) CreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userInput auth.CreateUserDto
	if err := decoder.Decode(&userInput); err != nil {
		BadRequestErrorHandler(w, err)
		return
	}

	// Validate input
	if errors := c.validationService.ValidateUserInput(userInput); len(errors) > 0 {
		ValidationErrorHandler(w, errors)
		return
	}

	user := &auth.User{
		UUID:     c.uuidService.New(),
		Email:    userInput.Email,
		Password: c.cryptoService.HashPassword(userInput.Password),
	}

	createdUser, err := c.service.CreateUser(user)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateEmail) {
			ConflictErrorHandler(w, err)
			return
		}
		InternalErrorHandler(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(createdUser); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userUUID, ok := r.Context().Value(CtxUUIDKey{}).(auth.UUID)
	if !ok {
		InternalErrorHandler(w, fmt.Errorf("failed to get UUID from context"))
		return
	}
	user, err := c.service.GetUser(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userInput auth.UpdateUserDto
	if err := decoder.Decode(&userInput); err != nil {
		BadRequestErrorHandler(w, err)
		return
	}

	// Validate input
	if errors := c.validationService.ValidateUserInput(userInput); len(errors) > 0 {
		ValidationErrorHandler(w, errors)
		return
	}

	if userInput.Email != "" {
		user.Email = userInput.Email
	}
	if userInput.Password != "" {
		user.Password = c.cryptoService.HashPassword(userInput.Password)
	}

	updatedUser, err := c.service.UpdateUser(user)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateEmail) {
			ConflictErrorHandler(w, err)
			return
		}
		InternalErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(updatedUser); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userUUID, ok := r.Context().Value(CtxUUIDKey{}).(auth.UUID)
	if !ok {
		InternalErrorHandler(w, fmt.Errorf("failed to get UUID from context"))
		return
	}
	err := c.service.DeleteUser(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userInput auth.CreateUserDto
	if err := decoder.Decode(&userInput); err != nil {
		BadRequestErrorHandler(w, err)
		return
	}

	// Validate input
	if errors := c.validationService.ValidateUserInput(userInput); len(errors) > 0 {
		ValidationErrorHandler(w, errors)
		return
	}

	user := &auth.User{
		UUID:     c.uuidService.New(),
		Email:    userInput.Email,
		Password: c.cryptoService.HashPassword(userInput.Password),
	}

	// JWT
	issuedAt := time.Now()
	refreshPayload := &auth.RefreshPayload{Jti: c.uuidService.New()}
	accessPayload := &auth.AccessPayload{Jti: refreshPayload.Jti, IP: c.getIp(r), Sub: user.UUID, Iat: issuedAt.Unix(), Exp: issuedAt.Add(accessTokenExpireTime).Unix()}

	accessTokenStr, refreshTokenStr, err := c.jwtService.GenerateTokens(refreshPayload, accessPayload)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	_, err = c.service.CreateUser(user)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateEmail) {
			ConflictErrorHandler(w, err)
			return
		}
		InternalErrorHandler(w, err)
		return
	}

	refreshToken := &auth.RefreshToken{
		UUID:        c.uuidService.New(),
		HashedToken: c.cryptoService.HashPassword(refreshTokenStr),
		UserUUID:    user.UUID,
		Active:      true,
		CreatedAt:   issuedAt,
	}

	err = c.service.RevokeRefreshTokensByUser(refreshToken.UserUUID)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	err = c.service.AddRefreshToken(refreshToken)
	if err != nil {
		c.service.DeleteUser(user.UUID)
		InternalErrorHandler(w, err)
		return
	}

	tokens := &auth.Tokens{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tokens); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var loginInput auth.Login
	if err := decoder.Decode(&loginInput); err != nil {
		BadRequestErrorHandler(w, err)
		return
	}
	user, err := c.service.GetUserByEmail(loginInput.Email)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	if err = c.cryptoService.ComparePasswords(user.Password, loginInput.Password); err != nil {
		ForbiddenErrorHandler(w, err)
		return
	}

	// JWT
	issuedAt := time.Now()
	refreshPayload := &auth.RefreshPayload{Jti: c.uuidService.New()}
	accessPayload := &auth.AccessPayload{Jti: refreshPayload.Jti, IP: c.getIp(r), Sub: user.UUID, Iat: issuedAt.Unix(), Exp: issuedAt.Add(accessTokenExpireTime).Unix()}

	accessTokenStr, refreshTokenStr, err := c.jwtService.GenerateTokens(refreshPayload, accessPayload)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	refreshToken := &auth.RefreshToken{
		UUID:        c.uuidService.New(),
		HashedToken: c.cryptoService.HashPassword(refreshTokenStr),
		UserUUID:    user.UUID,
		Active:      true,
		CreatedAt:   issuedAt,
	}

	err = c.service.RevokeRefreshTokensByUser(refreshToken.UserUUID)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	err = c.service.AddRefreshToken(refreshToken)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	tokens := &auth.Tokens{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tokens); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var refreshInput auth.Tokens
	if err := decoder.Decode(&refreshInput); err != nil {
		BadRequestErrorHandler(w, err)
		return
	}

	accessPayload, err := c.jwtService.GetAccessTokenPayload(refreshInput.AccessToken)
	if err != nil {
		// TODO handle errors better
		BadRequestErrorHandler(w, err)
		return
	}

	user, err := c.service.GetUser(accessPayload.Sub)
	if err != nil {
		ForbiddenErrorHandler(w, err)
		return
	}

	refreshToken, err := c.service.GetActiveRefreshTokenByUser(user.UUID)
	if err != nil {
		ForbiddenErrorHandler(w, err)
		return
	}
	err = c.cryptoService.ComparePasswords(refreshToken.HashedToken, refreshInput.RefreshToken)
	if err != nil {
		ForbiddenErrorHandler(w, err)
		return
	}

	refreshPayload, err := c.jwtService.GetRefreshTokenPayload(refreshInput.RefreshToken)
	if err != nil {
		ForbiddenErrorHandler(w, err)
		return
	}

	if refreshPayload.Jti != accessPayload.Jti {
		ForbiddenErrorHandler(w, fmt.Errorf("jti in accessToken and refreshToken does not match"))
		return
	}

	// JWT
	issuedAt := time.Now()
	newRefreshPayload := &auth.RefreshPayload{Jti: c.uuidService.New()}
	newAccessPayload := &auth.AccessPayload{Jti: refreshPayload.Jti, IP: c.getIp(r), Sub: user.UUID, Iat: issuedAt.Unix(), Exp: issuedAt.Add(accessTokenExpireTime).Unix()}

	if accessPayload.IP != newAccessPayload.IP {
		log.Println("NEW IP DETECTED!", accessPayload.IP, "vs", newAccessPayload.IP)
	}

	newAccessTokenStr, newRefreshTokenStr, err := c.jwtService.GenerateTokens(newRefreshPayload, newAccessPayload)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	newRefreshToken := &auth.RefreshToken{
		UUID:        c.uuidService.New(),
		HashedToken: c.cryptoService.HashPassword(newRefreshTokenStr),
		UserUUID:    user.UUID,
		Active:      true,
		CreatedAt:   issuedAt,
	}

	err = c.service.RevokeRefreshTokensByUser(newRefreshToken.UserUUID)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	err = c.service.AddRefreshToken(newRefreshToken)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	tokens := &auth.Tokens{
		AccessToken:  newAccessTokenStr,
		RefreshToken: newRefreshTokenStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tokens); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

// TODO create middleware
// err := c.jwtService.VerifyAccessToken(refreshInput.AccessToken)
// if err != nil {
// 	ForbiddenErrorHandler(w, err)
// 	return
// }

// Returns string with either IPv4 or IPv6
func (c *AuthController) getIp(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
