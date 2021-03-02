PROJECT_NAME:=zinx_server
VERSION:=v1



.PHONY: image run build clean

build:
	bash build.sh ${PROJECT_NAME}

image:
	docker build -t ${PROJECT_NAME}:${VERSION} .

run:
	docker run  -itd \
	-p 8999:8999 \
	${PROJECT_NAME}:${VERSION}


clean:
	rm -rf ${PROJECT_NAME} 

