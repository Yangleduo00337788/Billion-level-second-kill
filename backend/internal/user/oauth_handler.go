package user

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"inference-engine/internal/config"
	"inference-engine/internal/pkg/jwt"
	"inference-engine/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type OAuthHandler struct {
	db          *gorm.DB
	cfg         *config.OAuthConfig
	frontendURL string
	httpClient  *http.Client
}

func NewOAuthHandler(db *gorm.DB, cfg *config.OAuthConfig, frontendURL string) *OAuthHandler {
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		DialContext:     (&net.Dialer{Timeout: 15 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
	}

	// Support proxy from config or environment
	proxyURL := cfg.ProxyURL
	if proxyURL == "" {
		proxyURL = getEnvProxy()
	}
	if proxyURL != "" {
		if u, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(u)
		}
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return &OAuthHandler{db: db, cfg: cfg, frontendURL: frontendURL, httpClient: client}
}

func getEnvProxy() string {
	for _, key := range []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy", "ALL_PROXY", "all_proxy"} {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Google OAuth
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	if h.cfg.Google.ClientID == "" {
		response.Error(c, response.ErrBadRequest, "Google OAuth not configured")
		return
	}
	state := generateState()
	// 将 state 存储到 cookie 中
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
	redirectURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		url.QueryEscape(h.cfg.Google.ClientID),
		url.QueryEscape(h.cfg.Google.RedirectURL),
		url.QueryEscape("openid email profile"),
		state,
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, response.ErrBadRequest, "missing code")
		return
	}

	// 校验 state 参数
	state := c.Query("state")
	savedState, _ := c.Cookie("oauth_state")
	if savedState != "" && state != savedState {
		response.Error(c, response.ErrBadRequest, "invalid state parameter")
		return
	}

	token, err := h.exchangeGoogleCode(code)
	if err != nil {
		response.Error(c, response.ErrInternal, "Google auth failed: "+err.Error())
		return
	}

	userInfo, err := h.getGoogleUserInfo(token)
	if err != nil {
		response.Error(c, response.ErrInternal, "Failed to get Google user info: "+err.Error())
		return
	}

	user, err := h.findOrCreateOAuthUser("google", userInfo.Email, userInfo.Name, userInfo.Avatar)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	jwtToken, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login?token="+jwtToken)
}

func (h *OAuthHandler) exchangeGoogleCode(code string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {h.cfg.Google.ClientID},
		"client_secret": {h.cfg.Google.ClientSecret},
		"redirect_uri":  {h.cfg.Google.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	resp, err := h.httpClient.Post("https://oauth2.googleapis.com/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func (h *OAuthHandler) getGoogleUserInfo(token string) (*oauthUserInfo, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info oauthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GitHub OAuth
func (h *OAuthHandler) GitHubLogin(c *gin.Context) {
	if h.cfg.GitHub.ClientID == "" {
		response.Error(c, response.ErrBadRequest, "GitHub OAuth not configured")
		return
	}
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s",
		url.QueryEscape(h.cfg.GitHub.ClientID),
		url.QueryEscape(h.cfg.GitHub.RedirectURL),
		url.QueryEscape("user:email"),
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *OAuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, response.ErrBadRequest, "missing code")
		return
	}

	token, err := h.exchangeGitHubCode(code)
	if err != nil {
		response.Error(c, response.ErrInternal, "GitHub auth failed: "+err.Error())
		return
	}

	userInfo, err := h.getGitHubUserInfo(token)
	if err != nil {
		response.Error(c, response.ErrInternal, "Failed to get GitHub user info: "+err.Error())
		return
	}

	user, err := h.findOrCreateOAuthUser("github", userInfo.Email, userInfo.Name, userInfo.Avatar)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	jwtToken, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login?token="+jwtToken)
}

func (h *OAuthHandler) exchangeGitHubCode(code string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {h.cfg.GitHub.ClientID},
		"client_secret": {h.cfg.GitHub.ClientSecret},
	}

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func (h *OAuthHandler) getGitHubUserInfo(token string) (*oauthUserInfo, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
		Name      string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	// GitHub may not return email in user endpoint, try /user/emails
	if info.Email == "" {
		info.Email = h.getGitHubEmail(token)
	}

	name := info.Name
	if name == "" {
		name = info.Login
	}

	return &oauthUserInfo{
		Email:  info.Email,
		Name:   name,
		Avatar: info.AvatarURL,
	}, nil
}

func (h *OAuthHandler) getGitHubEmail(token string) string {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	json.Unmarshal(body, &emails)

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	if len(emails) > 0 {
		return emails[0].Email
	}
	return ""
}

// WeChat OAuth
func (h *OAuthHandler) WeChatLogin(c *gin.Context) {
	if h.cfg.WeChat.ClientID == "" {
		response.Error(c, response.ErrBadRequest, "WeChat OAuth not configured")
		return
	}
	redirectURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=wechat#wechat_redirect",
		url.QueryEscape(h.cfg.WeChat.ClientID),
		url.QueryEscape(h.cfg.WeChat.RedirectURL),
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *OAuthHandler) WeChatCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Error(c, response.ErrBadRequest, "missing code")
		return
	}

	tokenResp, err := h.exchangeWeChatCode(code)
	if err != nil {
		response.Error(c, response.ErrInternal, "WeChat auth failed: "+err.Error())
		return
	}

	userInfo, err := h.getWeChatUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
	if err != nil {
		response.Error(c, response.ErrInternal, "Failed to get WeChat user info: "+err.Error())
		return
	}

	email := userInfo.OpenID + "@wechat.local"
	user, err := h.findOrCreateOAuthUser("wechat", email, userInfo.Nickname, userInfo.Avatar)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	jwtToken, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/login?token="+jwtToken)
}

type wechatTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	ExpiresIn    int    `json:"expires_in"`
}

func (h *OAuthHandler) exchangeWeChatCode(code string) (*wechatTokenResp, error) {
	reqURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		url.QueryEscape(h.cfg.WeChat.ClientID),
		url.QueryEscape(h.cfg.WeChat.ClientSecret),
		code,
	)

	resp, err := h.httpClient.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result wechatTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.AccessToken == "" {
		return nil, fmt.Errorf("wechat: empty access token")
	}
	return &result, nil
}

type wechatUserInfo struct {
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"headimgurl"`
}

func (h *OAuthHandler) getWeChatUserInfo(accessToken, openID string) (*wechatUserInfo, error) {
	reqURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID,
	)

	resp, err := h.httpClient.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info wechatUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Shared
type oauthUserInfo struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"picture"`
}

func (h *OAuthHandler) findOrCreateOAuthUser(provider, email, name, avatar string) (*User, error) {
	var user User
	err := h.db.Where("email = ?", email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// Create new user
		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("oauth_"+provider+"_"+time.Now().String()), bcrypt.DefaultCost)
		user = User{
			Username: name,
			Email:    email,
			Password: string(hashedPwd),
			Avatar:   avatar,
			Role:     "user",
			Level:    1,
			Status:   1,
		}
		if err := h.db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		return &user, nil
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Update avatar if empty
	if user.Avatar == "" && avatar != "" {
		user.Avatar = avatar
		h.db.Model(&user).Update("avatar", avatar)
	}

	return &user, nil
}
