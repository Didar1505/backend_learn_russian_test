package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	goauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Credentials stores google client-ids.
type Credentials struct {
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"secret"`
}

// googleInstalledCredentials mirrors the official Google OAuth client JSON structure.
type googleInstalledCredentials struct {
	Installed struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
	} `json:"installed"`
	Web struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
	} `json:"web"`
}

const (
	stateKey  = "state"
	sessionID = "ginoauth_google_session"
)

var (
	conf  *oauth2.Config
	store sessions.Store
)

func init() {
	gob.Register(goauth.Userinfo{})
}

var loginURL string

func randToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		glog.Fatalf("[Gin-OAuth] Failed to read rand: %v", err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Setup the authorization path
func Setup(redirectURL, credFile string, scopes []string, secret []byte) {
	store = cookie.NewStore(secret)

	var c Credentials
	file, err := os.ReadFile(credFile)
	if err != nil {
		glog.Fatalf("[Gin-OAuth] File error: %v", err)
	}
	if err := json.Unmarshal(file, &c); err != nil || c.ClientID == "" || c.ClientSecret == "" {
		var gc googleInstalledCredentials
		if err := json.Unmarshal(file, &gc); err != nil {
			glog.Fatalf("[Gin-OAuth] Failed to unmarshal client credentials: %v", err)
		}
		switch {
		case gc.Web.ClientID != "" && gc.Web.ClientSecret != "":
			c.ClientID = gc.Web.ClientID
			c.ClientSecret = gc.Web.ClientSecret
			if redirectURL == "" && len(gc.Web.RedirectURIs) > 0 {
				redirectURL = gc.Web.RedirectURIs[0]
			}
		case gc.Installed.ClientID != "" && gc.Installed.ClientSecret != "":
			c.ClientID = gc.Installed.ClientID
			c.ClientSecret = gc.Installed.ClientSecret
			if redirectURL == "" && len(gc.Installed.RedirectURIs) > 0 {
				redirectURL = gc.Installed.RedirectURIs[0]
			}
		default:
			glog.Fatalf("[Gin-OAuth] Failed to unmarshal client credentials: missing client_id/client_secret")
		}
	}

	conf = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

// SetupFromString accepts string values for ouath2 Configs
func SetupFromString(redirectURL, clientID string, clientSecret string, scopes []string, secret []byte) {
	store = cookie.NewStore(secret)

	conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

func Session(name string) gin.HandlerFunc {
	return sessions.Sessions(name, store)
}

func LoginHandler(ctx *gin.Context) {
	stateValue := randToken()
	session := sessions.Default(ctx)
	session.Set(stateKey, stateValue)
	session.Save()
	ctx.Writer.Write([]byte(`
	<html>
		<head>
			<title>Golang Google</title>
		</head>
	  <body>
			<a href='` + GetLoginURL(stateValue) + `'>
				<button>Login with Google!</button>
			</a>
		</body>
	</html>`))
}

func LoginRedirectHandler(ctx *gin.Context) {
	stateValue := randToken()
	session := sessions.Default(ctx)
	session.Set(stateKey, stateValue)
	if err := session.Save(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save session: %w", err))
		return
	}
	ctx.Redirect(http.StatusFound, GetLoginURL(stateValue))
}

func GetLoginURL(state string) string {
	return conf.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "select_account"))
}

func WithLoginURL(s string) error {
	s = strings.TrimSpace(s)
	url, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	loginURL = url.String()
	return nil
}

// Auth is the google authorization middleware. You can use them to protect a routergroup.
// Example:
//
//	       private.Use(google.Auth())
//	       private.GET("/", UserInfoHandler)
//	       private.GET("/api", func(ctx *gin.Context) {
//	           ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
//	       })
//
//	   // Requires google oauth pkg to be imported as `goauth "google.golang.org/api/oauth2/v2"`
//	   func UserInfoHandler(ctx *gin.Context) {
//		      var (
//		      	res goauth.Userinfo
//		      	ok  bool
//		      )
//
//		      val := ctx.MustGet("user")
//		      if res, ok = val.(goauth.Userinfo); !ok {
//		      	res = goauth.Userinfo{Name: "no user"}
//		      }
//
//		      ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": res.Email})
//	   }
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Handle the exchange code to initiate a transport.
		session := sessions.Default(ctx)

		retrievedState := session.Get(stateKey)
		if retrievedState != ctx.Query(stateKey) {
			if loginURL != "" {
				ctx.Redirect(302, loginURL)
			} else {
				ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state: %s", retrievedState))
			}
			return
		}

		tok, err := conf.Exchange(context.TODO(), ctx.Query("code"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to exchange code for oauth token: %w", err))
			return
		}

		oAuth2Service, err := goauth.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, tok)))
		if err != nil {
			glog.Errorf("[Gin-OAuth] Failed to create oauth service: %v", err)
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create oauth service: %w", err))
			return
		}

		userInfo, err := oAuth2Service.Userinfo.Get().Do()
		if err != nil {
			glog.Errorf("[Gin-OAuth] Failed to get userinfo for user: %v", err)
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get userinfo for user: %w", err))
			return
		}

		ctx.Set("user", userInfo)
		ctx.Next()
	}
}

func ClearSession(ctx *gin.Context) error {
	if _, ok := ctx.Get(sessions.DefaultKey); !ok {
		return nil
	}
	session := sessions.Default(ctx)
	session.Clear()
	session.Delete(sessionID)
	session.Delete(stateKey)
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})
	return session.Save()
}

func LogoutHandler(ctx *gin.Context) {
	_ = ClearSession(ctx)
	next := strings.TrimSpace(ctx.Query("next"))
	if next != "" {
		ctx.Redirect(http.StatusFound, next)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
