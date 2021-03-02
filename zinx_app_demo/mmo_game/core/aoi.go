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

/*
   AOI管理模块
*/
type AOIManager struct {
	MinX  int           //区域左边界坐标
	MaxX  int           //区域右边界坐标
	CntsX int           //x方向格子的数量
	MinY  int           //区域上边界坐标
	MaxY  int           //区域下边界坐标
	CntsY int           //y方向的格子数量
	grIDs map[int]*GrID //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

/*
	初始化一个AOI区域
*/
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

	//给AOI初始化区域中所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子ID
			//格子编号：ID = IDy *nx + IDx  (利用格子坐标得到格子编号)
			gID := y*cntsX + x

			//初始化一个格子放在AOI中的map里，key是当前格子的ID
			aoiMgr.grIDs[gID] = NewGrID(gID,
				aoiMgr.MinX+x*aoiMgr.grIDWIDth(),
				aoiMgr.MinX+(x+1)*aoiMgr.grIDWIDth(),
				aoiMgr.MinY+y*aoiMgr.grIDLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.grIDLength())
		}
	}

	return aoiMgr
}

//得到每个格子在x轴方向的宽度
func (m *AOIManager) grIDWIDth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在x轴方向的长度
func (m *AOIManager) grIDLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n GrIDs in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grID := range m.grIDs {
		s += fmt.Sprintln(grID)
	}

	return s
}

//根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGrIDsByGID(gID int) (grIDs []*GrID) {
	//判断gID是否存在
	if _, ok := m.grIDs[gID]; !ok {
		return
	}

	//将当前gID添加到九宫格中
	grIDs = append(grIDs, m.grIDs[gID])

	// 根据gID, 得到格子所在的坐标
	x, y := gID%m.CntsX, gID/m.CntsX

	// 新建一个临时存储周围格子的数组
	surroundGID := make([]int, 0)

	// 新建8个方向向量: 左上: (-1, -1), 左中: (-1, 0), 左下: (-1,1), 中上: (0,-1), 中下: (0,1), 右上:(1, -1)
	// 右中: (1, 0), 右下: (1, 1), 分别将这8个方向的方向向量按顺序写入x, y的分量数组
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	// 根据8个方向向量, 得到周围点的相对坐标, 挑选出没有越界的坐标, 将坐标转换为gID
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < m.CntsX && newY >= 0 && newY < m.CntsY {
			surroundGID = append(surroundGID, newY*m.CntsX+newX)
		}
	}

	// 根据没有越界的gID, 得到格子信息
	for _, gID := range surroundGID {
		grIDs = append(grIDs, m.grIDs[gID])
	}

	return
}

//通过横纵坐标获取对应的格子ID
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.grIDWIDth()
	gy := (int(y) - m.MinY) / m.grIDLength()

	return gy*m.CntsX + gx
}

//通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGIDByPos(x, y)

	//根据格子ID得到周边九宫格的信息
	grIDs := m.GetSurroundGrIDsByGID(gID)
	for _, v := range grIDs {
		playerIDs = append(playerIDs, v.GetPlyerIDs()...)
		//fmt.Printf("===> grID ID : %d, pIDs : %v  ====", v.GID, v.GetPlyerIDs())
	}

	return
}

//通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPIDsByGID(gID int) (playerIDs []int) {
	playerIDs = m.grIDs[gID].GetPlyerIDs()
	return
}

//移除一个格子中的PlayerID
func (m *AOIManager) RemovePIDFromGrID(pID, gID int) {
	m.grIDs[gID].Remove(pID)
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPIDToGrID(pID, gID int) {
	m.grIDs[gID].Add(pID)
}

//通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGrIDByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grID := m.grIDs[gID]
	grID.Add(pID)
}

//通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGrIDByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grID := m.grIDs[gID]
	grID.Remove(pID)
}
