package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int
	MinX      int
	MinY      int
	MaxX      int
	MaxY      int
	playerIDS map[int]bool //当前格子内的玩家或者物体成员ID
	pidLock   sync.RWMutex //playerIDs的保护map的锁
}

func NewGrid(gID, minX, minY, maxX, maxY int) *Grid {
	grid := &Grid{
		GID:       gID,
		MinX:      minX,
		MinY:      minY,
		MaxX:      maxX,
		MaxY:      maxY,
		playerIDS: make(map[int]bool),
	}

	return grid
}

// Add 向当前格子中添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pidLock.Lock()
	defer g.pidLock.Unlock()

	g.playerIDS[playerID] = true
}

// Remove 从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pidLock.Lock()
	defer g.pidLock.Unlock()

	delete(g.playerIDS, playerID)
}

// GetPlyerIDs 得到当前格子中所有的玩家
func (g *Grid) GetPlyerIDs() (playerIDs []int) {
	g.pidLock.RLock()
	defer g.pidLock.RUnlock()

	for k := range g.playerIDS {
		playerIDs = append(playerIDs, k)
	}

	return
}

//打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("GrID ID: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDS)
}
