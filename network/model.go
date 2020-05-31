package network

import "github.com/joshprzybyszewski/cribbage/model"

type CreateInteractionRequest struct {
	PlayerID      model.PlayerID `json:"playerID"`
	LocalhostPort string         `json:"localhost_port,omitempty"`
	NPCType       model.PlayerID `json:"npc_type,omitempty"`
}

type CreateGameRequest struct {
	PlayerIDs []model.PlayerID `json:"playerIDs"`
}

// TODO figure out the minimum info the client will need
type GameResponse struct {
	Players         []model.Player                       `json:"players"`
	PlayerColors    map[model.PlayerID]model.PlayerColor `json:"player_colors,omitempty"`
	CurrentScores   map[model.PlayerColor]int            `json:"current_scores"`
	LagScores       map[model.PlayerColor]int            `json:"lag_scores"`
	Phase           model.Phase                          `json:"phase"`
	BlockingPlayers map[model.PlayerID]model.Blocker     `json:"blocking_players,omitempty"`
	CurrentDealer   model.PlayerID                       `json:"current_dealer"`
	Hands           map[model.PlayerID][]model.Card      `json:"hands,omitempty"`
	Crib            []model.Card                         `json:"crib,omitempty"`
	CutCard         model.Card                           `json:"cut_card"`
	PeggedCards     []model.PeggedCard                   `json:"pegged_cards,omitempty"`
	Actions         []model.PlayerAction                 `json:"actions"`
}

func NewGameResponse(g model.Game) GameResponse {
	return GameResponse{
		Players:         g.Players,
		PlayerColors:    g.PlayerColors,
		CurrentScores:   g.CurrentScores,
		LagScores:       g.LagScores,
		Phase:           g.Phase,
		BlockingPlayers: g.BlockingPlayers,
		CurrentDealer:   g.CurrentDealer,
		Hands:           g.Hands,
		Crib:            g.Crib,
		CutCard:         g.CutCard,
		PeggedCards:     g.PeggedCards,
		Actions:         g.Actions,
	}
}

type CreatePlayerRequest struct {
	ID   model.PlayerID `json:"id"`
	Name string         `json:"name"`
}

type PlayerResponse struct {
	ID    model.PlayerID                     `json:"id"`
	Name  string                             `json:"name"`
	Games map[model.GameID]model.PlayerColor `json:"games"`
}
