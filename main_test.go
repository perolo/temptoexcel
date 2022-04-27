package main

import (
	"fmt"
	"testing"
)

func Test_Main_Dummytest(t *testing.T) {

	str := fmt.Sprintf("Dummy string \n")
	if str != "Dummy string \n" {
		t.Error(str)
	}

}
