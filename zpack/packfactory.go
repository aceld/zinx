package zpack

import (
	"github.com/aceld/zinx/ziface"
	"sync"
)

var pack_once sync.Once

type pack_factory struct{}

var factoryInstance *pack_factory

/*
	Generates different packaging and unpackaging methods, singleton
	(生成不同封包解包的方式，单例)
*/
func Factory() *pack_factory {
	pack_once.Do(func() {
		factoryInstance = new(pack_factory)
	})

	return factoryInstance
}

// NewPack creates a concrete packaging and unpackaging object
// (NewPack 创建一个具体的拆包解包对象)
func (f *pack_factory) NewPack(kind string) ziface.IDataPack {
	var dataPack ziface.IDataPack

	switch kind {
	// Zinx standard default packaging and unpackaging method
	// (Zinx 标准默认封包拆包方式)
	case ziface.ZinxDataPack:
		dataPack = NewDataPack()
		break
	case ziface.ZinxDataPackOld:
		dataPack = NewDataPackLtv()
		break

		// case for custom packaging and unpackaging methods
		// (case 自定义封包拆包方式case)

	default:
		dataPack = NewDataPack()
	}

	return dataPack
}
