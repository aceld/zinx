# gofmt格式化
# run in terminal:
# make fmt
# win系统中，在git bash中如果出现make包没有找到
# 管理员运行git bash，运行以下命令
# choco install make

GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

fmt:
	gofmt -w $(GOFMT_FILES)
