package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	// Games stores the state of a game at a given time.
	//   Each Action will update a games state and we keep a full history of all actions.
	// The columns act as follows:
	// GameID is a UUID to identify a game
	// NumActions is how many actions have occurred in the game before this one
	// ScoreBlue, ScoreRed, and ScoreGreen are the scores for each color
	// ScoreBlueLag, ScoreRedLag, and ScoreGreenLag are the previous scores for each color
	// Phase is the model.Phase that the game is currently in
	// CutCard is a number representation of the card that's been cut
	// Crib is a number representation of the (up to 4) cards in the crib
	// CurrentDealer is the PlayerID for the dealer
	// BlockingPlayers is a json encoded map of who's blocking and why
	// Hands is a json encoded map of slices for player hands
	// PeggedCards is the json-encoded slice of previously pegged cards
	// Action is the json encoded model.PlayerAction
	createGameTable = `CREATE TABLE IF NOT EXISTS Games (
		GameID INT(1) UNSIGNED,
		NumActions INT(1) UNSIGNED,
		ScoreBlue TINYINT(1) UNSIGNED,
		ScoreRed TINYINT(1) UNSIGNED,
		ScoreGreen TINYINT(1) UNSIGNED,
		ScoreBlueLag TINYINT(1) UNSIGNED,
		ScoreRedLag TINYINT(1) UNSIGNED,
		ScoreGreenLag TINYINT(1) UNSIGNED,
		Phase TINYINT(1) UNSIGNED,
		CutCard TINYINT(1) UNSIGNED,
		Crib TINYINT(4) UNSIGNED,
		CurrentDealer VARCHAR(` + maxPlayerUUIDLenStr + `),
		BlockingPlayers BLOB,
		Hands BLOB,
		PeggedCards BLOB,
		Action BLOB,
		PRIMARY KEY (GameID, NumActions)
	) ENGINE = INNODB;`

	createGamePlayersTable = `CREATE TABLE IF NOT EXISTS GamePlayers (
		GameID INT(1) UNSIGNED,
		Player1ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Player2ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Player3ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Player4ID VARCHAR(` + maxPlayerUUIDLenStr + `),
		PRIMARY KEY (GameID)
	) ENGINE = INNODB;`

	queryLatestGame = `SELECT 
		gp.Player1ID, gp.Player2ID, gp.Player3ID, gp.Player4ID,
		g.ScoreBlue, g.ScoreRed, g.ScoreGreen,
		g.ScoreBlueLag, g.ScoreRedLag, g.ScoreGreenLag,
		g.Phase, g.BlockingPlayers, g.CurrentDealer,
		g.Hands, g.Crib, g.CutCard,
		g.PeggedCards,
		g.NumActions, g.Action
	FROM Games g
	INNER JOIN GamePlayers gp
		ON g.GameID = gp.GameID
		WHERE g.GameID = ? 
	ORDER BY
		NumActions DESC
	LIMIT 1;`

	queryGameAtNumActions = `SELECT 
		g.GameID,
		gp.Player1ID, gp.Player2ID, gp.Player3ID, gp.Player4ID,
		g.ScoreBlue, g.ScoreRed, g.ScoreGreen,
		g.ScoreBlueLag, g.ScoreRedLag, g.ScoreGreenLag,
		g.Phase, g.BlockingPlayers, g.CurrentDealer,
		g.Hands, g.Crib, g.CutCard,
		g.PeggedCards,
		g.NumActions, g.Action
	FROM Games g
	INNER JOIN GamePlayers gp
		ON g.GameID = gp.GameID
	WHERE g.GameID = ? AND
		g.NumActions = ?
	;`

	queryPlayerActionsBefore = `SELECT 
		NumActions, Action
	FROM Games
	WHERE GameID = ? AND
		NumActions <= ?
	;`

	insertGameAt = `INSERT INTO Games
		(
			GameID, NumActions, 
			ScoreBlue, ScoreRed, ScoreGreen,
			ScoreBlueLag, ScoreRedLag, ScoreGreenLag,
			Phase, CutCard, Crib,
			CurrentDealer,
			BlockingPlayers, Hands, PeggedCards, Action
		)
	VALUES
		(
			?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?,
			?, ?, ?, ?
		)
	`
)

var (
	gamesCreateStmts = []string{
		createGameTable,
		createGamePlayersTable,
	}
)

var _ persistence.GameService = (*gameService)(nil)

type gameService struct {
	db *sql.DB
}

func getGameService(
	ctx context.Context,
	db *sql.DB,
) (persistence.GameService, error) {

	for _, createStmt := range gamesCreateStmts {
		_, err := db.ExecContext(ctx, createStmt)
		if err != nil {
			return nil, err
		}
	}

	return &gameService{
		db: db,
	}, nil
}

func (g *gameService) Get(id model.GameID) (model.Game, error) {
	r := g.db.QueryRow(queryLatestGame, id)
	return g.populateGameFromRow(id, r)
}

func (g *gameService) GetAt(id model.GameID, numActions uint) (model.Game, error) {
	r := g.db.QueryRow(queryGameAtNumActions, id, numActions)
	return g.populateGameFromRow(id, r)
}

func (g *gameService) populateGameFromRow(
	gID model.GameID,
	r *sql.Row,
) (model.Game, error) {

	var p1ID, p2ID, p3ID, p4ID, curDealerID model.PlayerID
	var scoreBlue, scoreRed, scoreGreen,
		lagScoreBlue, lagScoreRed, lagScoreGreen int
	var phase model.Phase
	var cribCardInts []int8 = make([]int8, 4)
	var cutCardInt int8
	var blockingPlayers, hands, peggedCards, action []byte
	var numActions uint32
	err := r.Scan(
		&p1ID, &p2ID, &p3ID, &p4ID,
		&scoreBlue, &scoreRed, &scoreGreen,
		&lagScoreBlue, &lagScoreRed, &lagScoreGreen,
		&phase, &blockingPlayers, &curDealerID,
		&hands, &cribCardInts, &cutCardInt,
		&peggedCards,
		&numActions, &action,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Game{}, persistence.ErrGameNotFound
		}
		return model.Game{}, err
	}

	curScores, lagScores := populateScores(
		scoreBlue, scoreRed, scoreGreen,
		lagScoreBlue, lagScoreRed, lagScoreGreen,
	)

	players, err := g.getPlayersForGame(p1ID, p2ID, p3ID, p4ID)
	if err != nil {
		return model.Game{}, err
	}

	pc, err := g.getPlayerColors(gID)
	if err != nil {
		return model.Game{}, err
	}

	cutCard, err := model.NewCardFromTinyInt(cutCardInt)
	if err != nil {
		// If we've errored here, just ignore it and continue
		fmt.Printf("errored card for cut: %+v\n", err)
	}

	cribCards := getCribCards(cribCardInts)

	bp, err := getBlockingPlayers(blockingPlayers)
	if err != nil {
		return model.Game{}, err
	}

	h, err := getHands(hands)
	if err != nil {
		return model.Game{}, err
	}

	p, err := getPeggedCards(peggedCards)
	if err != nil {
		return model.Game{}, err
	}

	pas, err := g.getActions(gID, int(numActions))
	if err != nil {
		return model.Game{}, err
	}

	game := model.Game{
		ID:              gID,
		CurrentScores:   curScores,
		LagScores:       lagScores,
		Players:         players,
		PlayerColors:    pc,
		Phase:           phase,
		CurrentDealer:   curDealerID,
		CutCard:         cutCard,
		Crib:            cribCards,
		BlockingPlayers: bp,
		Hands:           h,
		PeggedCards:     p,
		Actions:         pas,
	}

	return game, nil
}

func populateScores(
	scoreBlue, scoreRed, scoreGreen,
	lagScoreBlue, lagScoreRed, lagScoreGreen int,
) (cur, lag map[model.PlayerColor]int) {
	curScores := make(map[model.PlayerColor]int, 3)
	lagScores := make(map[model.PlayerColor]int, 3)
	if scoreBlue > 0 {
		curScores[model.Blue] = scoreBlue
		lagScores[model.Blue] = lagScoreBlue
	}
	if scoreRed > 0 {
		curScores[model.Red] = scoreRed
		lagScores[model.Red] = lagScoreRed
	}
	if scoreGreen > 0 {
		curScores[model.Green] = scoreGreen
		lagScores[model.Green] = lagScoreGreen
	}

	return curScores, lagScores
}

func getCribCards(cribCardInts []int8) []model.Card {
	var cribCards []model.Card
	for _, cci := range cribCardInts {
		c, err := model.NewCardFromTinyInt(cci)
		if err != nil {
			// If we've errored here, just ignore it and continue
			fmt.Printf("errored card while building crib: %+v\n", err)
			continue
		}
		cribCards = append(cribCards, c)
	}
	return cribCards
}

func getBlockingPlayers(ser []byte) (map[model.PlayerID]model.Blocker, error) {
	blockers := map[model.PlayerID]model.Blocker{}

	err := json.Unmarshal(ser, &blockers)
	if err != nil {
		return nil, err
	}

	return blockers, nil
}

func getHands(ser []byte) (map[model.PlayerID][]model.Card, error) {
	hands := map[model.PlayerID][]model.Card{}

	err := json.Unmarshal(ser, &hands)
	if err != nil {
		return nil, err
	}

	return hands, nil
}

func getPeggedCards(ser []byte) ([]model.PeggedCard, error) {
	peggedCards := []model.PeggedCard{}

	err := json.Unmarshal(ser, &peggedCards)
	if err != nil {
		return nil, err
	}

	return peggedCards, nil
}

func (g *gameService) getPlayersForGame(
	p1ID, p2ID, p3ID, p4ID model.PlayerID,
) ([]model.Player, error) {

	if len(p1ID) == 0 {
		return nil, errors.New(`at least one player required`)
	}
	if len(p2ID) == 0 {
		return nil, errors.New(`at least two players required`)
	}

	pIDs := []model.PlayerID{
		p1ID, p2ID,
	}

	if len(p3ID) > 0 {
		// The third and fourth players can only exist if the first two do
		pIDs = append(pIDs, p3ID)
		if len(p4ID) > 0 {
			pIDs = append(pIDs, p4ID)
		}
	}

	players := make([]model.Player, len(pIDs))
	for i, pID := range pIDs {
		// TODO determine if it's worth getting the entire player here, not just the ID
		players[i].ID = pID
	}
	return players, nil

}

func (g *gameService) getPlayerColors(
	gID model.GameID,
) (map[model.PlayerID]model.PlayerColor, error) {

	// populate pc with the colors for each player
	pc := make(map[model.PlayerID]model.PlayerColor, 4)

	rows, err := g.db.Query(getPlayerColorsForGame, gID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var pID model.PlayerID
		var color model.PlayerColor
		err := rows.Scan(&pID, &color)
		if err != nil {
			return nil, err
		}
		if color != model.UnsetColor {
			pc[pID] = color
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pc, nil
}

func (g *gameService) getActions(
	gID model.GameID,
	maxNumActions int,
) ([]model.PlayerAction, error) {

	rows, err := g.db.Query(queryPlayerActionsBefore, gID, maxNumActions)
	if err != nil {
		return nil, err
	}
	paMap := make(map[int][]byte, maxNumActions)
	for rows.Next() {
		var numActions int
		var action []byte
		err = rows.Scan(&numActions, &action)
		if err != nil {
			return nil, err
		}
		paMap[numActions] = action
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	pas := make([]model.PlayerAction, maxNumActions)
	for i := range pas {
		bytes, ok := paMap[i]
		if !ok {
			return nil, errors.New(`missing action`)
		}
		err = json.Unmarshal(bytes, &pas[i])
		if err != nil {
			return nil, err
		}
	}

	return pas, nil
}

func (g *gameService) UpdatePlayerColor(id model.GameID, pID model.PlayerID, color model.PlayerColor) error {
	// There should be nothing to do here because the player service should take care
	// of all of the persistence that needs to happen
	return nil
}

func (g *gameService) Save(mg model.Game) error {
	if mg.NumActions() == 0 {
		// add all of the players to be recognized in this game
		// if it's the first time we've saved it
		for _, p := range mg.Players {
			_, err := g.db.Exec(addPlayerToGame, mg.ID, p.ID)
			if err != nil {
				return err
			}
		}
	}

	// TODO ensure the player colors are accurate from mg.PlayerColors

	cut := mg.CutCard.ToTinyInt()
	crib := make([]int8, len(mg.Crib))
	for i, cc := range mg.Crib {
		crib[i] = cc.ToTinyInt()
	}
	var bp, h, pc, a []byte

	ifs := []interface{}{
		mg.ID, mg.NumActions(),
		mg.CurrentScores[model.Blue], mg.CurrentScores[model.Red], mg.CurrentScores[model.Green],
		mg.LagScores[model.Blue], mg.LagScores[model.Red], mg.LagScores[model.Green],
		mg.Phase, cut, crib,
		mg.CurrentDealer,
		bp, h, pc, a,
	}
	_, err := g.db.Exec(insertGameAt, ifs...)
	if err != nil {
		return err
	}

	return nil
}
