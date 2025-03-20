package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func TestETHConvertFunc(t *testing.T) {
	amount, err := ETHConvertFunc("20000000000000000")
	if err != nil {
		t.Error()
		return
	}
	t.Log(amount)
}

func TestToDecimal(t *testing.T) {
	v, _ := ToDecimal(strconv.FormatUint(6000000000000000, 10), 18)
	if v > 0.005 {
		fmt.Println("xxxx")
	}
	fmt.Println(v)
}

func TestSprintf(t *testing.T) {
	var userData string = `#!/bin/bash
mkdir ~/.jupyter
nohup jupyter-notebook --ip=0.0.0.0 --port=8888 --no-browser --allow-root --IdentityProvider.token=%s &`
	userData = fmt.Sprintf(userData, "123")
	fmt.Println(userData)
}
