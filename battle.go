package game

import (
	"log"
	"sync"
	"time"
)

// BattleState status pertempuran
type BattleState int

const (
	Waiting BattleState = iota
	InProgress
	Finished
)

// BattleRoom room pertempuran 5v5
type BattleRoom struct {
	ID        string
	State     BattleState
	Team1     []*Player
	Team2     []*Player
	Tower1HP  float64
	Tower2HP  float64
	Base1HP   float64
	Base2HP   float64
	StartTime time.Time
	mu        sync.RWMutex
}

// Player dalam battle
type Player struct {
	ID       string
	Name     string
	Hero     *Hero
	Team     int // 1 or 2
	HP       float64
	MP       float64
	Level    int
	Exp      float64
	Gold     float64
	Kills    int
	Deaths   int
	Assists  int
	Position Position
	Alive    bool
	RespawnTimer time.Time
}

// NewBattleRoom membuat room battle baru
func NewBattleRoom(id string) *BattleRoom {
	return &BattleRoom{
		ID:       id,
		State:    Waiting,
		Tower1HP: 5000,
		Tower2HP: 5000,
		Base1HP:  10000,
		Base2HP:  10000,
	}
}

// StartBattle memulai pertempuran
func (b *BattleRoom) StartBattle() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if len(b.Team1) == 5 && len(b.Team2) == 5 {
		b.State = InProgress
		b.StartTime = time.Now()
		log.Printf("Battle %s started! 5v5", b.ID)
		go b.gameLoop()
	}
}

// gameLoop loop utama pertempuran (tick-based)
func (b *BattleRoom) gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 ticks per second
	defer ticker.Stop()

	for b.State == InProgress {
		<-ticker.C
		b.update()
		b.checkWinCondition()
	}
}

// update update state game setiap tick
func (b *BattleRoom) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Update semua player
	allPlayers := append(b.Team1, b.Team2...)
	for _, p := range allPlayers {
		if !p.Alive {
			if time.Now().After(p.RespawnTimer) {
				p.Alive = true
				p.HP = p.Hero.BaseHP
				p.MP = p.Hero.BaseMP
				log.Printf("Player %s respawned!", p.Name)
			}
			continue
		}

		// Regen HP/MP
		if p.HP < p.Hero.BaseHP {
			p.HP += 5 // regen per tick
		}
		if p.MP < p.Hero.BaseMP {
			p.MP += 3
		}
	}
}

// checkWinCondition cek kondisi menang
func (b *BattleRoom) checkWinCondition() {
	if b.Base1HP <= 0 {
		b.State = Finished
		log.Printf("Team 2 WINS! Battle %s finished", b.ID)
	} else if b.Base2HP <= 0 {
		b.State = Finished
		log.Printf("Team 1 WINS! Battle %s finished", b.ID)
	}
}

// UseSkill player menggunakan skill
func (b *BattleRoom) UseSkill(playerID string, skillIndex int, targetID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var attacker *Player
	var target *Player

	allPlayers := append(b.Team1, b.Team2...)
	for _, p := range allPlayers {
		if p.ID == playerID {
			attacker = p
		}
		if p.ID == targetID {
			target = p
		}
	}

	if attacker == nil || target == nil {
		return fmt.Errorf("player not found")
	}

	if !attacker.Alive {
		return fmt.Errorf("player is dead")
	}

	skill := attacker.Hero.Skills[skillIndex]
	
	// Check cooldown
	if time.Since(skill.LastUsed) < skill.Cooldown {
		return fmt.Errorf("skill on cooldown")
	}

	// Check mana
	if attacker.MP < skill.ManaCost {
		return fmt.Errorf("not enough mana")
	}

	// Apply damage
	damage := CalculateDamage(attacker.Hero, target.Hero, skill)
	target.HP -= damage
	attacker.MP -= skill.ManaCost
	skill.LastUsed = time.Now()

	log.Printf("%s used %s on %s for %.0f damage!", attacker.Name, skill.Name, target.Name, damage)

	if target.HP <= 0 {
		target.HP = 0
		target.Alive = false
		target.RespawnTimer = time.Now().Add(5 * time.Second)
		attacker.Kills++
		target.Deaths++
		attacker.Gold += 300
		log.Printf("%s killed %s!", attacker.Name, target.Name)
	}

	return nil
}
