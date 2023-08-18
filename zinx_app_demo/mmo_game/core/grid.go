package core

import (
	"fmt"
	"sync"
)

// GrID A grid class in a map 一个地图中的格子类
type GrID struct {
	GID       int          // Grid ID
	MinX      int          // Left boundary coordinate of the grid
	MaxX      int          // Right boundary coordinate of the grid
	MinY      int          // Upper boundary coordinate of the grid
	MaxY      int          // Lower boundary coordinate of the grid
	playerIDs map[int]bool // IDs of players or objects in the current grid
	pIDLock   sync.RWMutex // Lock for protecting the playerIDs map
}

// NewGrID Initialize a grid
func NewGrID(gID, minX, maxX, minY, maxY int) *GrID {
	return &GrID{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add a player to the current grid
func (g *GrID) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// Remove a player from the grid
func (g *GrID) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlyerIDs Get all players in the current grid
func (g *GrID) GetPlyerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

// String Print information method
func (g *GrID) String() string {
	return fmt.Sprintf("GrID ID: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
