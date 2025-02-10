package repository

import (
	"context"
	_ "embed"
	"errors"
	"github.com/TicketsBot/export/internal/model"
	"github.com/jackc/pgx/v5"
)

type OAuth2Repository struct {
	tx pgx.Tx
}

var (
	//go:embed sql/oauth2/get_client.sql
	queryOauth2GetClient string

	//go:embed sql/oauth2/delete_client.sql
	queryOauth2DeleteClient string

	//go:embed sql/oauth2/update_secret.sql
	queryOauth2UpdateSecret string

	//go:embed sql/oauth2/validate_redirect_uri.sql
	queryOauth2ValidateRedirectUri string

	//go:embed sql/oauth2/create_code.sql
	queryOauth2CreateCode string

	//go:embed sql/oauth2/create_code_authority.sql
	queryOauth2CreateCodeScope string
)

func NewOAuth2Repository(tx pgx.Tx) *OAuth2Repository {
	return &OAuth2Repository{
		tx: tx,
	}
}

func (r *OAuth2Repository) GetClient(ctx context.Context, clientId string) (*model.OAuth2Client, error) {
	var client model.OAuth2Client
	if err := r.tx.QueryRow(ctx, queryOauth2GetClient, clientId).Scan(
		&client.ClientId,
		&client.ClientSecret,
		&client.OwnerId,
		&client.Label,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &client, nil
}

func (r *OAuth2Repository) DeleteClient(ctx context.Context, clientId string) error {
	_, err := r.tx.Exec(ctx, queryOauth2DeleteClient, clientId)
	return err
}

func (r *OAuth2Repository) UpdateSecret(ctx context.Context, clientId, secret string) error {
	_, err := r.tx.Exec(ctx, queryOauth2UpdateSecret, clientId, secret)
	return err
}

func (r *OAuth2Repository) ValidateRedirectUri(ctx context.Context, clientId, redirectUri string) (bool, error) {
	var valid bool
	if err := r.tx.QueryRow(ctx, queryOauth2ValidateRedirectUri, clientId, redirectUri).Scan(&valid); err != nil {
		return false, err
	}

	return valid, nil
}

func (r *OAuth2Repository) CreateCode(ctx context.Context, code model.OAuth2CodeData) error {
	_, err := r.tx.Exec(ctx, queryOauth2CreateCode, code.Code, code.ClientId, code.UserId, code.CreatedAt)
	return err
}

func (r *OAuth2Repository) CreateCodeAuthority(ctx context.Context, code, scope string) error {
	_, err := r.tx.Exec(ctx, queryOauth2CreateCodeScope, code, scope)
	return err
}
