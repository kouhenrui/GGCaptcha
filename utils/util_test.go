package utils

import "testing"

func Test_RandMath2(t *testing.T) {

	ex, an := RandMath2()
	t.Logf(ex)
	t.Log(an)

	//t.Log(10 % 5)
}

func Test_RandMath(t *testing.T) {

	ex, an := RandMath(3)
	t.Log(ex)
	t.Log(an)
}
