package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_version_ex/protoDemo/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name:   "XiaoYuer",
		Age:    16,
		Emails: []string{"xiao_yu_er@sina.com", "yu_er@sina.cn"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "13113111311",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "14141444144",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "19191919191",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	newdata := &pb.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		fmt.Println("unmarshal err:", err)
	}
	fmt.Println(newdata)
}
