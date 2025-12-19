package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/steipete/foodoracli/internal/browserauth"
	"github.com/steipete/foodoracli/internal/foodora"
	"github.com/steipete/foodoracli/internal/version"
	"golang.org/x/term"
)

func newLoginCmd(st *state) *cobra.Command {
	var email string
	var password string
	var passwordStdin bool
	var otpMethod string
	var otp string
	var mfaToken string
	var clientSecret string
	var storeClientSecret bool
	var browser bool
	var clientID string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login via oauth2/token (email + password; optional MFA)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if st.cfg.BaseURL == "" {
				return errors.New("missing base_url (run `foodoracli config set --country HU` or similar)")
			}
			if email == "" {
				return errors.New("--email required")
			}

			if mfaToken == "" && st.cfg.PendingMfaToken != "" && strings.EqualFold(st.cfg.PendingMfaEmail, email) {
				mfaToken = st.cfg.PendingMfaToken
			}
			if !cmd.Flags().Changed("otp-method") && st.cfg.PendingMfaChannel != "" {
				otpMethod = st.cfg.PendingMfaChannel
			}

			if clientID == "" {
				clientID = strings.TrimSpace(st.cfg.OAuthClientID)
			}
			if clientID == "" {
				clientID = "android"
			}

			if clientSecret == "" {
				// Prefer cached/env; if missing, auto-fetch via Remote Config and cache.
				sec, err := st.resolveClientSecret(cmd.Context(), clientID)
				if err != nil {
					return err
				}
				clientSecret = sec.Secret
			} else if storeClientSecret {
				st.cfg.ClientSecret = clientSecret
				st.cfg.OAuthClientID = clientID
				st.markDirty()
			}

			if passwordStdin {
				b, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
				if err != nil && !errors.Is(err, io.EOF) {
					return err
				}
				password = strings.TrimSpace(string(b))
			}
			if password == "" {
				fmt.Fprint(cmd.ErrOrStderr(), "Password: ")
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Fprintln(cmd.ErrOrStderr())
				if err != nil {
					return err
				}
				password = strings.TrimSpace(string(b))
			}
			if password == "" {
				return errors.New("empty password")
			}

			ctx := cmd.Context()
			var tok foodora.AuthToken
			var mfa *foodora.MfaChallenge

			req := foodora.OAuthPasswordRequest{
				Username:     email,
				Password:     password,
				ClientSecret: clientSecret,
				ClientID:     clientID,
				OTPMethod:    otpMethod,
				OTPCode:      otp,
				MfaToken:     mfaToken,
			}

			if browser {
				bt, bmfa, sess, err := browserauth.OAuthTokenPassword(ctx, req, browserauth.PasswordOptions{
					BaseURL:   st.cfg.BaseURL,
					DeviceID:  st.cfg.DeviceID,
					Timeout:   10 * time.Minute,
					LogWriter: cmd.ErrOrStderr(),
				})
				if err != nil {
					return err
				}
				tok, mfa = bt, bmfa

				if sess.CookieHeader != "" {
					if st.cfg.CookiesByHost == nil {
						st.cfg.CookiesByHost = map[string]string{}
					}
					st.cfg.CookiesByHost[strings.ToLower(sess.Host)] = sess.CookieHeader
					st.markDirty()
				}
				if sess.UserAgent != "" {
					st.cfg.HTTPUserAgent = sess.UserAgent
					st.markDirty()
				}
			} else {
				_, cookie := st.cookieHeaderForBaseURL()
				prof := st.appHeaders()
				ua := st.cfg.HTTPUserAgent
				if ua == "" && prof.UserAgent != "" {
					ua = prof.UserAgent
				}
				if ua == "" {
					ua = "foodoracli/" + version.Version
				}

				c, err := foodora.New(foodora.Options{
					BaseURL:          st.cfg.BaseURL,
					DeviceID:         st.cfg.DeviceID,
					GlobalEntityID:   st.cfg.GlobalEntityID,
					TargetCountryISO: st.cfg.TargetCountryISO,
					UserAgent:        ua,
					CookieHeader:     cookie,
					FPAPIKey:         prof.FPAPIKey,
					AppName:          prof.AppName,
					OriginalUserAgent: func() string {
						if strings.HasPrefix(ua, "Android-app-") {
							return ua
						}
						return ""
					}(),
				})
				if err != nil {
					return err
				}

				tok, mfa, err = c.OAuthTokenPassword(ctx, req)
				if err != nil {
					return err
				}
			}

			if mfa != nil && tok.AccessToken == "" {
				st.cfg.PendingMfaToken = mfa.MfaToken
				st.cfg.PendingMfaChannel = mfa.Channel
				if mfa.Email != "" {
					st.cfg.PendingMfaEmail = mfa.Email
				} else {
					st.cfg.PendingMfaEmail = email
				}
				st.cfg.PendingMfaCreatedAt = time.Now()
				st.markDirty()

				fmt.Fprintf(cmd.ErrOrStderr(), "MFA triggered (%s). Check your %s. Retry with:\n", mfa.Channel, mfa.Channel)
				fmt.Fprintf(cmd.ErrOrStderr(), "  foodoracli login --email %s --otp-method %s --otp <CODE>\n", email, mfa.Channel)
				fmt.Fprintf(cmd.ErrOrStderr(), "rate limit reset: %ds\n", mfa.RateLimitReset)
				return nil
			}

			now := time.Now()
			st.cfg.OAuthClientID = clientID
			st.cfg.PendingMfaToken = ""
			st.cfg.PendingMfaChannel = ""
			st.cfg.PendingMfaEmail = ""
			st.cfg.PendingMfaCreatedAt = time.Time{}
			st.cfg.AccessToken = tok.AccessToken
			st.cfg.RefreshToken = tok.RefreshToken
			st.cfg.ExpiresAt = tok.ExpiresAt(now)
			st.markDirty()
			fmt.Fprintln(cmd.OutOrStdout(), "ok")
			return nil
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "account email")
	cmd.Flags().StringVar(&password, "password", "", "password (discouraged; prefer --password-stdin or prompt)")
	cmd.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read password from stdin (first line)")
	cmd.Flags().StringVar(&otpMethod, "otp-method", "sms", "OTP/MFA channel (e.g. sms, call)")
	cmd.Flags().StringVar(&mfaToken, "mfa-token", "", "MFA token from a prior mfa_triggered response")
	cmd.Flags().StringVar(&otp, "otp", "", "OTP code for MFA verification")
	cmd.Flags().StringVar(&clientID, "client-id", "", "oauth client_id (default: config or android)")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "oauth client secret (optional; otherwise auto-fetched)")
	cmd.Flags().BoolVar(&storeClientSecret, "store-client-secret", false, "persist --client-secret into config file")
	cmd.Flags().BoolVar(&browser, "browser", false, "use an interactive Playwright browser session (helps with Cloudflare)")
	return cmd
}

func newLogoutCmd(st *state) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Forget stored tokens",
		Run: func(cmd *cobra.Command, args []string) {
			st.cfg.PendingMfaToken = ""
			st.cfg.PendingMfaChannel = ""
			st.cfg.PendingMfaEmail = ""
			st.cfg.PendingMfaCreatedAt = time.Time{}
			st.cfg.AccessToken = ""
			st.cfg.RefreshToken = ""
			st.cfg.ExpiresAt = time.Time{}
			st.markDirty()
			fmt.Fprintln(cmd.OutOrStdout(), "ok")
		},
	}
}
