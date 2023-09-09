package convertutil

import (
	"fmt"
	"testing"
	"time"
)

type PetPo struct {
	Id        int64
	PetName   string
	CreatedAt time.Time
}

type PetVo struct {
	Id        int64
	PetName   string
	CreatedAt string
}

type UserPo struct {
	Id        int64
	Username  string
	Age       int64
	Pet       PetPo
	CreatedAt time.Time
}

type UserVo struct {
	Id       int64
	Username string
	Age      int64
	Pet      PetVo
}

func TestConvertStruct(t *testing.T) {
	userPo := &UserPo{
		Id:        1,
		Username:  "cotton",
		Age:       11,
		CreatedAt: time.Now(),
		Pet: PetPo{
			Id:      2,
			PetName: "sss",
		},
	}
	fmt.Printf("source: %+v\n", userPo)
	target := ConvertStruct[UserVo](*userPo)
	fmt.Printf("============\n")
	fmt.Printf("target: %+v\n", target)
}
