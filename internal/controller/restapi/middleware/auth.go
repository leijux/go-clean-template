package middleware

import (
	"net/http"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type errorResponse struct {
	Error string `json:"error"`
}

// jwtClaims carries the registered claims that identify the caller.
type jwtClaims struct {
	jwt.RegisteredClaims
}

// Auth returns a JWT authentication middleware built on
// github.com/gofiber/contrib/v3/jwt. The secret is taken from the config.
// After a valid token is found, its "sub" claim is stored as the caller id in
// ctx.Locals("userID"), which is what every REST handler reads.
func Auth(secret []byte) func(fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwt.SigningMethodHS256.Alg(),
			Key:    secret,
		},
		Claims: &jwtClaims{},
		SuccessHandler: func(c fiber.Ctx) error {
			if token := jwtware.FromContext(c); token != nil {
				if cl, ok := token.Claims.(*jwtClaims); ok {
					c.Locals("userID", cl.Subject)
				}
			}

			return c.Next()
		},
		ErrorHandler: func(c fiber.Ctx, _ error) error {
			return c.Status(http.StatusUnauthorized).JSON(errorResponse{Error: "invalid or expired token"})
		},
	})
}
