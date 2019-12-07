package memory

import (
	"sync"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

var _ persistence.DB = (*memory)(nil)

type memory struct {
	lock sync.Mutex

	games        map[model.GameID][]model.Game
	players      map[model.PlayerID]model.Player
	interactions map[model.PlayerID]interaction.PlayerMeans
}

func New() persistence.DB {
	return &memory{
		games:        map[model.GameID][]model.Game{},
		players:      map[model.PlayerID]model.Player{},
		interactions: map[model.PlayerID]interaction.PlayerMeans{},
	}
}

func (m *memory) GetGame(id model.GameID) (model.Game, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if games, ok := m.games[id]; ok {
		g := games[len(games)-1]
		// Persistence should never know about the Deck in a game
		// make sure that memory mimics real persistence
		g.Deck = nil
		return g, nil
	}
	return model.Game{}, persistence.ErrGameNotFound
}

func (m *memory) GetGameAction(id model.GameID, numActions uint) (model.Game, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if games, ok := m.games[id]; ok {
		if int(numActions) >= len(games) {
			return model.Game{}, persistence.ErrGameNotFound
		}
		g := games[numActions]
		// Persistence should never know about the Deck in a game
		// make sure that memory mimics real persistence
		g.Deck = nil
		return g, nil
	}
	return model.Game{}, persistence.ErrGameNotFound
}

func (m *memory) GetPlayer(id model.PlayerID) (model.Player, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if p, ok := m.players[id]; ok {
		return p, nil
	}
	return model.Player{}, persistence.ErrPlayerNotFound
}

func (m *memory) GetInteraction(id model.PlayerID) (interaction.PlayerMeans, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if i, ok := m.interactions[id]; ok {
		return i, nil
	}
	return interaction.PlayerMeans{}, persistence.ErrInteractionNotFound
}

func (m *memory) SaveGame(g model.Game) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	id := g.ID

	if len(m.games[id]) != g.NumActions() {
		return persistence.ErrGameActionsOutOfOrder
	}

	m.games[id] = append(m.games[id], g)

	return nil
}

func (m *memory) CreatePlayer(p model.Player) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	id := p.ID
	if _, ok := m.players[id]; ok {
		return persistence.ErrPlayerAlreadyExists
	}

	m.players[id] = p
	return nil
}

func (m *memory) AddPlayerColorToGame(id model.PlayerID, color model.PlayerColor, gID model.GameID) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.players[id]; !ok {
		return persistence.ErrPlayerNotFound
	}

	m.players[id].Games[gID] = color
	return nil

}

func (m *memory) SaveInteraction(pm interaction.PlayerMeans) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.interactions[pm.PlayerID] = pm
	return nil
}
