package utils

import (
	util "github.com/yottalabsai/endorphin/utils"
	"testing"
)

func TestGetSDKVersion(t *testing.T) {
	version := util.GetSDKVersion()
	t.Log(version)
}
