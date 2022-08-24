package game

type Resource string

const (
	Resource_Iron   Resource = "iron"
	Resource_Copper Resource = "copper"
	Resource_Gold   Resource = "gold"
)

type Factory struct {
	Level    int      `json:"level"`
	Resource Resource `json:"resource"`
}

func NewFactory(level int, resource Resource) Factory {
	return Factory{
		Level:    level,
		Resource: resource,
	}
}

type Factories struct {
	IronFactory   Factory `json:"iron_factory"`
	CopperFactory Factory `json:"copper_factory"`
	GoldFactory   Factory `json:"gold_Factory"`
}

func LevelToProduction(level int, resource Resource) int {
	switch resource {
	case Resource_Iron:
		return IronProductionRates[level]
	case Resource_Copper:
		return CopperProductionRates[level]
	case Resource_Gold:
		return GoldProductionRates[level]
	}

	return -1
}

var (
	IronProductionRates = map[int]int{
		1: 10,
		2: 20,
		3: 40,
		4: 80,
		5: 150,
	}

	CopperProductionRates = map[int]int{
		1: 3,
		2: 7,
		3: 14,
		4: 30,
		5: 60,
	}

	GoldProductionRates = map[int]int{
		1: 2,
		2: 3,
		3: 4,
		4: 6,
		5: 8,
	}
)
