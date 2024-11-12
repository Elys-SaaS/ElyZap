package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func CheckJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
		if err != nil {
			log.Fatalf("Failed to parse the issuer url: %v", err)
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{os.Getenv("AUTH0_AUDIENCE")},
			validator.WithCustomClaims(
				func() validator.CustomClaims {
					return &CustomClaims{}
				},
			),
			validator.WithAllowedClockSkew(time.Minute),
		)
		if err != nil {
			log.Fatalf("Failed to set up the jwt validator")
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Encountered error while validating JWT: %v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"Failed to validate JWT."}`))
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)
		encounteredError := true

		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			c.SetRequest(r)
			err := next(c)
			if err != nil {
				return
			}
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Response(), c.Request())

		if encounteredError {
			err := c.JSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
			if err != nil {
				return err
			}
		}
		return nil
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

func ValidateServiceToken(token string) bool {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{"https://elysaas.eu.auth0.com/api/v2/"},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	_, err = jwtValidator.ValidateToken(context.Background(), token)
	return err == nil
}
