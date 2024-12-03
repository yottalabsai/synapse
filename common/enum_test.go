package common

import "testing"

func Test_UnmarshalJSON(t *testing.T) {
	s := "GPU_1"
	switch s {
	case string(GPU_1), string(GPU_2), string(GPU_4), string(GPU_8), string(GPU_16), string(GPU_32):
		t.Log("success")
	}
}
