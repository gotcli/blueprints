package blueprints

// Pro tier: OAuth2/OIDC single sign-on, layered on top of the existing JWT auth.Manager (auth.file.go).
var SSOManager = `package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type SSOProvider struct {
	config      *oauth2.Config
	userInfoURL string
}

func NewSSOProviderFromViper() *SSOProvider {
	return &SSOProvider{
		config: &oauth2.Config{
			ClientID:     viper.GetString("SSO_CLIENT_ID"),
			ClientSecret: viper.GetString("SSO_CLIENT_SECRET"),
			RedirectURL:  viper.GetString("SSO_REDIRECT_URL"),
			Scopes:       viper.GetStringSlice("SSO_SCOPES"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  viper.GetString("SSO_AUTH_URL"),
				TokenURL: viper.GetString("SSO_TOKEN_URL"),
			},
		},
		userInfoURL: viper.GetString("SSO_USERINFO_URL"),
	}
}

func (p *SSOProvider) LoginURL(state string) string {
	return p.config.AuthCodeURL(state)
}

type SSOUserInfo struct {
	Subject string ` + "`json:\"sub\"`" + `
	Email   string ` + "`json:\"email\"`" + `
	Name    string ` + "`json:\"name\"`" + `
}

func (p *SSOProvider) Callback(ctx context.Context, code string) (*SSOUserInfo, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("sso: exchange code: %w", err)
	}

	client := p.config.Client(ctx, token)
	resp, err := client.Get(p.userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("sso: fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("sso: identity provider rejected the user info request")
	}

	var info SSOUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("sso: decode user info: %w", err)
	}
	if info.Subject == "" {
		return nil, errors.New("sso: identity provider did not return a subject")
	}

	return &info, nil
}
`

var SSOHandler = `package auth

import (
	"github.com/gofiber/fiber/v2"

	"{{ .ModuleName }}/pkg/response"
)

func SSOLogin(provider *SSOProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Redirect(provider.LoginURL(c.Query("state")))
	}
}

func SSOCallback(provider *SSOProvider, manager *Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing authorization code")
		}

		info, err := provider.Callback(c.Context(), code)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "sso authentication failed")
		}

		pair, err := manager.IssuePair(info.Subject, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to issue session")
		}

		return response.Send(c, fiber.StatusOK, pair)
	}
}
`
