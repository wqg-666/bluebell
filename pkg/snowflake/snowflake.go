package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}

//func main() {
//	if err := Init("2026-11-01", 1); err != nil {
//		fmt.Printf("init failed, err:%v\n", err)
//		return
//	}
//	id := GenID()
//	fmt.Println(id)
//}

//package main
//
//import (
//	"fmt"
//	"time"
//
//	"github.com/sony/sonyflake"
//)
//
//var (
//	sonyFlake     *sonyflake.Sonyflake
//	sonyMachineID uint16
//)
//
//func getMachineID() (uint16, error) {
//	return sonyMachineID, nil
//}
//
//// Init 需传入当前的机器ID
//func Init(startTime string, machineId uint16) (err error) {
//	sonyMachineID = machineId
//	var st time.Time
//	st, err = time.Parse("2006-01-02", startTime)
//	if err != nil {
//		return err
//	}
//	settings := sonyflake.Settings{
//		StartTime: st,
//		MachineID: getMachineID,
//	}
//	sonyFlake = sonyflake.NewSonyflake(settings)
//	return
//}
//
//// GenID 生成id
//func GenID() (id uint64, err error) {
//	if sonyFlake == nil {
//		err = fmt.Errorf("snoy flake not inited")
//		return
//	}
//	id, err = sonyFlake.NextID()
//	return
//}
//
//func main() {
//	if err := Init("2026-10-01", 1); err != nil {
//		fmt.Printf("Init failed, err:%v\n", err)
//		return
//	}
//	id, _ := GenID()
//	fmt.Println(id)
//}
