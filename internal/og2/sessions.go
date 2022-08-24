package og2

type User struct {
	Name string `json:"name"`
}

type Session struct {
	User User `json:"user"`
}

type Sessions struct {
	sessions []*Session
}

func NewSessions() *Sessions {
	return &Sessions{
		sessions: make([]*Session, 0),
	}
}

func (s *Sessions) Create(user User) *Session {
	session := &Session{
		User: user,
	}
	s.sessions = append(s.sessions, session)

	return session
}
