package game

import (
	"fmt"
	"time"
)

type User struct {
	Name string `json:"name"`
}

type Upgrade struct {
	User    User     `json:"user"`
	Factory Resource `json:"factory"`
}

type Resources struct {
	Iron   int `json:"iron"`
	Copper int `json:"copper"`
	Gold   int `json:"gold"`
}

type Session struct {
	User        User      `json:"user"`
	Resources   Resources `json:"resources"`
	Factories   Factories `json:"factories"`
	LastUpdated int64     `json:"last_updated"`
}

func (s Session) Update() Session {
	currentTime := time.Now().Unix()
	elapsed := int(currentTime - s.LastUpdated)

	newIron := elapsed * LevelToProduction(s.Factories.IronFactory.Level, Resource_Iron)
	newCopper := elapsed * LevelToProduction(s.Factories.CopperFactory.Level, Resource_Copper)
	newGold := elapsed * LevelToProduction(s.Factories.GoldFactory.Level, Resource_Gold)

	s.Resources.Iron += newIron
	s.Resources.Copper += newCopper
	s.Resources.Gold += newGold

	s.LastUpdated = currentTime
	return s
}

func (s Session) Upgrade(resource Resource) (Session, error) {
	// Upgrade logic goes here
	return s, fmt.Errorf("upgrades not implemented")
}
