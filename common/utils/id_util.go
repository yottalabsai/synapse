package utils

import (
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"synapse/common/log"
	"sync"
)

var (
	node     *snowflake.Node
	nodeOnce sync.Once
)

func GetSnowflakeNode() *snowflake.Node {
	nodeOnce.Do(func() {
		var err error
		node, err = snowflake.NewNode(1)
		if err != nil {
			log.Log.Fatalw("Failed to create snowflake node", zap.Error(err))
		}
	})
	return node
}

func GetNextId() int64 {
	return GetSnowflakeNode().Generate().Int64()
}
