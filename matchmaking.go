package match

import (
	"log"
	"sync"
	"time"
)

// PlayerRank data rank pemain
type PlayerRank struct {
	PlayerID string
	Rank     int
	MMR      float64 // Matchmaking Rating
}

// MatchmakingQueue antrian matchmaking
type MatchmakingQueue struct {
	mu       sync.Mutex
	queue    []PlayerRank
	minPlayers int
	maxPlayers int
}

// NewMatchmakingQueue membuat queue baru
func NewMatchmakingQueue(min, max int) *MatchmakingQueue {
	q := &MatchmakingQueue{
		minPlayers: min,
		maxPlayers: max,
	}
	go q.matchmakingLoop()
	return q
}

// AddPlayer menambah player ke queue
func (q *MatchmakingQueue) AddPlayer(p PlayerRank) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, p)
	log.Printf("Player %s added to queue. Queue size: %d", p.PlayerID, len(q.queue))
}

// matchmakingLoop loop untuk mencocokkan pemain
func (q *MatchmakingQueue) matchmakingLoop() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		q.mu.Lock()
		if len(q.queue) >= q.minPlayers {
			match := q.findMatch()
			if match != nil {
				log.Printf("Match found! %d players", len(match))
				// Start game room
				q.startGame(match)
			}
		}
		q.mu.Unlock()
	}
}

// findMatch mencari match yang cocok berdasarkan MMR
func (q *MatchmakingQueue) findMatch() []PlayerRank {
	if len(q.queue) < q.minPlayers {
		return nil
	}

	// Simple MMR-based matching
	var match []PlayerRank
	avgMMR := 0.0
	for _, p := range q.queue[:q.minPlayers] {
		avgMMR += p.MMR
	}
	avgMMR /= float64(q.minPlayers)

	for _, p := range q.queue {
		if len(match) >= q.maxPlayers {
			break
		}
		// Toleransi MMR ±200
		if p.MMR >= avgMMR-200 && p.MMR <= avgMMR+200 {
			match = append(match, p)
		}
	}

	if len(match) >= q.minPlayers {
		// Remove matched players from queue
		q.queue = q.queue[len(match):]
		return match
	}
	return nil
}

func (q *MatchmakingQueue) startGame(players []PlayerRank) {
	log.Printf("Starting game with players: %v", players)
	// Create battle room, assign teams, etc.
}
