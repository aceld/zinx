package core

import "fmt"

const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

// AOIManager AOI management module(AOI管理模块)
type AOIManager struct {
	MinX  int           // Left boundary coordinate of the area(区域左边界坐标)
	MaxX  int           // Right boundary coordinate of the area(区域右边界坐标)
	CntsX int           // Number of grids in the x direction(x方向格子的数量)
	MinY  int           // Upper boundary coordinate of the area(区域上边界坐标)
	MaxY  int           // Lower boundary coordinate of the area(区域下边界坐标)
	CntsY int           // Number of grids in the y direction(y方向的格子数量)
	grIDs map[int]*GrID // Which grids are present in the current area, key = grid ID, value = grid object(当前区域中都有哪些格子，key=格子ID， value=格子对象)
}

// NewAOIManager Initialize an AOI area(初始化一个AOI区域)
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grIDs: make(map[int]*GrID),
	}

	// Initialize all grids in the AOI region (给AOI初始化区域中所有的格子)

	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// Calculate the grid ID
			// Grid number: ID = IDy *nx + IDx (obtain grid number using grid coordinates)
			// 计算格子ID
			// 格子编号：ID = IDy *nx + IDx  (利用格子坐标得到格子编号)
			gID := y*cntsX + x

			// Initialize a grid in the AOI map, where the key is the current grid's ID
			// 初始化一个格子放在AOI中的map里，key是当前格子的ID
			aoiMgr.grIDs[gID] = NewGrID(gID,
				aoiMgr.MinX+x*aoiMgr.grIDWIDth(),
				aoiMgr.MinX+(x+1)*aoiMgr.grIDWIDth(),
				aoiMgr.MinY+y*aoiMgr.grIDLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.grIDLength())
		}
	}

	return aoiMgr
}

// grIDWIDth Get the width of each grid in the x-axis direction
// (得到每个格子在x轴方向的宽度)
func (m *AOIManager) grIDWIDth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// grIDLength Get the length of each grid in the x-axis direction
// (得到每个格子在x轴方向的长度)
func (m *AOIManager) grIDLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// String Print information method
// (打印信息方法)
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n GrIDs in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grID := range m.grIDs {
		s += fmt.Sprintln(grID)
	}

	return s
}

// GetSurroundGrIDsByGID Get the surrounding nine grids information based on the grid's gID
// 根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGrIDsByGID(gID int) (grIDs []*GrID) {
	// Check if gID exists
	// 判断gID是否存在
	if _, ok := m.grIDs[gID]; !ok {
		return
	}

	// Add the current gID to the nine grids
	// 将当前gID添加到九宫格中
	grIDs = append(grIDs, m.grIDs[gID])

	// Get the coordinates of the grid based on gID
	// 根据gID, 得到格子所在的坐标
	x, y := gID%m.CntsX, gID/m.CntsX

	// Create a temporary array to store the coordinates of the surrounding grids
	// 新建一个临时存储周围格子的数组
	surroundGID := make([]int, 0)

	// Create eight direction vectors: Upper left: (-1, -1), Left middle: (-1, 0), Upper right: (-1,1),
	// Middle upper: (0,-1), Middle lower: (0,1), Right upper: (1, -1), Right middle: (1, 0), Right lower: (1, 1),
	// respectively insert these eight direction vectors into the x and y component arrays in order
	// 新建8个方向向量: 左上: (-1, -1), 左中: (-1, 0), 左下: (-1,1), 中上: (0,-1), 中下: (0,1), 右上:(1, -1)
	// 右中: (1, 0), 右下: (1, 1), 分别将这8个方向的方向向量按顺序写入x, y的分量数组
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	// Get the relative coordinates of the surrounding points based on the eight direction vectors,
	// select the coordinates that do not go out of bounds, and convert the coordinates to gID
	// 根据8个方向向量, 得到周围点的相对坐标, 挑选出没有越界的坐标, 将坐标转换为gID
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < m.CntsX && newY >= 0 && newY < m.CntsY {
			surroundGID = append(surroundGID, newY*m.CntsX+newX)
		}
	}

	// Get grid information based on valid gID
	// 根据没有越界的gID, 得到格子信息
	for _, gID := range surroundGID {
		grIDs = append(grIDs, m.grIDs[gID])
	}

	return
}

// GetGIDByPos Get the corresponding grid ID by horizontal and vertical coordinates
// 通过横纵坐标获取对应的格子ID
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.grIDWIDth()
	gy := (int(y) - m.MinY) / m.grIDLength()

	return gy*m.CntsX + gx
}

// GetPIDsByPos Get all PlayerIDs within the surrounding nine grids by horizontal and vertical coordinates
// 通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	// Get which grid ID the current coordinates belong to
	// 根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGIDByPos(x, y)

	// Get information about the surrounding nine grids based on the grid ID
	// 根据格子ID得到周边九宫格的信息
	grIDs := m.GetSurroundGrIDsByGID(gID)
	for _, v := range grIDs {
		playerIDs = append(playerIDs, v.GetPlyerIDs()...)
		//fmt.Printf("===> grID ID : %d, pIDs : %v  ====", v.GID, v.GetPlyerIDs())
	}

	return
}

// GetPIDsByGID Get all PlayerIDs within the surrounding nine grids by horizontal and vertical coordinates
// 通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPIDsByGID(gID int) (playerIDs []int) {
	playerIDs = m.grIDs[gID].GetPlyerIDs()
	return
}

// RemovePIDFromGrID Remove a PlayerID from a grid
// 移除一个格子中的PlayerID
func (m *AOIManager) RemovePIDFromGrID(pID, gID int) {
	m.grIDs[gID].Remove(pID)
}

// AddPIDToGrID Add a PlayerID to a grid
// 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPIDToGrID(pID, gID int) {
	m.grIDs[gID].Add(pID)
}

// AddToGrIDByPos Add a Player to a grid based on horizontal and vertical coordinates
// 通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGrIDByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grID := m.grIDs[gID]
	grID.Add(pID)
}

// RemoveFromGrIDByPos Remove a Player from the corresponding grid based on horizontal and vertical coordinates
// 通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGrIDByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grID := m.grIDs[gID]
	grID.Remove(pID)
}
