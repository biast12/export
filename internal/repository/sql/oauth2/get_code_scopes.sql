SELECT oauth2_code_scopes.scope
FROM oauth2_code_scopes
INNER JOIN oauth2_codes ON oauth2_code_scopes.code = oauth2_codes.code
WHERE oauth2_codes.code = $1 AND oauth2_codes.client = $2;