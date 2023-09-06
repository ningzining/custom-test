package convertutil

import (
	"fmt"
	"testing"
	"time"
)

func TestConvertStruct(t *testing.T) {
	userPo := &UserPo{
		Id:        1,
		Username:  "cotton",
		Age:       11,
		CreatedAt: time.Now(),
	}
	fmt.Printf("source: %+v\n", userPo)
	target := ConvertStruct[*UserPo, UserVo](userPo)
	fmt.Printf("============\n")
	fmt.Printf("target: %+v\n", target)
}
