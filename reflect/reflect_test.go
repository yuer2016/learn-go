package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

/**
* 反射（reflect）机制能在运行期探知对象的类型信息和内存结构
 */
func TestReflectFeild(t *testing.T) {
	type X int
	var a X = 100
	t1 := reflect.TypeOf(a)
	fmt.Println(t1.Name(), t1.Kind())

	//除了通过实际对象获取类型外，还可直接构造一些基础复合类型
	a1 := reflect.ArrayOf(10, reflect.TypeOf(byte(0)))
	m := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(0))
	fmt.Println(a1, m)

	//方法 Elem 返回指针、数组、切片、字典（值）或通道的基类型
	fmt.Println(reflect.TypeOf(map[string]int{}).Elem())
	fmt.Println(reflect.TypeOf([]int32{}).Elem())
}

func TestReflectElemValue(t *testing.T) {
	a := 100
	va, vp := reflect.ValueOf(a), reflect.ValueOf(&a).Elem()

	fmt.Println(va.CanAddr(), va.CanSet())
	fmt.Println(vp.CanAddr(), vp.CanSet())
}

func TestReflectStruct(t *testing.T) {
	type user struct {
		name string
		age  int
	}

	type manager struct {
		user
		title string
	}

	var m manager
	r := reflect.TypeOf(&m)

	// 获取指针的基类型
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		fmt.Println(f.Name, f.Type, f.Offset)

		// 输出匿名字段结构
		if f.Anonymous {
			for x := 0; x < f.Type.NumField(); x++ {
				af := f.Type.Field(x)
				fmt.Println("  ", af.Name, af.Type)
			}
		}
	}
}
