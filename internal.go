package game

import "time"

// HeroType tipe hero
type HeroType string

const (
	Tank     HeroType = "tank"
	Fighter  HeroType = "fighter"
	Assassin HeroType = "assassin"
	Mage     HeroType = "mage"
	Marksman HeroType = "marksman"
	Support  HeroType = "support"
)

// Skill kemampuan hero
type Skill struct {
	Name        string
	Damage      float64
	Cooldown    time.Duration
	ManaCost    float64
	Range       float64
	LastUsed    time.Time
}

// Hero data hero
type Hero struct {
	ID        string
	Name      string
	Type      HeroType
	BaseHP    float64
	BaseMP    float64
	BaseAD    float64 // Attack Damage
	BaseAP    float64 // Ability Power
	BaseArmor float64
	BaseMR    float64 // Magic Resist
	BaseSpeed float64 // Movement speed
	Skills    []Skill
}

// HeroDB database hero (contoh)
var HeroDB = map[string]*Hero{
	"li_bai": {
		ID:        "li_bai",
		Name:      "Li Bai",
		Type:      Assassin,
		BaseHP:    2800,
		BaseMP:    450,
		BaseAD:    170,
		BaseAP:    0,
		BaseArmor: 80,
		BaseMR:    50,
		BaseSpeed: 370,
		Skills: []Skill{
			{Name: "Jianghu Ershi", Damage: 300, Cooldown: 8 * time.Second, ManaCost: 60, Range: 500},
			{Name: "Shenjiang Xing", Damage: 500, Cooldown: 12 * time.Second, ManaCost: 80, Range: 800},
			{Name: "Qinglian Jian Ge", Damage: 800, Cooldown: 30 * time.Second, ManaCost: 100, Range: 1000},
		},
	},
	"diaochan": {
		ID:        "diaochan",
		Name:      "Diaochan",
		Type:      Mage,
		BaseHP:    2600,
		BaseMP:    500,
		BaseAD:    150,
		BaseAP:    120,
		BaseArmor: 60,
		BaseMR:    40,
		BaseSpeed: 350,
		Skills: []Skill{
			{Name: "Zhen Hua", Damage: 250, Cooldown: 5 * time.Second, ManaCost: 50, Range: 600},
			{Name: "Yuan Wu", Damage: 400, Cooldown: 10 * time.Second, ManaCost: 70, Range: 700},
			{Name: "Sheng Kai", Damage: 600, Cooldown: 25 * time.Second, ManaCost: 90, Range: 900},
		},
	},
}

// GetHero mendapatkan hero by ID
func GetHero(id string) *Hero {
	return HeroDB[id]
}

// CalculateDamage menghitung damage
func CalculateDamage(attacker *Hero, defender *Hero, skill Skill) float64 {
	baseDamage := skill.Damage + attacker.BaseAD*0.5
	
	// Reduksi armor
	armorReduction := defender.BaseArmor / (defender.BaseArmor + 100)
	finalDamage := baseDamage * (1 - armorReduction)
	
	return finalDamage
}
