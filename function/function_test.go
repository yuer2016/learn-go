package function

import (
	"fmt"
	"testing"
)

/*golang 函数不⽀持嵌套 (nested)、重载 (overload) 和 默认参数 (default parameter)*/
func test(x, y int, s string) (int, string) {
	n := x + y
	return n, fmt.Sprintf(s, n)
}

func TestDefineFunc(t *testing.T) {
	_, s := test(1, 2, "testFunction:%d")
	t.Logf("result:%s", s)
}

/*函数是第⼀类对象，可作为参数传递。建议将复杂签名定义为函数类型*/
func testAnonymous(fn func() int) int {
	return fn()
}

type FormatFunc func(s string, x, y int) string

func format(fn FormatFunc, s string, x, y int) string {
	return fn(s, x, y)
}

func TestAnonymousFunc(t *testing.T) {
	s1 := testAnonymous(func() int { return 100 })

	s2 := format(func(s string, x, y int) string {
		return fmt.Sprintf(s, x, y)
	}, "%d,%d", 10, 20)

	t.Log(s1)
	t.Logf(s2)
}

func variants(s string, n ...int) string {
	var x int
	for _, i := range n {
		x += i
	}

	return fmt.Sprintf(s, x)
}

func TestVariants(t *testing.T) {
	s := []int{1, 2, 3}
	t.Log(variants("sum1:%d", 1, 2, 3))
	t.Log(variants("sum2:%d", s...))
}
