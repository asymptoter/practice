package game

type impl struct {
	mysql *gorm.DB
}

func NewStore(db *gorm.DB) Service {
	return &impl{
		mysql: db,
	}
}

func (s *impl) Result(context ctx.CTX, userInfo *user.User, gameResult Result) {
}
