package method

import (
	"fmt"
	"testing"
)

/*方法是与对象实例绑定的特殊函数。方法是有关联状态的，而函数通常没有。*/
//可以为当前包,以及除接口和指针以外的任何类型定义方法,方法同样不支持重载(overload)
type N int

func (n N) toString() string {
	return fmt.Sprintf("%#x", n)
}

func (n N) value() {
	n++
	fmt.Printf("v: %p, %v\n", &n, n)
}

func (n *N) pointer() {
	(*n)++
	fmt.Printf("p: %p, %v\n", n, *n)
}

func TestMethodDefine(t *testing.T) {
	var a N = 25
	p := &a
	println(a.toString())

	println()

	a.value()
	a.pointer()

	println()

	p.value()
	p.pointer()

	println()
	fmt.Printf("a: %p, %v\n", &a, a)
}

/*
如何选择方法的 receiver 类型？
要修改实例状态,用 *T
无须修改状态的小对象或固定值，建议用 T
大对象建议用 *T, 以减少复制成本
引用类型、字符串、函数等指针包装对象，直接用 T
若包含 Mutex 等同步字段, 用 *T，避免因复制造成锁操作无效
其他无法确定的情况, 都用 *T*/
