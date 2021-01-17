package core

import (
	"fmt"
	"sync"
)

/*
	一个地图中的格子类
*/
type GrID struct {
	GID       int          //格子ID
	MinX      int          //格子左边界坐标
	MaxX      int          //格子右边界坐标
	MinY      int          //格子上边界坐标
	MaxY      int          //格子下边界坐标
	playerIDs map[int]bool //当前格子内的玩家或者物体成员ID
	pIDLock   sync.RWMutex //playerIDs的保护map的锁
}

//初始化一个格子
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

//向当前格子中添加一个玩家
func (g *GrID) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

//从格子中删除一个玩家
func (g *GrID) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

//得到当前格子中所有的玩家
func (g *GrID) GetPlyerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

//打印信息方法
func (g *GrID) String() string {
	return fmt.Sprintf("GrID ID: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
