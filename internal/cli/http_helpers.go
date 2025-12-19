package cli

import (
	"net/url"
	"strings"
)

type appHeaderProfile struct {
	FPAPIKey  string
	AppName   string
	UserAgent string
}

func cookieHost(baseURL string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Hostname())
}

func (s *state) cookieHeaderForBaseURL() (host string, cookie string) {
	host = cookieHost(s.cfg.BaseURL)
	if host == "" || len(s.cfg.CookiesByHost) == 0 {
		return host, ""
	}
	return host, strings.TrimSpace(s.cfg.CookiesByHost[host])
}

func (s *state) appHeaders() appHeaderProfile {
	p := appHeaderProfile{
		FPAPIKey: "android",
	}
	if strings.EqualFold(s.cfg.TargetCountryISO, "AT") || strings.HasPrefix(strings.ToUpper(s.cfg.GlobalEntityID), "MJM_") || strings.Contains(strings.ToLower(s.cfg.BaseURL), "mj.fd-api.com") {
		p.AppName = "at.mjam"
		// From the provided at.mjam APKM (v25.3.0 / build 250300134).
		p.UserAgent = "Android-app-25.3.0(250300134)"
	}
	return p
}
