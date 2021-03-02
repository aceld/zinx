FROM centos:8
COPY zinx_server /zinx-server
COPY /conf/zinx.json /conf/zinx.json
WORKDIR /
EXPOSE  8999

ENTRYPOINT [ "/zinx-server" ]

