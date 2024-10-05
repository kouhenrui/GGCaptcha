package GGCaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OAuthProvider string

const (
	Google OAuthProvider = "google"
	GitHub OAuthProvider = "github"
)

type OAuthClient struct {
	Provider     OAuthProvider
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// NewOAuthClient 初始化 OAuth 客户端
func NewOAuthClient(provider OAuthProvider, clientID, clientSecret, redirectURI string) *OAuthClient {
	return &OAuthClient{
		Provider:     provider,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

// GetAuthURL 获取 OAuth 授权 URL
func (client *OAuthClient) GetAuthURL(state string) string {
	switch client.Provider {
	case Google:
		return fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email&state=%s", client.ClientID, url.QueryEscape(client.RedirectURI), state)
	case GitHub:
		return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s", client.ClientID, url.QueryEscape(client.RedirectURI), state)
	default:
		return ""
	}
}

// ExchangeCode 交换授权码获取访问令牌
func (client *OAuthClient) ExchangeCode(code string) (string, error) {
	var tokenURL string
	var payload map[string]string

	switch client.Provider {
	case Google:
		tokenURL = "https://oauth2.googleapis.com/token"
		payload = map[string]string{
			"client_id":     client.ClientID,
			"client_secret": client.ClientSecret,
			"redirect_uri":  client.RedirectURI,
			"code":          code,
			"grant_type":    "authorization_code",
		}
	case GitHub:
		tokenURL = "https://github.com/login/oauth/access_token"
		payload = map[string]string{
			"client_id":     client.ClientID,
			"client_secret": client.ClientSecret,
			"code":          code,
			"redirect_uri":  client.RedirectURI,
		}
	default:
		return "", fmt.Errorf("unsupported OAuth provider")
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	// 发送 POST 请求获取访问令牌
	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(payloadBytes)) // 需要构造有效的请求体
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	// 根据不同的提供商解析访问令牌
	switch client.Provider {
	case Google:
		return tokenResponse["access_token"].(string), nil
	case GitHub:
		return tokenResponse["access_token"].(string), nil
	default:
		return "", fmt.Errorf("unsupported OAuth provider")
	}
}
