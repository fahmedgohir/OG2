package game

import (
	"encoding/json"
	"fmt"
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

func (s Session) Update() bool {
	currentTime := time.Now().Unix()
	elapsed := int(currentTime - s.LastUpdated)
	if elapsed <= 0 {
		return false
	}

	newIron := elapsed * LevelToProduction(s.Factories.IronFactory.Level, Resource_Iron)
	newCopper := elapsed * LevelToProduction(s.Factories.CopperFactory.Level, Resource_Copper)
	newGold := elapsed * LevelToProduction(s.Factories.GoldFactory.Level, Resource_Gold)

	s.Resources.Iron += newIron
	s.Resources.Copper += newCopper
	s.Resources.Gold += newGold

	s.LastUpdated = currentTime
	return true
}

func (s Session) Upgrade(resource Resource) error {
	// Upgrade logic goes here
	return fmt.Errorf("upgrades not implemented")
}
