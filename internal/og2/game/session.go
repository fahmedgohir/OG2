package game

import (
	"encoding/json"
	"time"
)

type User struct {
	Name string `json:"name"`
}

type Resources struct {
	Iron   int `json:"iron"`
	Copper int `json:"copper"`
	Gold   int `json:"gold"`
}

type Factories struct {
	IronFactory   Factory `json:"iron_factory"`
	CopperFactory Factory `json:"copper_factory"`
	GoldFactory   Factory `json:"gold_Factory"`
}

type Session struct {
	User        User      `json:"user"`
	Resources   Resources `json:"resources"`
	Factories   Factories `json:"factories"`
	LastUpdated int64     `json:"last_updated"`
}

func NewSession(user User) Session {
	return Session{
		User: user,
		Factories: Factories{
			NewFactory(1, Resource_Iron),
			NewFactory(1, Resource_Copper),
			NewFactory(1, Resource_Gold),
		},
		LastUpdated: time.Now().Unix(),
	}
}

func Marshal(s Session) ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Unmarshal(b []byte) (Session, error) {
	var session Session
	if err := json.Unmarshal(b, &session); err != nil {
		return Session{}, err
	}

	return session, nil
}

func (s Session) Update() (Session, bool) {
	currentTime := time.Now().Unix()
	elapsed := currentTime - s.LastUpdated
	if elapsed <= 0 {
		return s, false
	}

	newIron := elapsed * IronProductionRates[s.Factories.IronFactory.Level]
	newCopper := elapsed * CopperProductionRates[s.Factories.CopperFactory.Level]
	newGold := elapsed * GoldProductionRates[s.Factories.GoldFactory.Level]

	s.Resources.Iron += int(newIron)
	s.Resources.Copper += int(newCopper)
	s.Resources.Gold += int(newGold)

	s.LastUpdated = currentTime
	return s, true
}

func (s Session) Upgrade(resource Resource) (Session, error) {
	var err error
	switch resource {
	case Resource_Iron:
		s.Factories.IronFactory, err = s.Factories.IronFactory.Upgrade(s.Resources)
		break
	case Resource_Copper:
		s.Factories.CopperFactory, err = s.Factories.CopperFactory.Upgrade(s.Resources)
		break
	case Resource_Gold:
		s.Factories.GoldFactory, err = s.Factories.GoldFactory.Upgrade(s.Resources)
		break
	}

	return s, err
}
