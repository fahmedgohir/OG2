package og2

import (
	"database/sql"
	"errors"
	"hunter.io/og2/internal/og2/game"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // Register sql driver
)

var (
	// ErrCouldNotCreateTable occurs if sessions table could not be created at initialisation.
	ErrCouldNotCreateTable = errors.New("could not create sessions table")

	// ErrCouldNotAddSession occurs if a session could not be added.
	ErrCouldNotAddSession = errors.New("could not add session")

	// ErrCouldNotFindSession occurs if a session could not be found.
	ErrCouldNotFindSession = errors.New("could not find session")

	// ErrCouldNotUpdateSession occurs if a session could not be updated.
	ErrCouldNotUpdateSession = errors.New("could not update session")
)

type Sessions interface {
	Start()
	Create(user game.User) error
	Get(user game.User) (game.Session, error)
	Update() error
	Close()
}

type sessions struct {
	db    *sql.DB
	mutex sync.RWMutex
	quit  chan struct{}
}

func NewSessions(db *sql.DB) (Sessions, error) {
	query := `CREATE TABLE IF NOT EXISTS sessions(
			name TEXT PRIMARY KEY,
			iron INT,
			copper INT,
			gold INT,
			iron_level INT,
			copper_level INT,
			gold_level INT,
			last_updated BIGINT
		);`

	_, err := db.Exec(query)
	if err != nil {
		return nil, ErrCouldNotCreateTable
	}

	return &sessions{
		db: db,
	}, nil
}

func (s *sessions) Start() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.Update()
			case <-s.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *sessions) Close() {
	close(s.quit)
}

func (s *sessions) Create(user game.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := `
		INSERT INTO sessions(
			name,
			iron,
		    copper,
			gold,
			iron_level,
		    copper_level,
			gold_level,
			last_updated
		) values(?, ?, ?, ?, ?, ?, ?, ?);
	`

	res, err := s.db.Exec(query, user.Name, 0, 0, 0, 1, 1, 1, time.Now().Unix())
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows != 1 {
		return ErrCouldNotAddSession
	}

	return nil
}

func (s *sessions) Get(user game.User) (game.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	query := `
		SELECT name, iron, copper, gold, iron_level, copper_level, gold_level, last_updated FROM sessions
		WHERE name = $1
	`

	row := s.db.QueryRow(query, user.Name)
	if row == nil {
		return game.Session{}, ErrCouldNotFindSession
	}

	var name string
	var iron, copper, gold int
	var iron_level, copper_level, gold_level int
	var last_updated int64

	err := row.Scan(
		&name,
		&iron,
		&copper,
		&gold,
		&iron_level,
		&copper_level,
		&gold_level,
		&last_updated,
	)

	if err != nil {
		return game.Session{}, err
	}

	return toSession(
		name,
		iron, copper, gold,
		iron_level, copper_level, gold_level,
		last_updated,
	), nil
}

func (s *sessions) Update() error {
	query := `
		SELECT name, iron, copper, gold, iron_level, copper_level, gold_level, last_updated FROM sessions
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	sessions := make([]game.Session, 0)
	for rows.Next() {
		var name string
		var iron, copper, gold int
		var iron_level, copper_level, gold_level int
		var last_updated int64

		err = rows.Scan(
			&name,
			&iron,
			&copper,
			&gold,
			&iron_level,
			&copper_level,
			&gold_level,
			&last_updated,
		)
		if err != nil {
			return err
		}

		sessions = append(sessions, toSession(
			name,
			iron, copper, gold,
			iron_level, copper_level, gold_level,
			last_updated,
		))
	}

	for _, session := range sessions {
		updated := session.Update()
		s.update(updated)
	}

	return nil
}

func (s *sessions) update(session game.Session) error {
	query := `
		UPDATE sessions
		SET name = $1, iron = $2, copper = $3, gold = $4, iron_level = $5, copper_level = $6, gold_level = $7, last_updated = $8
		WHERE name = $1
	`

	res, err := s.db.Exec(
		query,
		session.User.Name,
		session.Resources.Iron,
		session.Resources.Copper,
		session.Resources.Gold,
		session.Factories.IronFactory.Level,
		session.Factories.CopperFactory.Level,
		session.Factories.GoldFactory.Level,
		session.LastUpdated,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows != 1 {
		return ErrCouldNotUpdateSession
	}

	return nil
}

func toSession(
	name string,
	iron, copper, gold int,
	iron_level, copper_level, gold_level int,
	last_updated int64,
) game.Session {
	return game.Session{
		User: game.User{
			Name: name,
		},
		Resources: game.Resources{
			Iron:   iron,
			Copper: copper,
			Gold:   gold,
		},
		Factories: game.Factories{
			IronFactory:   game.NewFactory(iron_level, game.Resource_Iron),
			CopperFactory: game.NewFactory(copper_level, game.Resource_Copper),
			GoldFactory:   game.NewFactory(gold_level, game.Resource_Gold),
		},
		LastUpdated: last_updated,
	}
}
