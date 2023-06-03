package function

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

/*golang 函数不⽀持嵌套 (nested)、重载 (overload) 和 默认参数 (default parameter)*/
func funcDefine(x, y int, s string) (int, string) {
	n := x + y
	return n, fmt.Sprintf(s, n)
}

func TestDefineFunc(t *testing.T) {
	_, s := funcDefine(1, 2, "testFunction:%d")
	t.Logf("result:%s", s)
}

/*
*有返回值的函数,必须有明确的 return 终止语句
 */
func returnFunc(x int) int {
	if x > 0 {
		return 1
	} else {
		return -1
	}
}

func TestReturnFunc(t *testing.T) {
	var num = returnFunc(10)
	t.Logf("return num:%d", num)
}

/*命名返回值*/
func paging(sql string, index int) (count, pages int, err error) {
	count = index * 10
	pages = index
	return
}

// 命名返回函数
func TestNameingRetrun(t *testing.T) {
	count, pages, err := paging("select", 10)

	if err != nil {
		t.Fatal("not found this err")
	}

	t.Logf("count %d,pages %d", count, pages)
}

/*函数在 golang 第一类对象 [first-class object] 可在运行期创建，可用作函数参数或返回值，可存入变量的实体。最常见的用法就是匿名函数*/
func testAnonymous(fn func() int) int {
	return fn()
}

/*将复杂签名定义为函数类型*/
type FormatFunc func(s string, x, y int) string

func format(fn FormatFunc, s string, x, y int) string {
	return fn(s, x, y)
}

/*匿名函数做返回值*/
func anonymousReturnFunc() func(int, int) int {
	return func(x, y int) int {
		return x * y
	}
}

/*匿名函数,匿名函数是指没有定义名字符号的函数*/
func TestAnonymousFunc(t *testing.T) {
	s1 := testAnonymous(func() int { return 100 })

	s2 := format(func(s string, x, y int) string {
		return fmt.Sprintf(s, x, y)
	}, "%d,%d", 10, 20)

	t.Log(s1)
	t.Logf(s2)

	/*函数内部定义匿名函数，形成类似嵌套效果。匿名函数可直接调用，作为参数或返回值*/
	func(s string) {
		println(s)
	}("hello,world!")

	/*匿名函数保存到变量*/
	add := func(x, y int) int {
		return x + y
	}
	println(add(10, 20))

	mul := anonymousReturnFunc()

	fmt.Printf("mul:%d \n", mul(10, 20))
}

/*函数可变参数*/
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

/*golang 函数 closure 闭包*/
func closureFunc(x int) func() {
	println(&x)
	return func() {
		println(&x, x)
	}
}

func TestClosureFunc(t *testing.T) {
	f := closureFunc(520)
	f()
}

/*延迟函数*/
func deferFile() error {
	f, err := os.Create("hello.txt")

	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString("hello world!")
	return nil
}

func TestDeferFile(t *testing.T) {
	defer os.Remove("hello.txt")
	deferFile()
}

/*4.6 错误处理, golang 官方推荐的标准做法是返回 error 状态*/
var errDivByZero = errors.New("division by zero")

func div(x, y int) (int, error) {
	if y == 0 {
		return 0, errDivByZero
	} else {
		return x / y, nil
	}
}

func TestErrorFunc(t *testing.T) {
	z, err := div(5, 2)
	if err == errDivByZero {
		t.Error(err)
	}
	println(z)
}

//自定义错误类型
type DivError struct {
	x, y int
}

//实现 error 接口方法
func (DivError) Error() string {
	return "division by zero"
}

//返回自定义错误类型
func divStruct(x, y int) (int, error) {
	if y == 0 {
		return 0, DivError{x, y}
	}

	return x / y, nil
}

func TestDivStruct(t *testing.T) {
	z, err := div(15, 0)

	if err != nil {
		// 根据类型匹配
		switch e := err.(type) {
		case DivError:
			fmt.Println(e, e.x, e.y)
		default:
			fmt.Println(e)
		}

		fmt.Println(err)
	}

	println(z)
}

/*panic，recover 内置函数而非语句。
panic 会立即中断当前函数流程，执行延迟调用。
连续调用 panic，仅最后一个会被 recover 捕获。
而在延迟调用函数中，recover 可捕获并返回 panic 提交的错误对象。
在延迟函数中 panic，不会影响后续延迟调用执行。
而 recover 之后 panic，可被再次捕获。
recover 必须在延迟调用函数中执行才能正常工作。*/
func TestRecover(t *testing.T) {
	defer func() {
		// 捕获错误
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	panic("i am dead") // 引发错误
	println("exit.")   // 永不会执行
}

func panicFunc() {
	defer println("test.1")
	defer println("test.2")

	panic("i am dead")
}

func TestPanicFunc(t *testing.T) {
	defer func() {
		t.Fatal(recover())
	}()

	panicFunc()
}
