package jsonutils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/play"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

func jsonCopyGame(input model.Game) model.Game {
	output := input
	output.Deck = nil

	if len(output.PlayerColors) == 0 {
		output.PlayerColors = nil
	}
	if len(output.BlockingPlayers) == 0 {
		output.BlockingPlayers = nil
	}
	if len(output.Hands) == 0 {
		output.Hands = nil
	}
	if len(output.Crib) == 0 {
		output.Crib = nil
	}
	if len(output.PeggedCards) == 0 {
		output.PeggedCards = nil
	}

	return output
}

func TestUnmarshalGame(t *testing.T) {
	alice, bob, pAPIs := testutils.EmptyAliceAndBob()

	g5 := model.GameID(5)

	gPeg := gameAtPegging(t, alice, bob, pAPIs)

	testCases := []struct {
		msg  string
		game model.Game
	}{{
		msg: `deal`,
		game: model.Game{
			ID:              g5,
			Players:         []model.Player{alice, bob},
			Deck:            model.NewDeck(),
			BlockingPlayers: map[model.PlayerID]model.Blocker{alice.ID: model.CountCrib},
			CurrentDealer:   alice.ID,
			PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
			CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
			LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
			Phase:           model.CribCounting,
			Hands: map[model.PlayerID][]model.Card{
				alice.ID: {
					model.NewCardFromString(`7s`),
					model.NewCardFromString(`8s`),
					model.NewCardFromString(`9s`),
					model.NewCardFromString(`10s`),
				},
				bob.ID: {
					model.NewCardFromString(`7c`),
					model.NewCardFromString(`8c`),
					model.NewCardFromString(`9c`),
					model.NewCardFromString(`10c`),
				},
			},
			CutCard: model.NewCardFromString(`7h`),
			Crib: []model.Card{
				model.NewCardFromString(`7d`),
				model.NewCardFromString(`8d`),
				model.NewCardFromString(`9d`),
				model.NewCardFromString(`10d`),
			},
			PeggedCards: make([]model.PeggedCard, 0, 8),
			Actions: []model.PlayerAction{{
				GameID:    g5,
				ID:        alice.ID,
				Overcomes: model.DealCards,
				Action: model.DealAction{
					NumShuffles: 543,
				},
			}},
		},
	}, {
		msg:  `game at pegging`,
		game: gPeg,
	}}

	for _, tc := range testCases {
		gameCopy := jsonCopyGame(tc.game)

		b, err := json.Marshal(tc.game)
		require.NoError(t, err, tc.msg)

		actGame, err := UnmarshalGame(b)
		require.NoError(t, err, tc.msg)
		assert.Equal(t, gameCopy, actGame, tc.msg)
	}
}

func gameAtPegging(t *testing.T, alice, bob model.Player, pAPIs map[model.PlayerID]interaction.Player) model.Game {
	g, err := play.CreateGame([]model.Player{alice, bob}, pAPIs)
	require.NoError(t, err)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.DealCards,
		Action:    model.DealAction{NumShuffles: 10},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[alice.ID][0], g.Hands[alice.ID][1]}},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[bob.ID][0], g.Hands[bob.ID][1]}},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CutCard,
		Action:    model.CutDeckAction{Percentage: 0.314},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[bob.ID][0]},
	}, pAPIs))

	return g
}
