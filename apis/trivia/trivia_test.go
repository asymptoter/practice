package trivia

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asymptoter/practice-backend/models"
	mTrivia "github.com/asymptoter/practice-backend/store/trivia/mocks"
	mUser "github.com/asymptoter/practice-backend/store/user/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	any      = mock.AnythingOfType
	anyCTX   = any("ctx.CTX")
	anyUUID  = any("uuid.UUID")
	anyQuiz  = any("*models.Quiz")
	anyGame  = any("*models.Game")
	mockUser = &models.User{
		ID:    uuid.New(),
		Token: uuid.New(),
	}
)

type triviaSuite struct {
	suite.Suite
	router  *gin.Engine
	mUser   *mUser.Store
	mTrivia *mTrivia.Store
}

func TestTriviaSuite(t *testing.T) {
	suite.Run(t, new(triviaSuite))
}

func (s *triviaSuite) SetupSuite() {
	s.mTrivia = &mTrivia.Store{}
	s.mUser = &mUser.Store{}
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	SetHttpHandler(s.router.Group("/api/v1"), s.mTrivia, s.mUser)
}

func (s *triviaSuite) SetupTest() {
	s.mUser.On("GetByToken", anyCTX, anyUUID).Return(mockUser, nil).Once()
}

func (s *triviaSuite) TearDownTest() {

}

func (s *triviaSuite) TearDownSuite() {
	s.mTrivia.AssertExpectations(s.T())
	s.mUser.AssertExpectations(s.T())
}

func (s *triviaSuite) TestCreateQuiz() {
	quiz := &models.Quiz{}
	b, err := json.Marshal(quiz)
	s.NoError(err)
	r := bytes.NewReader(b)
	req := httptest.NewRequest("POST", "/api/v1/quiz", r)
	req.Header.Add("token", mockUser.Token.String())
	recorder := httptest.NewRecorder()
	s.mTrivia.On("CreateQuiz", anyCTX, anyQuiz).Return(nil).Once()
	s.router.ServeHTTP(recorder, req)
	s.Equal(http.StatusCreated, recorder.Code)
}

func (s *triviaSuite) TestGetQuizzes() {
	req := httptest.NewRequest("GET", "/api/v1/quizzes", nil)
	req.Header.Add("token", mockUser.Token.String())
	recorder := httptest.NewRecorder()
	s.mTrivia.On("GetQuizzes", anyCTX, mockUser.ID, "", "").Return([]*models.Quiz{}, nil).Once()
	s.router.ServeHTTP(recorder, req)
	s.Equal(http.StatusOK, recorder.Code)
}

func (s *triviaSuite) TestCreateGame() {
	game := &models.Game{}
	b, err := json.Marshal(game)
	s.NoError(err)
	r := bytes.NewReader(b)
	req := httptest.NewRequest("POST", "/api/v1/game", r)
	req.Header.Add("token", mockUser.Token.String())
	recorder := httptest.NewRecorder()
	s.mTrivia.On("CreateGame", anyCTX, anyGame).Return(nil).Once()
	s.router.ServeHTTP(recorder, req)
	s.Equal(http.StatusCreated, recorder.Code)
}

func (s *triviaSuite) TestGetGames() {
	req := httptest.NewRequest("GET", "/api/v1/games", nil)
	req.Header.Add("token", mockUser.Token.String())
	recorder := httptest.NewRecorder()
	s.mTrivia.On("GetGames", anyCTX, mockUser.ID, "").Return([]*models.Game{}, nil).Once()
	s.router.ServeHTTP(recorder, req)
	s.Equal(http.StatusOK, recorder.Code)
}
