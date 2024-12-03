package utils

import (
	"synapse/config"
)

func InitRedis() {
	config.InitConfig("../resources", "")
	config.InitRedis()
}
