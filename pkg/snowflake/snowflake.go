package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node

func Init(startTime string, machineId int64) (err error) {
	var st time.Time
	if st, err = time.Parse("2006-01-02", startTime); err != nil {
		return err
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineId)
	return err
}

// GenID 生成唯一的uid
func GenID() int64 {
	return node.Generate().Int64()
}
