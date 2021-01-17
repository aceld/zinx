package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(100, 300, 4, 200, 450, 5)
	fmt.Println(aoiMgr)
}

func TestAOIManagerSuroundGrIDsByGID(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	for k, _ := range aoiMgr.grIDs {
		//得到当前格子周边的九宫格
		grIDs := aoiMgr.GetSurroundGrIDsByGID(k)
		//得到九宫格所有的IDs
		fmt.Println("gID : ", k, " grIDs len = ", len(grIDs))
		gIDs := make([]int, 0, len(grIDs))
		for _, grID := range grIDs {
			gIDs = append(gIDs, grID.GID)
		}
		fmt.Printf("grID ID: %d, surrounding grID IDs are %v\n", k, gIDs)
	}
}
