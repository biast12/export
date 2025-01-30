package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/data-self-service/internal/api"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

func Authenticate(a *api.Core) func(handler http.Handler) http.Handler {
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

			token, err := jwt.Parse(split[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(a.Config.Jwt.Secret), nil
			})

			if err != nil {
				err := api.NewError(err, http.StatusUnauthorized, "Invalid token")
				a.HandleError(r.Context(), w, err)
				return
			}

			if !token.Valid {
				err := api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
				a.HandleError(r.Context(), w, err)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				err := api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
				a.HandleError(r.Context(), w, err)
				return
			}

			userId, extractErr := extractUserId(claims)
			if extractErr != nil {
				a.HandleError(r.Context(), w, extractErr)
				return
			}

			ownedGuilds, extractErr := extractOwnedGuilds(claims)
			if extractErr != nil {
				a.HandleError(r.Context(), w, extractErr)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", userId)
			ctx = context.WithValue(ctx, "ownedGuilds", ownedGuilds)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractUserId(claims jwt.MapClaims) (uint64, *api.Error) {
	userIdRaw, ok := claims["sub"]
	if !ok {
		return 0, api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token: missing subject")
	}

	userIdStr, ok := userIdRaw.(string)
	if !ok {
		return 0, api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token: invalid subject")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return 0, api.NewError(err, http.StatusUnauthorized, "Invalid token: invalid subject")
	}

	return userId, nil
}

func extractOwnedGuilds(claims jwt.MapClaims) ([]uint64, *api.Error) {
	guildsRaw, ok := claims["owned_guilds"]
	if !ok {
		return nil, api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
	}

	guildsSlice, ok := guildsRaw.([]interface{})
	if !ok {
		return nil, api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
	}

	ownedGuilds := make([]uint64, len(guildsSlice))
	for i, guildIdRaw := range guildsSlice {
		guild, ok := guildIdRaw.(string)
		if !ok {
			return nil, api.NewError(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
		}

		guildId, err := strconv.ParseUint(guild, 10, 64)
		if err != nil {
			return nil, api.NewError(err, http.StatusUnauthorized, "Invalid token")
		}

		ownedGuilds[i] = guildId
	}

	return ownedGuilds, nil
}
