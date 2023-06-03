package reflect

import (
	"fmt"
	"net/http"
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

//对于匿名字段，可用多级索引（按定义顺序）直接访问
func TestReflectAnonymousField(t *testing.T) {
	type user struct {
		name string
		age  int
	}

	type manager struct {
		user
		title string
	}

	var m manager
	mt := reflect.TypeOf(m)

	// 按名称查找
	name, _ := mt.FieldByName("name")
	t.Log(name.Name, name.Type)

	//按多级索引查找
	age := mt.FieldByIndex([]int{0, 1})
	t.Log(age.Name, age.Type)
}

// 反射能探知当前包或外包的非导出结构成员
func TestReflectNumField(t *testing.T) {
	var s http.Server

	st := reflect.ValueOf(s).Type()

	for i := 0; i < st.NumField(); i++ {
		t.Log(st.Field(i).Name)
	}
}

type X int

func (X) String() string {
	return ""
}

func TestReflectImplements(t *testing.T) {
	var a X
	it := reflect.TypeOf(a)

	//Implements不能直接使用类型作为参数
	st := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	fmt.Println(it.Implements(st))

	// ConvertibleTo() 方法用于检查一个类型是否可以被转换为另一个类型
	ti := reflect.TypeOf(0)
	fmt.Println(ti.ConvertibleTo(it))

	// AssignableTo() 方法用于检查一个类型是否可以被赋值给另一个类型
	fmt.Println(it.AssignableTo(st), it.AssignableTo(ti))
}

// Value 专注于对象实例数据读写
func TestReflectValue(t *testing.T) {
	a := 100
	va, vp := reflect.ValueOf(a), reflect.ValueOf(&a).Elem()

	fmt.Println(va.CanAddr(), va.CanSet())
	fmt.Println(vp.CanAddr(), vp.CanSet())
}

type M struct{}

func (M) Test(x, y int) (int, error) {
	return x + y, nil
}

// 动态调用方法
func TestReflectMethod(t *testing.T) {
	var a M

	v := reflect.ValueOf(&a)
	m := v.MethodByName("Test")

	in := []reflect.Value{
		reflect.ValueOf(1),
		reflect.ValueOf(2),
	}

	out := m.Call(in)

	for _, r := range out {
		t.Log(r)
	}
}

//对于变参来说，用 CallSlice 要更方便一些
type S struct{}

func (S) Format(s string, a ...interface{}) string {
	return fmt.Sprintf(s, a...)
}

func TestReflectCallMethodBySlice(t *testing.T) {
	var a S

	v := reflect.ValueOf(&a)
	m := v.MethodByName("Format")

	// 所有参数都须处理
	out := m.Call([]reflect.Value{
		reflect.ValueOf("%s= %d"),
		reflect.ValueOf("x"),
		reflect.ValueOf(100),
	})
	t.Log(out)

	// 仅一个 []interface{} 即可
	out = m.CallSlice([]reflect.Value{
		reflect.ValueOf("%s= %d"),
		reflect.ValueOf([]interface{}{"x", 100}),
	})
	t.Log(out)
}
