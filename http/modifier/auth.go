package modifier

import (
	"net/http"

	"github.com/ling-server/core/errors"
	"github.com/ling-server/core/secret"
)

// Authorizer is a kind of Modifier used to authorize the requests
type Authorizer Modifier

// SecretAuthorizer authorizes the requests with the specified secret
type SecretAuthorizer struct {
	secret string
}

// NewSecretAuthorizer returns an instance of SecretAuthorizer
func NewSecretAuthorizer(secret string) *SecretAuthorizer {
	return &SecretAuthorizer{
		secret: secret,
	}
}

// Modify the request by adding secret authentication information
func (s *SecretAuthorizer) Modify(req *http.Request) error {
	if req == nil {
		return errors.New("the request is null")
	}
	err := secret.AddToRequest(req, s.secret)
	return err
}
