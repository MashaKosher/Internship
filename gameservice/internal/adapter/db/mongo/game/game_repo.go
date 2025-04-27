package game

import "go.mongodb.org/mongo-driver/mongo"

type GameRepo struct {
	*mongo.Database
}

func New(db *mongo.Database) *GameRepo {
	return &GameRepo{db}
}

func (r *GameRepo) AddGame() error {
	return nil
}

func (r *GameRepo) GetUserGames() error {
	return nil
}
