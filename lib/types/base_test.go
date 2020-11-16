package types

import (
	"reflect"
	"testing"
)

func Test_CreateBase(t *testing.T) {
	got := CreateBase("Workshop1", ".", "../../misc/test/workshop1", "I. ", 1)
	if !reflect.DeepEqual(got, B1) {
		t.Errorf("CreateBase //\n# Got: %v\n# Exp: %v", got, B1)
	}
}
