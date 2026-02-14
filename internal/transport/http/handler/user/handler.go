package userapi

import (
	"context"

	"github.com/karabasBaRaBaS-1276/user-service/internal/service"
)

// PersonHandler реализует интерфейс StrictServerInterface
type PersonHandler struct {
	service service.PersonService // указатель на сервис, реализующий бизнес логику
}

var _ StrictServerInterface = (*PersonHandler)(nil)

// CreateUserFull implements [StrictServerInterface].
func (p *PersonHandler) CreateUserFull(ctx context.Context, request CreateUserFullRequestObject) (CreateUserFullResponseObject, error) {
	panic("unimplemented")
}

// DeleteUserByUsername implements [StrictServerInterface].
func (p *PersonHandler) DeleteUserByUsername(ctx context.Context, request DeleteUserByUsernameRequestObject) (DeleteUserByUsernameResponseObject, error) {
	panic("unimplemented")
}

// ExchangeAuthorizationCode implements [StrictServerInterface].
func (p *PersonHandler) ExchangeAuthorizationCode(ctx context.Context, request ExchangeAuthorizationCodeRequestObject) (ExchangeAuthorizationCodeResponseObject, error) {
	panic("unimplemented")
}

// GetOauthJwks implements [StrictServerInterface].
func (p *PersonHandler) GetOauthJwks(ctx context.Context, request GetOauthJwksRequestObject) (GetOauthJwksResponseObject, error) {
	panic("unimplemented")
}

// GetUserByUsername implements [StrictServerInterface].
func (p *PersonHandler) GetUserByUsername(ctx context.Context, request GetUserByUsernameRequestObject) (GetUserByUsernameResponseObject, error) {
	panic("unimplemented")
}

// InitiateOAuthAuthorization implements [StrictServerInterface].
func (p *PersonHandler) InitiateOAuthAuthorization(ctx context.Context, request InitiateOAuthAuthorizationRequestObject) (InitiateOAuthAuthorizationResponseObject, error) {
	panic("unimplemented")
}

// LoginUser implements [StrictServerInterface].
func (p *PersonHandler) LoginUser(ctx context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error) {
	panic("unimplemented")
}

// RestoreUserByUsername implements [StrictServerInterface].
func (p *PersonHandler) RestoreUserByUsername(ctx context.Context, request RestoreUserByUsernameRequestObject) (RestoreUserByUsernameResponseObject, error) {
	panic("unimplemented")
}

// UpdateUserByUsername implements [StrictServerInterface].
func (p *PersonHandler) UpdateUserByUsername(ctx context.Context, request UpdateUserByUsernameRequestObject) (UpdateUserByUsernameResponseObject, error) {
	panic("unimplemented")
}

func NewPersonHandler(s service.PersonService) *PersonHandler {
	return &PersonHandler{service: s}
}
