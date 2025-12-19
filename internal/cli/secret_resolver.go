package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/steipete/foodoracli/internal/firebase"
)

type resolvedSecret struct {
	Secret     string
	FromEnv    bool
	FromConfig bool
	FromFetch  bool
}

func (s *state) resolveClientSecret(ctx context.Context) (resolvedSecret, error) {
	if s.cfg.ClientSecret != "" {
		return resolvedSecret{Secret: s.cfg.ClientSecret, FromConfig: true}, nil
	}
	if v := os.Getenv("FOODORA_CLIENT_SECRET"); v != "" {
		return resolvedSecret{Secret: v, FromEnv: true}, nil
	}

	secret, err := fetchClientSecretFromRemoteConfig(ctx, strings.ToUpper(s.cfg.TargetCountryISO))
	if err != nil {
		return resolvedSecret{}, err
	}
	if secret == "" {
		return resolvedSecret{}, errors.New("fetched empty client secret")
	}

	// Cache for next run.
	s.cfg.ClientSecret = secret
	s.markDirty()
	return resolvedSecret{Secret: secret, FromFetch: true}, nil
}

func fetchClientSecretFromRemoteConfig(ctx context.Context, countryISO string) (string, error) {
	if countryISO == "" {
		countryISO = "HU"
	}

	rc := firebase.NewRemoteConfigClient(firebase.NetPincerHU)
	resp, err := rc.Fetch(ctx)
	if err != nil {
		return "", err
	}

	raw, ok := resp.Entries["client_secrets"]
	if !ok || strings.TrimSpace(raw) == "" {
		return "", errors.New("remote config key client_secrets missing/empty")
	}

	var m map[string]string
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return "", fmt.Errorf("client_secrets not JSON map: %w", err)
	}

	rawByCountry := strings.TrimSpace(m[countryISO])
	if rawByCountry == "" {
		return "", fmt.Errorf("client_secrets.%s missing/empty", countryISO)
	}

	// Newer configs: per-country JSON blob containing {android: "...", corp_android: "..."}.
	if strings.HasPrefix(rawByCountry, "{") {
		var per map[string]string
		if err := json.Unmarshal([]byte(rawByCountry), &per); err != nil {
			return "", fmt.Errorf("client_secrets.%s not JSON map: %w", countryISO, err)
		}
		return strings.TrimSpace(per["android"]), nil
	}

	// Older configs: per-country value is the secret itself.
	return rawByCountry, nil
}
