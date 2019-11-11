package model

type Suit int

const (
	Spades Suit = iota
	Clubs
	Diamonds
	Hearts
)

type Card struct {
	Suit  Suit
	Value int // Ace is 1, King is 13
}

const NumCardsPerDeck = 52

type PeggedCard struct {
	Card
	PlayerID PlayerID
}

type PlayerID int64
type GameID int64

type PlayerColor int8

const (
	Green PlayerColor = iota
	Blue
	Red
)

func (c PlayerColor) String() string {
	switch c {
	case Blue:
		return `blue`
	case Red:
		return `red`
	case Green:
		return `green`
	}
	return `notacolor`
}

type Player struct {
	ID    PlayerID
	Name  string
	Color PlayerColor // TODO map[GameID]PlayerColor
}

type BlockingPlayer struct {
	ID    PlayerID
	Reason Blocker
}

type Blocker int

const (
	DealCards Blocker = iota
	CribCard
	CutCard
	PegCard
	CountHand
	CountCrib
)

type CribBlocker struct {
	Desired int
	Dealer PlayerID
	PlayerColors map[PlayerID]PlayerColor
}

type PlayerAction struct {
	GameID GameID
	ID PlayerID
	Overcomes Blocker
	Action interface{}
}

type DealAction struct {
	NumShuffles int
}

type BuildCribAction struct {
	Cards []Card
}

type CutDeckAction struct {
	Percentage float64
}

type Phase int

const (
	Deal Phase = iota
	BuildCrib
	Cut
	Pegging
	Counting
	CribCounting
	Done
)

const (
	WinningScore int = 121
)

type Game struct {
	ID              GameID
	Players         []Player
	Deck		Deck
	BlockingPlayers []BlockingPlayer
	CurrentDealer   PlayerID
	PlayerColors    map[PlayerID]PlayerColor
	CurrentScores   map[PlayerColor]int
	LagScores       map[PlayerColor]int
	Phase           Phase
	Hands           map[PlayerID][]Card
	CutCard         Card
	Crib            []Card
	PeggedCards     []PeggedCard
}
