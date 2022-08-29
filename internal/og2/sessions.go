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
	Set(session game.Session) error
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
			state TEXT
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
				s.update()
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

	session := game.NewSession(user)
	b, err := game.Marshal(session)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO sessions(
			name,
			state
		) values(?, ?);
	`

	res, err := s.db.Exec(query, user.Name, string(b))
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
		SELECT state FROM sessions
		WHERE name = $1
	`

	row := s.db.QueryRow(query, user.Name)
	if row == nil {
		return game.Session{}, ErrCouldNotFindSession
	}

	var state string
	if err := row.Scan(&state); err != nil {
		return game.Session{}, err
	}

	return game.Unmarshal([]byte(state))
}

func (s *sessions) update() error {
	s.mutex.RLock()
	query := `
		SELECT state FROM sessions
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	sessions := make([]game.Session, 0)
	for rows.Next() {
		var state string
		if err := rows.Scan(&state); err != nil {
			return err
		}

		session, err := game.Unmarshal([]byte(state))
		if err != nil {
			return err
		}

		sessions = append(sessions, session)
	}
	s.mutex.RUnlock()

	for _, session := range sessions {
		if updated, ok := session.Update(); ok {
			s.Set(updated)
		}
	}

	return nil
}

func (s *sessions) Set(session game.Session) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := `
		UPDATE sessions
		SET state = $1
		WHERE name = $2
	`

	b, err := game.Marshal(session)
	if err != nil {
		return err
	}

	res, err := s.db.Exec(query, string(b), session.User.Name)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows != 1 {
		return ErrCouldNotUpdateSession
	}

	return nil
}
