package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/export/internal/api"
	"github.com/TicketsBot/export/internal/api/constants"
	"github.com/TicketsBot/export/internal/utils"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"net/http"
	"strconv"
	"strings"
)

type AuthType string

const (
	AuthTypeUser AuthType = "user"
	AuthTypeApp  AuthType = "app"
)

func Authenticate(a *api.Core, acceptedAuthTypes ...AuthType) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers, ok := r.Header["Authorization"]
			if !ok || len(headers) == 0 {
				err := api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Missing authorization header")
				a.HandleError(r.Context(), w, err)
				return
			}

			split := strings.Split(headers[0], " ")
			if len(split) != 2 || split[0] != "Bearer" {
				err := api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid authorization scheme")
				a.HandleError(r.Context(), w, err)
				return
			}

			decoded, err := utils.Base64Decode(split[1])
			if err != nil {
				err := api.NewError(fmt.Errorf("invalid token, bad encoding: %w", err), http.StatusUnauthorized, "Invalid token")
				a.HandleError(r.Context(), w, err)
				return
			}

			tokenRaw, err := jwe.Decrypt(decoded, jwe.WithKey(jwa.DIRECT(), []byte(a.Config.Jwt.EncryptionKey)))
			if err != nil {
				err := api.NewError(fmt.Errorf("invalid token, failed to decode JWE: %w", err), http.StatusUnauthorized, "Invalid token")
				a.HandleError(r.Context(), w, err)
				return
			}

			token, err := jwt.Parse(tokenRaw, jwt.WithKey(jwa.HS256(), []byte(a.Config.Jwt.Secret)),
				jwt.WithVerify(true), jwt.WithValidate(true))
			if err != nil {
				err := api.NewError(fmt.Errorf("invalid token, failed to parse JWT: %w", err), http.StatusUnauthorized, "Invalid token")
				a.HandleError(r.Context(), w, err)
				return
			}

			authType, extractErr := extractAuthType(token)
			if extractErr != nil {
				a.HandleError(r.Context(), w, extractErr)
				return
			}

			if !utils.Contains(acceptedAuthTypes, authType) {
				err := api.NewError(nil, http.StatusUnauthorized, "Invalid token: invalid auth type")
				a.HandleError(r.Context(), w, err)
				return
			}

			subject, extractErr := extractSubject(token)
			if extractErr != nil {
				a.HandleError(r.Context(), w, extractErr)
				return
			}

			var userId uint64
			var authorizedGuilds []uint64
			if authType == AuthTypeUser {
				userId = subject

				guilds, extractErr := extractOwnedGuilds(token)
				if extractErr != nil {
					a.HandleError(r.Context(), w, extractErr)
					return
				}

				authorizedGuilds = guilds
			} else if authType == AuthTypeApp {
				onBehalfOfUser, err := extractUint64Claim(token, "on_behalf_of")
				if err != nil {
					a.HandleError(r.Context(), w, err)
					return
				}

				if onBehalfOfUser == 0 {
					err := api.NewError(nil, http.StatusUnauthorized, "Invalid token: missing on_behalf_of claim")
					a.HandleError(r.Context(), w, err)
					return
				}

				userId = onBehalfOfUser
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", userId)
			ctx = context.WithValue(ctx, "ownedGuilds", authorizedGuilds)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractAuthType(claims jwt.Token) (AuthType, *api.Error) {
	if !claims.Has(constants.JwtClaimAuthType) {
		return AuthTypeUser, nil
	}

	var authTypeRaw string
	if err := claims.Get(constants.JwtClaimAuthType, &authTypeRaw); err != nil {
		return "", api.NewError(err, http.StatusUnauthorized, "Invalid token: invalid auth type")
	}

	authType := AuthType(authTypeRaw)
	return authType, nil
}

func extractSubject(claims jwt.Token) (uint64, *api.Error) {
	return extractUint64Claim(claims, "sub")
}

func extractUint64Claim(claims jwt.Token, key string) (uint64, *api.Error) {
	var raw string
	if err := claims.Get(key, &raw); err != nil {
		return 0, api.NewError(err, http.StatusUnauthorized, "Invalid token: invalid claim")
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, api.NewError(err, http.StatusUnauthorized, "Invalid token: invalid claim")
	}

	return value, nil
}

func extractOwnedGuilds(claims jwt.Token) ([]uint64, *api.Error) {
	var guildsSlice []interface{}
	if err := claims.Get(constants.JwtClaimOwnedGuilds, &guildsSlice); err != nil {
		return nil, api.NewError(err, http.StatusUnauthorized, "Invalid token: invalid owned guilds")
	}

	ownedGuilds := make([]uint64, len(guildsSlice))
	for i, guildIdRaw := range guildsSlice {
		guild, ok := guildIdRaw.(string)
		if !ok {
			return nil, api.NewError(fmt.Errorf("invalid token, guild id was not a string"),
				http.StatusUnauthorized, "Invalid token")
		}

		guildId, err := strconv.ParseUint(guild, 10, 64)
		if err != nil {
			return nil, api.NewError(fmt.Errorf("invalid token, guild ID was not a uint: %w", err),
				http.StatusUnauthorized, "Invalid token")
		}

		ownedGuilds[i] = guildId
	}

	return ownedGuilds, nil
}
