package persistence

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

type DB interface {
	CreatePlayer(p model.Player) error
	GetPlayer(id model.PlayerID) (model.Player, error)
	AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error

	GetGame(id model.GameID) (model.Game, error)
	GetGameAction(id model.GameID, numActions uint) (model.Game, error)
	SaveGame(g model.Game) error

	GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error)
	SaveInteraction(pm interaction.PlayerMeans) error
}

var (
	ErrPlayerNotFound      error = errors.New(`player not found`)
	ErrPlayerAlreadyExists error = errors.New(`username already exists`)

	ErrGameNotFound          error = errors.New(`game not found`)
	ErrGameInitialSave       error = errors.New(`game must be saved with no actions`)
	ErrGameActionsOutOfOrder error = errors.New(`game actions out of order`)

	ErrInteractionNotFound error = errors.New(`interaction not found`)
)
