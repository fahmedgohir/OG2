package game

import (
	"errors"
	"time"
)

var (
	ErrUpgradeNotEnoughResources = errors.New("not enough resources to upgrade")
	ErrUpgradeInWaitTime         = errors.New("upgrade still in wait time")
)

type Resource string

const (
	Resource_Iron   Resource = "iron"
	Resource_Copper Resource = "copper"
	Resource_Gold   Resource = "gold"
)

type Factory struct {
	Level       int      `json:"level"`
	Resource    Resource `json:"resource"`
	LastUpdated int64    `json:"last_updated"`
}

func NewFactory(level int, resource Resource) Factory {
	return Factory{
		Level:       level,
		Resource:    resource,
		LastUpdated: time.Now().Unix(),
	}
}

func (f Factory) Upgrade(resources Resources) (Factory, error) {
	var required Resources
	switch f.Resource {
	case Resource_Iron:
		required = IronUpgradeCost[f.Level]
		break
	case Resource_Copper:
		required = CopperUpgradeCost[f.Level]
		break
	case Resource_Gold:
		required = GoldUpgradeCost[f.Level]
		break
	}

	if required.Iron > resources.Iron ||
		required.Copper > resources.Copper ||
		required.Gold > resources.Gold {
		return f, ErrUpgradeNotEnoughResources
	}

	currentTime := time.Now()
	elapsed := currentTime.Sub(time.Unix(f.LastUpdated, 0))

	var waitTime time.Duration
	switch f.Resource {
	case Resource_Iron:
		waitTime = IronUpgradeDuration[f.Level]
		break
	case Resource_Copper:
		waitTime = CopperUpgradeDuration[f.Level]
		break
	case Resource_Gold:
		waitTime = GoldUpgradeDuration[f.Level]
		break
	}

	if elapsed < waitTime {
		return f, ErrUpgradeInWaitTime
	}

	f.Level += 1
	f.LastUpdated = currentTime.Unix()
	return f, nil
}

var (
	IronProductionRates = map[int]int64{
		1: 10,
		2: 20,
		3: 40,
		4: 80,
		5: 150,
	}

	CopperProductionRates = map[int]int64{
		1: 3,
		2: 7,
		3: 14,
		4: 30,
		5: 60,
	}

	GoldProductionRates = map[int]int64{
		1: 2,
		2: 3,
		3: 4,
		4: 6,
		5: 8,
	}

	IronUpgradeDuration = map[int]time.Duration{
		1: 15 * time.Second,
		2: 30 * time.Second,
		3: 60 * time.Second,
		4: 90 * time.Second,
		5: 120 * time.Second,
	}

	CopperUpgradeDuration = map[int]time.Duration{
		1: 15 * time.Second,
		2: 30 * time.Second,
		3: 60 * time.Second,
		4: 90 * time.Second,
		5: 120 * time.Second,
	}

	GoldUpgradeDuration = map[int]time.Duration{
		1: 15 * time.Second,
		2: 30 * time.Second,
		3: 60 * time.Second,
		4: 90 * time.Second,
		5: 120 * time.Second,
	}

	IronUpgradeCost = map[int]Resources{
		1: {Iron: 300, Copper: 100, Gold: 1},
		2: {Iron: 800, Copper: 250, Gold: 2},
		3: {Iron: 1600, Copper: 500, Gold: 4},
		4: {Iron: 3000, Copper: 1000, Gold: 8},
	}

	CopperUpgradeCost = map[int]Resources{
		1: {Iron: 200, Copper: 70},
		2: {Iron: 400, Copper: 150},
		3: {Iron: 800, Copper: 300},
		4: {Iron: 1600, Copper: 600},
	}

	GoldUpgradeCost = map[int]Resources{
		1: {Copper: 100, Gold: 0},
		2: {Copper: 200, Gold: 0},
		3: {Copper: 400, Gold: 0},
		4: {Copper: 800, Gold: 0},
	}
)
