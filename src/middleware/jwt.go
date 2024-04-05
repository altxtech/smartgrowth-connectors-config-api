package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)
// Custom Claims
type CustomClaims struct {
	Scope string `json:"scope"`
	Sub string `json:"sub"`
}

// Implement the validator.CustomClaims interface
func (c *CustomClaims) Validate(context.Context) error {
	return nil
}
func NewCustomClaims() validator.CustomClaims{
	return &CustomClaims{}
}

// Errors
type AuthError struct {
	Error string `json:"error"`
}

func EnsureValidToken(issDomain string, identifier string) gin.HandlerFunc {
	issuerURL, err := url.Parse("https://" + issDomain + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}
	log.Println(issuerURL.String())
	log.Println(identifier)

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	keyFunc, err  := provider.KeyFunc(context.Background())
	if err != nil {
		log.Fatalf("Key func err: %v", err)
	}
	log.Println(keyFunc)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{identifier},
		validator.WithAllowedClockSkew(time.Minute),
		validator.WithCustomClaims(NewCustomClaims),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	return func(c *gin.Context) {

		// Retrieve the token
		authHeader := c.Request.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		token := parts[len(parts) -1]
		log.Println(token)

		validClaims, err := jwtValidator.ValidateToken(context.Background(), token)
		if err != nil {
			error := AuthError{fmt.Sprintf("Authorization error: %v", err)}
			c.AbortWithStatusJSON(http.StatusUnauthorized, error)
			return
		}

		// Save claims information to request context
		claims, ok := validClaims.(*validator.ValidatedClaims)
		if !ok {
			error := AuthError{fmt.Sprintf("Invalid Claims: %v", err)}
			c.AbortWithStatusJSON(http.StatusUnauthorized, error)
			return
		}
		customClaims, ok := claims.CustomClaims.(*CustomClaims)
		if !ok {
			error := AuthError{fmt.Sprintf("Invalid Claims: %v", err)}
			c.AbortWithStatusJSON(http.StatusUnauthorized, error)
			return
		}

		c.Set("scope", customClaims.Scope)
		c.Set("sub", customClaims.Sub)

		
		c.Next()
	}
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
