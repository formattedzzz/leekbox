package structtomap

import (
	"reflect"
	"testing"
)

// struct递归转换为map[string]interface{}
func TestStructToMap(t *testing.T) {
	type demo struct {
		input interface{}
		want  SupMap
	}

	type T1 struct {
		Name  string `json:"name"`
		Age   int
		Marry bool
	}
	type T2_1 struct {
		Info string `json:"info"`
	}
	type T2 struct {
		Array []interface{} `json:"array"`
		Other T2_1
	}

	demo1 := demo{
		input: T1{
			Name: "leooo",
		},
		want: SupMap{
			"name":  "leooo",
			"Age":   0,
			"Marry": false,
		},
	}
	demo2 := demo{
		input: T2{
			Array: []interface{}{1, "", true, 1.2},
			Other: T2_1{
				Info: "shiyu",
			},
		},
		want: SupMap{
			"array": []interface{}{1, "", true, 1.2},
			"Other": SupMap{
				"info": "shiyu",
			},
		},
	}
	demo3 := demo{
		input: &T1{
			Name: "leooo",
		},
		want: SupMap{
			"name":  "leooo",
			"Age":   0,
			"Marry": false,
		},
	}
	demo4 := demo{
		input: []interface{}{1, true, ""},
		want:  nil,
	}
	demos := []demo{demo1, demo2, demo3, demo4}
	for _, demo := range demos {
		got, _ := StructToMap(demo.input, "json")
		if !reflect.DeepEqual(got, demo.want) {
			t.Errorf("excepted:%#v, got:%#v", demo.want, got)
		}
	}
}
