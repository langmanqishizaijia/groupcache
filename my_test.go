package groupcache

import (
	"fmt"
	"groupcache/consistenthash"
	"testing"
)
const (
	replicas = 4
)
func TestNewGroup(t *testing.T) {
	mp := consistenthash.New(replicas,nil )
	//判空操作
	fmt.Printf("%v", mp.IsEmpty())

	//Add keys test
	//keys := []string{"zhao","qian","sun","li","zhou","wu","zheng","wang"}
	mp.Add("zhao","qian","sun","li","zhou","wu","zheng")

	for i:=0;i<10;i++ {
		ret := mp.Get("hello")
		fmt.Printf("i=%v, ret=%v\n", i, ret)
	}
}
