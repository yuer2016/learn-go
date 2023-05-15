package method

import (
	"fmt"
	"testing"
	"unsafe"
)

/**
* 接口代表一种调用契约，是多个方法声明的集合。
* Go 接口实现机制很简洁，只要目标类型方法集内包含接口声明的全部方法，就被视为实现了该接口，无须做显示声明。目标类型可实现多个接口。
*
* 接口自身也是一种结构类型，只是编译器会对其做出很多限制:
* 1. 不能有字段
* 2. 不能定义自己的方法。
* 3. 只能声明方法，不能实现。
* 4. 可嵌入其他接口类型。
 */
type tester interface {
	test()
	string() string
}

type data struct{}

func (*data) test()         {}
func (data) string() string { return "interface impl" }

func TestDataInterface(t *testing.T) {
	var d data
	var tdt tester = &d
	tdt.test()
	println(tdt.string())
}

/*
* 如果接口没有任何方法声明，那么就是一个空接口（interface{}），它的用途类似面向对象里的根类型 Object，可被赋值为任何类型的对象。
* 接口变量默认值是 nil。如果实现接口的类型支持，可做相等运算。
 */
func TestEmptyInterface(t *testing.T) {
	var t1, t2 interface{}
	println(t1 == nil, t1 == t2)

	t1, t2 = 100, 100
	println(t1 == t2)

	//t1, t2 = map[string]int{}, map[string]int{}
	//println(t1 == t2)
}

/**
* 嵌入其他接口类型，相当于将其声明的方法集导入。这就要求不能有同名方法，因为方法不支持重载。不能嵌入自身或循环嵌入，那会导致递归错误。
 */
type stringer interface {
	str() string
}

type testinger interface {
	stringer
	test()
}

type datas struct{}

func (*datas) test() {}
func (datas) str() string {
	return "embedded interface"
}

func TestEmbeddedInterface(t *testing.T) {
	var data datas
	var tester = &data

	tester.test()
	t.Log(tester.str())
}

func pp(a stringer) {
	println(a.str())
}

/**
* 超集接口变量可隐式转换为子集，反过来不行
 */
func TestSuperInterface(t *testing.T) {
	var d datas
	var tt testinger = &d

	pp(tt)
}

type iface struct {
	tab  *itab          // 类型信息
	data unsafe.Pointer // 实际对象指针
}

type itab struct {
	inter *interfacetype // 接口类型
	_type *_type         // 实际对象类型
	fun   [1]uintptr     // 实际对象方法地址
}

type interfacetype interface{}
type _type interface{}

/**
* 类型转换
* 类型推断可将接口变量还原为原始类型,或用来判断是否实现了某个更具体的接口类型
 */
type dataInt int

func (d dataInt) String() string {
	return fmt.Sprintf("data:%d", d)
}

func TestTypingTransform(t *testing.T) {
	var d dataInt = 15
	var x interface{} = d

	// 转换为更具体的接口类型
	if n, ok := x.(fmt.Stringer); ok {
		fmt.Println(n)
	}

	// 转换回原始类型
	if d2, ok := x.(dataInt); ok {
		fmt.Println(d2)
	}
}
