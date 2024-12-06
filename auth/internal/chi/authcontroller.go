package chi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	auth "github.com/medods-technical-assessment"
	"github.com/medods-technical-assessment/internal/common"
)

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

	_, err := c.service.CreateUser(user)
	if err != nil {
		if errors.Is(err, common.ErrDuplicateEmail) {
			ConflictErrorHandler(w, err)
			return
		}
		InternalErrorHandler(w, err)
		return
	}

	// JWT
	issuedAt := time.Now()
	payload := auth.JWTPayload{IP: r.RemoteAddr, Iat: issuedAt.Unix(), Exp: issuedAt.Add(5 * time.Minute).Unix()}

	accessTokenStr, refreshTokenStr, err := c.jwtService.GenerateTokens(payload)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	refreshToken := &auth.RefreshToken{
		UUID:        c.uuidService.New(),
		TokenString: refreshTokenStr,
		UserUUID:    user.UUID,
		Revoked:     false,
		CreatedAt:   issuedAt,
	}
	err = c.service.AddRefreshTokenToWhitelist(refreshToken)
	if err != nil {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(user); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {

	issuedAt := time.Now()
	payload := auth.JWTPayload{IP: r.RemoteAddr, Iat: issuedAt.Unix(), Exp: issuedAt.Add(5 * time.Minute).Unix()}

	accessToken, err := c.jwtService.NewAccessToken(payload)

	refreshToken := &auth.RefreshToken{
		UUID:        c.uuidService.New(),
		TokenString: "",
		// UserUUID: ,
		Revoked:   false,
		CreatedAt: issuedAt,
	}
	log.Println("refreshToken", refreshToken)

	if err != nil {
		InternalErrorHandler(w, err)
		return
	}
	fmt.Println("accessToken", accessToken)
	fmt.Println("RemoteAddr", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
