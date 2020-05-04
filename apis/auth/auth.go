package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/email"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/asymptoter/practice-backend/store/auth"
	"github.com/asymptoter/practice-backend/store/user"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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

type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

func SetHttpHandler(r *gin.RouterGroup, db *sqlx.DB, redisService redis.Service, us user.Store, authStore auth.Store) {
	var c Credentials
	file, err := ioutil.ReadFile("../../creds.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &c)

	h := handler{
		sql:       db,
		redis:     redisService,
		us:        us,
		authStore: authStore,
		googleOAuth2Config: &oauth2.Config{
			ClientID:     c.Cid,
			ClientSecret: c.Csecret,
			RedirectURL:  "http://asymptoter-practice.nctu.me:8080/api/v1/auth/signupwithgoogle",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		/*
			facebookOAuth2Config: &oauth2.Config{
				ClientID:     c.Cid,
				ClientSecret: c.Csecret,
				RedirectURL:  "http://asymptoter-practice.nctu.me:8080/api/v1/auth/signupwithfacebook",
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint: facebook.Endpoint,
			},
		*/
	}

	r.Use(sessions.Sessions("goquestsession", sessions.NewCookieStore([]byte("secret"))))
	r.Handle("POST", "/signup", h.signup)
	r.Handle("GET", "/activation", h.activation)
	r.Handle("POST", "/login", h.login)

	r.Handle("POST", "/signupwithemail", h.signupWithEmail)
	r.Handle("GET", "/signupwithgoogle", h.signupWithGoogle)
	r.Handle("GET", "/loginwithgoogle", h.loginWithGoogle)
	//r.Handle("GET", "/signupwithfacebook", h.signupWithFacebook)
	//r.Handle("GET", "/loginwithfacebook", h.loginWithFacebook)
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
	context := ctx.Background()
	if err := c.ShouldBind(&signupInfo); err != nil {
		context.WithFields(logrus.Fields{
			"params": signupInfo,
			"error":  err,
		}).Error("c.ShouldBind failed")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user := &models.User{
		Email:    signupInfo.Email,
		Password: signupInfo.Password,
	}

	if _, err := h.us.Create(context, user); err != nil {
		context.WithFields(logrus.Fields{
			"err":   err,
			"email": user.Email,
		}).Error("signup failed at us.Create")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	activeToken := uuid.New().String()
	activeAccountKey := "auth:active:activeToken:" + activeToken
	if err := h.redis.Set(context, activeAccountKey, user.Token, 24*time.Hour); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"userID": user.ID,
		}).Error("redis.Set failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	cfg := config.Value.Server
	link := "http://" + cfg.Address + "/api/v1/auth/activation?id=" + user.ID.String() + "&activeToken=" + activeToken
	activeMessage := fmt.Sprintf(cfg.Email.ActivationMessage, link)

	if err := email.Send(context, signupInfo.Email, activeMessage); err != nil {
		context.WithField("err", err).Error("email.Send failed")
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
	context := ctx.Background()

	activeAccountKey := "auth:active:activeToken:" + activeToken
	token := ""
	if err := h.redis.Get(context, activeAccountKey, &token); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"userID": userID,
		}).Error("redis.Get failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := h.sql.Exec("UPDATE users SET token=$1 WHERE id=$2;", token, userID); err != nil {
		context.WithField("err", err).Error("sql.Update failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	context.Info("sql.Exec done.")

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
	context := ctx.Background()
	if err := c.ShouldBind(&loginInfo); err != nil {
		context.WithField("error", err).Error("c.ShouldBind failed")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	context = ctx.WithValue(context, "email", loginInfo.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		context.WithField("err", err).Error("bcrypt.GenerateFromPassword failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user := models.User{}
	if err := h.sql.Get(&user, "SELECT ID, email, password, token, register_date FROM users WHERE email = $1 LIMIT 1;", loginInfo.Email); err != nil {
		context.WithField("err", err).Error("login failed at sql.Get")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	context = ctx.WithValue(context, "userID", user.ID)

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(loginInfo.Password)); err != nil {
		context.WithField("err", err).Error("Invalid email or password")
		c.JSON(http.StatusUnauthorized, errors.New("Invalid email or password"))
		return
	}

	if len(user.Token) == 0 {
		context.Error("Account is not activated")
		c.JSON(http.StatusBadRequest, errors.New("Account is not activated"))
		return
	}

	userInfoKey := "user:" + user.ID.String()
	b, err := json.Marshal(user)
	if err != nil {
		context.WithField("err", err).Error("json.Marshal failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := h.redis.Set(context, userInfoKey, string(b), 24*time.Hour); err != nil {
		context.WithField("err", err).Error("redis.Set failed")
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
	context := ctx.Background()
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

	user, err := h.authStore.Signup(context, googleUserInfo)
	if err != nil {
		context.WithField("err", err).Error("signup failed at us.Create")
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

type SignupWithEmailRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReenteredPassword string `json:"reenterPassword"`
}

type SignupWithEmailResponse struct {
	Token string `json:"token"`
}

func (h *handler) signupWithEmail(c *gin.Context) {
	context := ctx.Background()

	signupInfo := &SignupWithEmailRequest{}
	if err := c.ShouldBind(&signupInfo); err != nil {
		context.WithField("err", err).Error("signupWithEmail failed at c.ShouldBind")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	context = ctx.WithValue(context, "signupInfo", signupInfo)

	// Check email format
	if err := checkmail.ValidateFormat(signupInfo.Email); err != nil {
		context.WithField("err", err).Error("signupWithEmail failed at checkmail.ValidateFormat")
		c.AbortWithError(http.StatusBadRequest, err)
	}

	// Check password length
	if len(signupInfo.Password) < 8 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("password should longer than 8 characters"))
		return
	}

	// Check reentered password is the same as password
	if signupInfo.Password != signupInfo.ReenteredPassword {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("reentered password is wrong"))
		return
	}

	// Store user information into db
	u := &models.User{
		Name:     signupInfo.Name,
		Email:    signupInfo.Email,
		Password: signupInfo.Password,
	}
	user, err := h.authStore.Signup(context, u)
	if err != nil {
		context.WithField("err", err).Error("signup failed at authStore.Signup")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, SignupWithEmailResponse{
		Token: user.ID.String(),
	})
}

/*
func (h *handler) signupWithFacebook(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	context := ctx.Background()
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	tok, err := h.facebookOAuth2Confing.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := h.facebookOAuth2ConfigConfig.Client(oauth2.NoContext, tok)
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

	user, err := h.authStore.Signup(context, googleUserInfo)
	if err != nil {
		context.WithField("err", err).Error("signup failed at us.Create")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, SignupWithGoogleResponse{
		Token: user.ID.String(),
	})
}
*/
