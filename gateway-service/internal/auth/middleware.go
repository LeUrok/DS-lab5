package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type ctxKey string

const UsernameKey ctxKey = "username"

var ErrUnauthorized = errors.New("unauthorized")

type Middleware struct {
	jwksURL  string
	issuer   string
	audience string
	cache    *jwk.Cache
}

func NewMiddleware(ctx context.Context, issuer, jwksURL, audience string) (*Middleware, error) {
	cache := jwk.NewCache(ctx)
	if err := cache.Register(jwksURL); err != nil {
		return nil, err
	}

	if _, err := cache.Get(ctx, jwksURL); err != nil {
		return nil, err
	}
	return &Middleware{
		jwksURL:  jwksURL,
		issuer:   issuer,
		audience: audience,
		cache:    cache,
	}, nil
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"message":"missing or invalid Authorization header"}`, http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		keyset, err := m.cache.Get(r.Context(), m.jwksURL)
		if err != nil {
			http.Error(w, `{"message":"failed to fetch JWKs"}`, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse([]byte(tokenString),
			jwt.WithKeySet(keyset),
			jwt.WithIssuer(m.issuer),
			jwt.WithAudience(m.audience),
			jwt.WithValidate(true),
		)
		if err != nil {
			http.Error(w, `{"message":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		username, ok := token.Get("sub")
		if !ok {
			http.Error(w, `{"message":"token has no sub claim"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UsernameKey, username.(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UsernameFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(UsernameKey).(string); ok {
		return v
	}
	return ""
}