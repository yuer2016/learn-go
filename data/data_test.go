package data

import (
	"fmt"
	"testing"
	"unsafe"
)

/* 字符串是不可变字节(byte)序列，其本身是一个复合结构。*/
type stringStruct struct {
	str unsafe.Pointer
	len int
}

/*字符串是不可变字节(byte)序列,其本身是一个复合结构。*/
func TestString(t *testing.T) {
	s := "鱼儿\x61\142\u0041"

	fmt.Printf("%s\n", s)
	fmt.Printf("%x,len: %d\n", s, len(s))

	//字符串默认值不是 nil，而是 ""
	var s1 string
	println(s1 == "")

	// 切片语法
	s2 := "abcdefg"

	s3 := s2[:3]  // 从头开始,仅指定结束索引位置
	s4 := s2[1:4] // 指定开始和结束位置，返回 [start,end)
	s5 := s2[2:]  // 指定开始位置，返回后面全部内容

	println(s3, s4, s5)
}

/*定义数组类型时,数组长度必须是非负整型常量表达式,长度是类型组成部分。
元素类型相同,但长度不同的数组不属于同一类型*/
func TestArray(t *testing.T) {
	var a [4]int // 元素自动初始化为零

	b := [4]int{2, 5}     //未提供初始值的元素自动初始化为 0
	c := [4]int{5, 3: 10} //可指定索引位置初始化

	d := [...]int{1, 2, 3}    //编译器按初始化值数量确定数组长度
	e := [...]int{10, 5: 100} //支持索引初始化，但注意数组长度与此有关

	fmt.Println(a, b, c, d, e)
}

/*对于结构等复合类型数组,可省略元素初始化类型标签*/
func TestStructArray(t *testing.T) {
	type user struct {
		name string
		age  byte
	}

	//省略了类型标签
	d := [...]user{
		{"Tom", 20},
		{"Mary", 18},
	}

	fmt.Printf("%#v\n", d)
}

/*在定义多维数组时，仅第一维度允许使用 “...”*/
func TestMulitArray(t *testing.T) {
	a := [2][2]int{
		{1, 2},
		{3, 4},
	}

	b := [...][2]int{
		{10, 20},
		{30, 40},
	}
	// 三维数组
	c := [...][2][2]int{
		{
			{1, 2},
			{3, 4},
		},
		{
			{10, 20},
			{30, 40},
		},
	}

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

	/*内置函数 len 和 cap 都返回第一维度长度*/
	println(len(a), cap(a))
	println(len(b), cap(b))
	println(len(b[1]), cap(b[1]))
}

/*与 C 数组变量隐式作为指针使用不同,Go 数组是值类型,赋值和传参操作都会复制整个数组数据*/
func TestCopyArray(t *testing.T) {
	a := [2]int{10, 20}
	var b [2]int
	b = a

	fmt.Printf("a: %p, %v\n", &a, a)
	fmt.Printf("b: %p, %v\n", &b, b)

	func(x [2]int) {
		fmt.Printf("x: %p, %v\n", &x, x)
	}(a)

	//可改用指针或切片，以此避免数据复制
	func(x *[2]int) {
		fmt.Printf("&x: %p, %v\n", x, *x)
	}(&a)
}

/*
切片(slice) 本身并非动态数组或数组指针。
它内部通过指针引用底层数组，设定相关属性将数据读写操作限定在指定区域内.
切片本身是个只读对象，其工作机制类似数组指针的一种包装
*/
type slice struct {
	array unsafe.Pointer
	len   int // len 用于限定可读的写元素数量
	cap   int // 属性 cap 表示切片所引用数组片段的真实长度
}

func TestSlice(t *testing.T) {
	x := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := x[2:5]

	println(cap(s), len(s))

	for i := 0; i < len(s); i++ {
		println(s[i])
	}
}

/*可直接创建切片对象,无须预先准备数组。
因为是引用类型,须使用 make 函数或显式初始化语句,它会自动完成底层数组内存分配*/
func TestInitSlice(t *testing.T) {
	s1 := make([]int, 3, 5)    // 指定 len、cap,底层数组初始化为零值
	s2 := make([]int, 3)       // 省略 cap, 和 len 相等
	s3 := []int{10, 20, 5: 30} // 按初始化元素分配底层数组,并设置 len、cap

	fmt.Println(s1, len(s1), cap(s1))
	fmt.Println(s2, len(s2), cap(s2))
	fmt.Println(s3, len(s3), cap(s3))
}

func TestCopySlice(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s1 := s[5:8]
	n := copy(s[4:], s1) //在同一底层数组的不同区间复制
	fmt.Println(n, s)

	s2 := make([]int, 6) //在不同数组间复制
	n = copy(s2, s)
	fmt.Println(n, s2)
}

/*
字典(哈希表)是一种使用频率极高的数据结构。
将其作为语言内置类型，从运行时层面进行优化,可获得更高效的性能。
作为无序键值对集合,字典要求 key 必须是支持相等运算符（==、!=）的数据类型,比如,数字、字符串、指针、数组、结构体,以及对应接口类型。
字典是引用类型，使用 make 函数或初始化表达语句来创建
*/
func TestMap(t *testing.T) {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2

	//值为匿名结构类型
	m2 := map[int]struct {
		x int
	}{
		1: {x: 100}, //可省略 key、value 类型标签
		2: {x: 200},
	}

	fmt.Println(m, m2)

	m["a"] = 10 //修改
	m["c"] = 30 //新增

	/*使用 ok-idiom 判断 key 是否存在,返回值.访问不存在的键值，默认返回零值，不会引发错误。
	但推荐使用 ok-idiom 模式，毕竟通过零值无法判断键值是否存在，或许存储的 value 本就是零*/
	if v, ok := m["d"]; ok {
		println(v)
	}

	//删除键值对。不存在时，不会报错
	delete(m, "d")
}

/*字典进行迭代，每次返回的键值次序都不相同*/
func TestMapRange(t *testing.T) {
	m := make(map[string]int)

	for i := 0; i < 8; i++ {
		m[string('a'+rune(i))] = i
	}

	for i := 0; i < 4; i++ {
		for k, v := range m {
			print(k, ":", v, "  ")
		}
		println()
	}
}
