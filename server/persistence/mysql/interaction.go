package mysql

import (
	"context"
	"database/sql"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
)

const (
	// Interactions stores the json-serialized InteractionMeans for a given player/mode
	createInteractionTable = `CREATE TABLE IF NOT EXISTS Interactions (
		PlayerID VARCHAR(` + maxPlayerUUIDLenStr + `),
		Mode INT,
		Means BLOB,
		PRIMARY KEY (PlayerID)
	) ENGINE = INNODB;`

	getPreferredPlayerMeans = `SELECT
		PreferredInteractionMode
	FROM Players
		WHERE PlayerID = ?
	;`

	getPlayerMeans = `SELECT
		Mode, Means
	FROM Interactions
		WHERE PlayerID = ?
	;`

	createPlayerMeans = `INSERT INTO Interactions
		(PlayerID, Mode, Means)
	VALUES
		(?, ?, ?)
	;`

	updatePlayerMeans = `INSERT INTO Interactions
		(PlayerID, Mode, Means)
	VALUES
		(?, ?, ?)
	ON DUPLICATE KEY UPDATE
		Means = ?
	;`
)

var (
	interactionCreateStmts = []string{
		createInteractionTable,
	}
)

var _ persistence.InteractionService = (*interactionService)(nil)

type interactionService struct {
	db *sql.DB
}

func getInteractionService(
	ctx context.Context,
	db *sql.DB,
) (persistence.InteractionService, error) {

	for _, createStmt := range interactionCreateStmts {
		_, err := db.ExecContext(ctx, createStmt)
		if err != nil {
			return nil, err
		}
	}

	return &interactionService{
		db: db,
	}, nil
}

func (s *interactionService) Get(id model.PlayerID) (interaction.PlayerMeans, error) {
	result := interaction.PlayerMeans{
		PlayerID: id,
	}

	r := s.db.QueryRow(getPreferredPlayerMeans, id)
	var preference int
	err := r.Scan(
		&preference,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return interaction.PlayerMeans{}, err
		}
		// This means the user doesn't have a preferred mode yet
		// let's just use a default
		preference = int(interaction.Unknown)
	}
	result.PreferredMode = interaction.Mode(preference)

	rows, err := s.db.Query(getPlayerMeans, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return interaction.PlayerMeans{}, persistence.ErrInteractionNotFound
		}
		return interaction.PlayerMeans{}, err
	}
	var serMeans []byte
	for rows.Next() {
		meansResult := interaction.Means{}
		err = rows.Scan(
			&meansResult.Mode,
			&serMeans,
		)
		if err != nil {
			return interaction.PlayerMeans{}, err
		}

		err = meansResult.AddSerializedInfo(serMeans)
		if err != nil {
			return interaction.PlayerMeans{}, err
		}
		result.Interactions = append(result.Interactions, meansResult)
	}

	if err := rows.Err(); err != nil {
		return interaction.PlayerMeans{}, err
	}

	return result, nil
}

func (s *interactionService) Create(pm interaction.PlayerMeans) error {
	var serMeans []byte
	var err error
	for _, means := range pm.Interactions {
		serMeans, err = means.GetSerializedInfo()
		if err != nil {
			return err
		}
		_, err = s.db.Exec(
			createPlayerMeans,
			pm.PlayerID,
			means.Mode,
			serMeans,
		)
		err = convertMysqlError(err)
		if err != nil {
			if err == errDuplicateEntry {
				return persistence.ErrPlayerAlreadyExists
			}
			return err
		}
	}
	return nil
}

func (s *interactionService) Update(pm interaction.PlayerMeans) error {
	var serMeans []byte
	var err error
	for _, means := range pm.Interactions {
		serMeans, err = means.GetSerializedInfo()
		if err != nil {
			return err
		}
		_, err = s.db.Exec(
			updatePlayerMeans,
			pm.PlayerID,
			means.Mode,
			serMeans,
			serMeans,
		)
		err = convertMysqlError(err)
		if err != nil {
			return err
		}
	}
	return nil
}