package core

import (
	"sync"
)

// WorldManager The overall management module of the current game world 当前游戏世界的总管理模块
type WorldManager struct {
	AoiMgr  *AOIManager       // AOI planning manager for the current world map(当前世界地图的AOI规划管理器)
	Players map[int32]*Player // Collection of currently online players(当前在线的玩家集合)
	pLock   sync.RWMutex      // Mutual exclusion mechanism to protect Players(保护Players的互斥读写机制)
}

// Provide an external handle to the world management module
// 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldManager

// Provide an initialization method for WorldManager
// 提供WorldManager 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
	}
}

// AddPlayer Provide the ability to add a player, adding the player to the player information table Players
// (提供添加一个玩家的的功能，将玩家添加进玩家信息表Players)
func (wm *WorldManager) AddPlayer(player *Player) {
	// Add the player to the world manager
	// 将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.Players[player.PID] = player
	wm.pLock.Unlock()

	// Add the player to the AOI network planning
	// 将player 添加到AOI网络规划中
	wm.AoiMgr.AddToGrIDByPos(int(player.PID), player.X, player.Z)
}

// RemovePlayerByPID Remove a player from the player information table by player ID
// 从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayerByPID(pID int32) {
	wm.pLock.Lock()
	delete(wm.Players, pID)
	wm.pLock.Unlock()
}

// GetPlayerByPID Get corresponding player information by player ID
// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPID(pID int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pID]
}

// GetAllPlayers Get information of all players
// 获取所有玩家的信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	// Create a slice to return player collection
	// 创建返回的player集合切片
	players := make([]*Player, 0)

	// Add to the slice
	// 添加切片
	for _, v := range wm.Players {
		players = append(players, v)
	}

	return players
}

// GetPlayersByGID Get information of all players in a specific gID
// 获取指定gID中的所有player信息
func (wm *WorldManager) GetPlayersByGID(gID int) []*Player {
	// Get all pIDs corresponding to the gID
	// 通过gID获取 对应 格子中的所有pID
	pIDs := wm.AoiMgr.grIDs[gID].GetPlyerIDs()

	// Get player objects corresponding to pIDs
	// 通过pID找到对应的player对象
	players := make([]*Player, 0, len(pIDs))
	wm.pLock.RLock()
	for _, pID := range pIDs {
		players = append(players, wm.Players[int32(pID)])
	}
	wm.pLock.RUnlock()

	return players
}
