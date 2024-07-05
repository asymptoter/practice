package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/external/email"
	"github.com/asymptoter/practice-backend/external/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/asymptoter/practice-backend/store/auth"
	"github.com/asymptoter/practice-backend/store/user"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type handler struct {
	sql                  *sqlx.DB
	redis                redis.Service
	us                   user.Store
	googleOAuth2Config   *oauth2.Config
	facebookOAuth2Config *oauth2.Config
	authStore            auth.Store
}

func SetHttpHandler(r *gin.RouterGroup, db *sqlx.DB, redisService redis.Service, us user.Store, authStore auth.Store) {
	gc := config.Value.Server.GoogleOAuth2

	h := handler{
		sql:       db,
		redis:     redisService,
		us:        us,
		authStore: authStore,
		googleOAuth2Config: &oauth2.Config{
			ClientID:     gc.ClientID,
			ClientSecret: gc.ClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  gc.RedirectURL,
			Scopes:       gc.Scopes,
		},
	}

	r.Use(sessions.Sessions("goquestsession", sessions.NewCookieStore([]byte("secret"))))
	r.Handle("POST", "/signup", h.signup)
	r.Handle("GET", "/activation", h.activation)
	r.Handle("POST", "/login", h.login)

	r.Handle("GET", "/signupwithgoogle", h.signupWithGoogle)
	r.Handle("GET", "/loginwithgoogle", h.loginWithGoogle)
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	UserID uuid.UUID `json:"userID"`
}

func (h *handler) signup(c *gin.Context) {
	var signupInfo SignupRequest
	ctx := ctx.Background()
	if err := c.ShouldBind(&signupInfo); err != nil {
		ctx.With("params", signupInfo).Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user := &models.User{
		Email:    signupInfo.Email,
		Password: signupInfo.Password,
	}

	if _, err := h.us.Create(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	activeToken := uuid.New().String()
	activeAccountKey := "auth:active:activeToken:" + activeToken
	if err := h.redis.Set(ctx, activeAccountKey, user.Token, 24*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	cfg := config.Value.Server
	link := "http://" + cfg.Address + "/api/v1/auth/activation?id=" + user.ID.String() + "&activeToken=" + activeToken
	activeMessage := fmt.Sprintf(cfg.Email.ActivationMessage, link)

	if err := email.Send(ctx, signupInfo.Email, activeMessage); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, SignupResponse{
		UserID: user.ID,
	})
}

type ActivationResponse struct {
	Token string `json:"token"`
}

func (h *handler) activation(c *gin.Context) {
	userID := c.Query("id")
	activeToken := c.Query("activeToken")
	ctx := ctx.Background()

	activeAccountKey := "auth:active:activeToken:" + activeToken
	token := ""
	if err := h.redis.Get(ctx, activeAccountKey, &token); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := h.sql.Exec("UPDATE users SET token=$1 WHERE id=$2;", token, userID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ActivationResponse{
		Token: token,
	})
}

type loginRequest struct {
	Email    string
	Password string
}

func (h *handler) login(c *gin.Context) {
	var loginInfo loginRequest
	ctx := ctx.Background()
	if err := c.ShouldBind(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ctx = ctx.With("email", loginInfo.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user := models.User{}
	if err := h.sql.Get(&user, "SELECT ID, email, password, token, register_date FROM users WHERE email = $1 LIMIT 1;", loginInfo.Email); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx = ctx.With("userID", user.ID)

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(loginInfo.Password)); err != nil {
		ctx.Error(err)
		c.JSON(http.StatusUnauthorized, errors.New("Invalid email or password"))
		return
	}

	if len(user.Token) == 0 {
		ctx.Error("Account is not activated")
		c.JSON(http.StatusBadRequest, errors.New("Account is not activated"))
		return
	}

	userInfoKey := "user:" + user.ID.String()
	b, err := json.Marshal(user)
	if err != nil {
		ctx.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := h.redis.Set(ctx, userInfoKey, string(b), 24*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

type SignupWithGoogleResponse struct {
	Token string `json:"token"`
}

type SignupWithGoogleUserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *handler) signupWithGoogle(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	ctx := ctx.Background()
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	tok, err := h.googleOAuth2Config.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := h.googleOAuth2Config.Client(oauth2.NoContext, tok)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	googleUserInfo := &models.User{}
	if err := json.Unmarshal(data, googleUserInfo); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := h.authStore.Signup(ctx, googleUserInfo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, SignupWithGoogleResponse{
		Token: user.ID.String(),
	})
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (h *handler) loginWithGoogle(c *gin.Context) {
	state := randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + h.googleOAuth2Config.AuthCodeURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}
