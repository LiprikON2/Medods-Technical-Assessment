package chi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"time"

	auth "github.com/medods-technical-assessment"
	"github.com/medods-technical-assessment/internal/common"
	"github.com/medods-technical-assessment/pkg/utils"
)

const accessTokenExpireTime = 5 * time.Minute

type AuthController struct {
	service           auth.AuthService
	validationService auth.ValidationService
	cryptoService     auth.CryptoService
	uuidService       auth.UUIDService
	jwtService        auth.JWTService
	mailService       auth.MailService
}

func NewAuthController(service auth.AuthService, validationService auth.ValidationService, cryptoService auth.CryptoService, uuidService auth.UUIDService, jwtService auth.JWTService, mailService auth.MailService) *AuthController {
	return &AuthController{
		service:           service,
		validationService: validationService,
		cryptoService:     cryptoService,
		uuidService:       uuidService,
		jwtService:        jwtService,
		mailService:       mailService,
	}
}

// ref: https://stackoverflow.com/a/68100270
type CtxUserUUIDKey struct{}

type CtxAccessTokenKey struct{}

func (c *AuthController) GetUser(w http.ResponseWriter, r *http.Request) {

	userUUID, err := c.getUserUUIDFromContext(r)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}
	user, err := c.service.GetUser(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	publicUser := user.ToPublic()

	if err = c.writeResponse(respParams{w: w, code: http.StatusOK, json: publicUser}); err != nil {
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
	publicUsers := utils.MapSlice(users, func(user *auth.User) *auth.PublicUser {
		return user.ToPublic()
	})

	if err = c.writeResponse(respParams{w: w, code: http.StatusOK, json: publicUsers}); err != nil {
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

	if err = c.writeResponse(respParams{w: w, code: http.StatusCreated, json: createdUser}); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userUUID, err := c.getUserUUIDFromContext(r)
	if err != nil {
		InternalErrorHandler(w, err)
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

	if err = c.writeResponse(respParams{w: w, code: http.StatusOK, json: updatedUser}); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userUUID, err := c.getUserUUIDFromContext(r)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}
	err = c.service.DeleteUser(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	if err = c.writeResponse(respParams{w: w, code: http.StatusNoContent}); err != nil {
		InternalErrorHandler(w, err)
		return
	}

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

	refreshPayload, accessPayload := c.createPayloads(r)

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

	refreshToken := c.makeRefreshToken(refreshTokenStr, accessPayload.Jti, user.UUID, accessPayload.Iat)

	err = c.service.RevokeRefreshTokensByUser(user.UUID)
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

	if err = c.writeResponse(respParams{w: w, code: http.StatusCreated, json: tokens}); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var loginInput auth.LoginUserDto
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

	c.handleSuccessfulAuth(w, r, user)
}

func (c *AuthController) LoginByUUID(w http.ResponseWriter, r *http.Request) {

	userUUID, err := c.getUserUUIDFromContext(r)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}
	user, err := c.service.GetUser(userUUID)

	if err != nil {
		NotFoundErrorHandler(w, err)
		return
	}

	c.handleSuccessfulAuth(w, r, user)
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
		BadRequestErrorHandler(w, err)
		return
	}

	refreshToken, err := c.service.GetActiveRefreshToken(accessPayload.Jti)
	if err != nil {
		ForbiddenErrorHandler(w, err)
		return
	}

	user, err := c.service.GetUser(refreshToken.UserUUID)
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
		ForbiddenErrorHandler(w, fmt.Errorf("jti in accessToken and refreshToken do not match"))
		return
	}

	newRefreshPayload, newAccessPayload := c.createPayloads(r)

	if accessPayload.IP != newAccessPayload.IP {
		c.mailService.Send(user.Email, "New login", fmt.Sprintf(`We noticed you logged in from a new ip address %s. If this was you, there's nothing for you to do right now.`, newAccessPayload.IP))
	}

	newAccessTokenStr, newRefreshTokenStr, err := c.jwtService.GenerateTokens(newRefreshPayload, newAccessPayload)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	newRefreshToken := c.makeRefreshToken(newRefreshTokenStr, newAccessPayload.Jti, user.UUID, accessPayload.Iat)

	err = c.service.RevokeRefreshTokensByUser(user.UUID)
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

	if err = c.writeResponse(respParams{w: w, code: http.StatusOK, json: tokens}); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) GetMe(w http.ResponseWriter, r *http.Request) {
	accessTokenStr, err := c.getAccessTokenFromContext(r)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}
	accessPayload, err := c.jwtService.GetAccessTokenPayload(accessTokenStr)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	refreshToken, err := c.service.GetActiveRefreshToken(accessPayload.Jti)
	if err != nil {
		ForbiddenErrorHandler(w, err)
		return
	}

	user, err := c.service.GetUser(refreshToken.UserUUID)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	if err = c.writeResponse(respParams{w: w, code: http.StatusOK, json: user}); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

// Returns string with either IPv4 or IPv6
func (c *AuthController) getIp(r *http.Request) (string, netip.Addr) {
	ipStr, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return "", netip.Addr{}
	}

	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return "", netip.Addr{}
	}
	return ipStr, ip
}

func (c *AuthController) createPayloads(r *http.Request) (*auth.RefreshPayload, *auth.AccessPayload) {
	issuedAt := time.Now()
	ipStr, ip := c.getIp(r)
	refreshPayload := &auth.RefreshPayload{Jti: c.uuidService.New(), IP: ip}
	accessPayload := &auth.AccessPayload{Jti: refreshPayload.Jti, IP: ipStr, Iat: issuedAt.Unix(), Exp: issuedAt.Add(accessTokenExpireTime).Unix()}

	return refreshPayload, accessPayload
}

func (c *AuthController) handleSuccessfulAuth(w http.ResponseWriter, r *http.Request, user *auth.User) {
	refreshPayload, accessPayload := c.createPayloads(r)

	accessTokenStr, refreshTokenStr, err := c.jwtService.GenerateTokens(refreshPayload, accessPayload)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	refreshToken := c.makeRefreshToken(refreshTokenStr, accessPayload.Jti, user.UUID, accessPayload.Iat)

	if err = c.service.RevokeRefreshTokensByUser(user.UUID); err != nil {
		InternalErrorHandler(w, err)
		return
	}

	if err = c.service.AddRefreshToken(refreshToken); err != nil {
		InternalErrorHandler(w, err)
		return
	}

	tokens := &auth.Tokens{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}

	if err = c.writeResponse(respParams{w: w, code: http.StatusOK, json: tokens}); err != nil {
		InternalErrorHandler(w, err)
		return
	}
}

func (c *AuthController) makeRefreshToken(refreshTokenStr string, uuid auth.UUID, userUUID auth.UUID, Iat int64) *auth.RefreshToken {
	refreshToken := &auth.RefreshToken{
		UUID:        uuid,
		HashedToken: c.cryptoService.HashPassword(refreshTokenStr),
		UserUUID:    userUUID,
		Active:      true,
		CreatedAt:   time.Unix(Iat, 0),
	}

	return refreshToken
}

func (c *AuthController) getUserUUIDFromContext(r *http.Request) (auth.UUID, error) {
	userUUID, ok := r.Context().Value(CtxUserUUIDKey{}).(auth.UUID)
	if !ok {
		return auth.UUID{}, fmt.Errorf("failed to get UUID from context")
	}

	return userUUID, nil
}

func (c *AuthController) getAccessTokenFromContext(r *http.Request) (string, error) {
	accessToken, ok := r.Context().Value(CtxAccessTokenKey{}).(string)
	if !ok {
		return "", fmt.Errorf("failed to get UUID from context")
	}

	return accessToken, nil
}

type respParams struct {
	w    http.ResponseWriter
	code int
	json any
}

func (c *AuthController) writeResponse(params respParams) error {
	if params.w == nil {
		return fmt.Errorf("w param is not provided")
	}
	if params.code == 0 {
		return fmt.Errorf("statusCode param is not provided")
	}

	params.w.Header().Set("Content-Type", "application/json")
	params.w.WriteHeader(params.code)

	if params.json == nil {
		return nil
	}

	encoder := json.NewEncoder(params.w)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(params.json)
	if err != nil {
		return fmt.Errorf("error writing response: %w", err)
	}

	return nil
}
