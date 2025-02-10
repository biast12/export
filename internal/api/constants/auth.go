package constants

const (
	JwtClaimAuthType    = "auth_type"
	JwtClaimOwnedGuilds = "owned_guilds"

	ScopeGuildData        = "guild_data"
	ScopeGuildTranscripts = "guild_transcripts"
)

var Scopes = []string{
	ScopeGuildData,
	ScopeGuildTranscripts,
}
