package Notify

import (
	"github.com/aceld/zinx/ziface"
)

//建立一个用户自定义ID和连接映射的结构
//map会存在 并发问题，大量数据循环读取问题
//暂时先用map结构存储，但是应该不是最好的选择，抛砖引玉

type ConnIDMap map[uint64]ziface.IConnection
