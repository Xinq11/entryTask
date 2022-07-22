package controller

import "testing"

func TestVerifyPassword(t *testing.T) {
	password := "123aaaaS"
	res := VerifyPassword(password)
	if res {
		t.Log("success")
	} else {
		t.Error("fail")
	}
}
