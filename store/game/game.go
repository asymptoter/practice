package game

type GameInfo struct {
}

type GameResult struct {
}

type Quiz struct {
}

type Store interface {
	Result(context ctx.CTX, userInfo *user.User, gameInfo *GameInfo) (GameResult, error)
}
