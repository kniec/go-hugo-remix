package types

import (
	"fmt"
	"io/ioutil"
	"reflect"
)

func compVal(i1, i2 interface{}, fails []string) (res []string) {
	res = fails
	t1 := reflect.TypeOf(i1)
	t2 := reflect.TypeOf(i2)
	if t1 != t2 {
		fmt.Printf("%v :i1 // i2: %v", i1, i2)
		fmt.Printf(" ==> Wrong TYPE\n")
		res = append(fails, fmt.Sprintf("Type: i1:%s != i2:%s", t1, t2))
		return
	}
	if i1 != i2 {
		fmt.Printf("%v :i1 // i2: %v", i1, i2)
		fmt.Printf(" ==> Do not match\n")
		res = append(fails, fmt.Sprintf("Value: i1:%v != i2:%v", i1, i2))
	}
	return
}

func readFile(path string) ([]byte, error) {
	yamlFile, err := ioutil.ReadFile(path)
	return yamlFile, err
}
