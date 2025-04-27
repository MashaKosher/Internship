package mongo

type (
	GameRepo interface {
		AddGame() error
		GetUserGames() error
	}
)
