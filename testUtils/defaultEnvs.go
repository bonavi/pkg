package testUtils

import "testing"

func SetDefaultEnvs(t *testing.T) {
	t.Setenv("ACCESS_TOKEN_TTL", "1h")
	t.Setenv("REFRESH_TOKEN_TTL", "24h")
}
