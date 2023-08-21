# Learn-Go

golang programming lanuage learning project

I think there are two difficulties in go language, one is network and the other is concurrency, so this project contains examples of two topics.

## Golang

Go 是一门编译型语言，Go 语言的工具链将源代码及其依赖转换成计算机的机器指令。

### Go 配置

```bash
export GOPATH="${HOME}/.go"
export GOROOT="/usr/local/opt/go/libexec"
export GOBIN=$GOPATH/bin
export GOPKG=$GOPATH/pkg
export GOPROXY=https://mirrors.tencent.com/go/
export GOPRIVATE=*.code.oa.com,git.woa.com
export GO111MODULE=on
export G_MIRROR=https://golang.google.cn/dl/
unset GODEBUG
```

```go
go env -w GOPROXY=https://mirrors.tencent.com/go/
```

## Go 依赖管理

### GOROOT

指的是 Go 的 **编译器** 和 **标准库** ,二者位于同一个安装包中。属于 Go 语言的顶级目录。

### GOPATH

实际中的 Go 语言项目是由一个或多个 package 组成的，这些 package 按照来源分为以下几种:

* 标准库
* 第三方库
* 项目私有库

标准库的 package 全部位于 GOROOT 环境变量指示的目录中，第三方库和项目私有库位于 GOPATH 环境变量所指示的目录中。

GOPATH 是用户工作空间的目录的环境变量，属于用户域范畴。

当某个 package 需要引用其他包时，编译器就会依次从 GOROOT/src 和 GOPATH/src 依次查找；从 GOROOT 下找到,就不会再到 GOPATH 目录下查找。

GOPATH 的问题在于，在实际的工程项目中，如果两个项目引用不同版本的 第三方库 则这两个项目无法共享同一个 GOPATH 。

### GoModule

为了解决 GOPATH 不同库版本依赖不能共享的问题。GO 1.11 中首次引入了 Module 特性。

GO Module 核心解决两个重要的问题:

* 准确的记录项目依赖
* 可重复的构建

A module is a collection of related Go package that are versioned together as a single unit。

### 工作空间

编译工具对源码目录有严格要求，每个工作空间 (workspace) 必须由 **bin、pkg、src** 三个目录组成。

``` go
workspace 
    | +--- bin // go install 安装目录 
             +--- learn  
    | +--- pkg // go build ⽣成静态库 (.a) 存放目录
              +--- darwin_amd64  
              +--- mylib.a  
              +--- mylib  
                      +--- sublib.a  
    | +--- src // 项目源码目录 
             +--- learn  
                      +--- main.go  
             +--- mylib  
                      +--- mylib.go  
                      +--- sublib  
                               +--- sublib.go
```

可在 **GOPATH** 环境变量列表中添加多个工作空间,但不能和 **GOROOT** 相同。

```bash
export GOPATH=$HOME/projects/golib:$HOME/projects/go

go mod init MODULE_NAME
```

通常 **go get** 使用第一个工作空间保存下载的第三方库。

### 源文件

* **编码**：源码文件必须是 UTF-8 格式，否则会导致编译器出错。
* **结束**：语句以 ";" 结束，多数时候可以省略。
* **注释**：⽀持 "//"、"/\*\*/" 两种注释方式，不能嵌套。
* **命名**：采用 camelCasing 风格，不建议使用下划线。

### 包结构

所有代码都必须组织在 package 中。

* **源文件** 头部以  package 声明包名称。
* 包由同一目录下的多个源码文件组成。
* 包名类似 **namespace**，与包所在目录名、编译文件名无关。
* 目录名最好不用 **main、all、std** 这三个保留名称。
* 可执行文件必须包含 **package main**，入口函数 **main**。

说明：**os.Args** 返回命令行参数，**os.Exit** 终止进程。
要获取正确的可执⾏文件路径，可用 **filepath.Abs(exec.LookPath(os.Args[0]))** 。

包中成员以名称首字母大小写决定访问权限。

* public: 首字母大写，可被包外访问。
* internal: 首字母小写，仅包内成员可以访问。

该规则适用于 **全局变量、全局常量、类型、结构字段、函数、方法** 等。

### 导入包

使用包成员前，必须先用 **import** 关键字导⼊，但不能形成导入循环。
**import "相对目录/包主文件名"**
相对目录是指从 **/pkg/** 开始的子目录，以标准库为例：

```go
import "fmt" -> /usr/local/go/pkg/darwin_amd64/fmt.a 
import "os/exec" -> /usr/local/go/pkg/darwin_amd64/os/exec.a
```

在导入时，可指定包成员访问方式。
比如对包重命名，以避免同名冲突。

```go
import "github/test" // 默认模式: test.A
import M "github/test" // 包重命名: M.A 
import . "github/test" // 简便模式: A 
import _ "github/test" // 非导⼊模式: 仅让该包执⾏初始化函数
```

未使用的导入包，会被编译器视为错误 (不包括 "import \_")。

***./main.go:4: imported and not used: "fmt"***

对于当前目录下的子包，除使用默认完整导入路径外，还可使用 local 方式。

```go
workspace  
    +--- src  
        +--- learn  
            +--- main.go  
            +--- test  
                +--- test.go 
                
main.go 
import "learn/test" // 正常模式 
import "./test" // 本地模式，仅对 go run main.go 有效 
```

## Golang 基础过程抽象

Go 语言和其他编程语言一样，一个大的 **程序** 是由很多小的基础构件组成的：

* 变量保存值 ;
* 简单的 **加法和减法** 运算被组合成较复杂的表达式 ;
* 基础类型被聚合为 **数组** 或 **结构体** 等更复杂数据结构 ;
* 使用 if 和 for 之类的控制语句来组织和控制表达式执行流程 ;
* 多个语句被组织到一个个函数中，以便代码的 **隔离** 和 **复用** ;
* 函数以 **源文件** 和 **包** 的方式被组织;

## Golang 全局函数

| 内建函数名           | 作用                                                                               |
| :------------------- | :--------------------------------------------------------------------------------- |
| append()             | 用于向切片追加元素，并返回一个新的切片（注意，原始切片的底层数组可能会被扩展）     |
| cap()                | 用于获取切片、数组或通道的容量。容量是指底层数组可以存储的元素数目                 |
| close()              | 用于关闭通道。关闭通道后，无法再向其发送值，但仍然可以从中读取已有的值             |
| complex()            | 用于创建一个复数值                                                                 |
| copy()               | 用于将一个切片的内容复制到另一个切片中。如果两个切片长度不一致，则会复制较短的那个 |
| delete()             | 用于删除一个 map 中的键值对                                                        |
| imag()               | 用于获取一个复数的虚部                                                             |
| len()                | 用于获取切片、数组、字符串、通道或映射的长度                                       |
| make()               | 用于创建切片、映射或通道，并返回它们的引用                                         |
| new()                | 用于创建变量的指针，并返回指向该变量的指针                                         |
| panic()              | 用于引发一个运行时错误，并打印错误信息                                             |
| print() 和 println() | 用于在标准输出中打印文本或变量的值                                                 |
| real()               | 用于获取一个复数的实部                                                             |
| recover()            | 用于从运行时错误中恢复，并返回该错误对象                                           |

## Golang 数据抽象

### 变量

Go 是 **静态类型** 语言，**不能在运行期改变变量类型**。
使用关键字 ***var*** **定义变量**，自动初始化为 **零值**。
如果 **提供初始化值**，可 **省略** 变量类型，由编译器 **自动推断** 。

```go
/*golang 定义变量的三种方式*/
var x int
var f float32 = 1.6
var s = "abc"

/*在函数内部，可用更简略的 ":=" 方式定义变量*/
func main() {
    //注意检查，是定义新局部变量，还是修改全局变量
    x := 123 
}

/*可一次定义多个变量*/
var x, y, z int 
var s, n = "abc", 123 
var ( 
    a int
    b float32
) 
 
/*多变量赋值时，先计算所有相关值，然后再从左到右依次赋值*/
func main() {
    data, i := [3]int{0, 1, 2}, 0
    i, data[i] = 2, 100 // (i = 0) -> (i = 2), (data[0] = 100)
}

/*特殊只写变量 "_"，用于忽略值占位*/
func main() {
    _, s := test() 
    println(s) 
}

func test() (int, string) { 
    return 1, "abc"
}

/*编译器会将 未使用 的 局部变量 当做错误*/
func main() { 
     // Error: i declared and not used。(可使用 "_ = i" 规避) 
    i := 0
}

/*常量值必须是编译期可确定的数字、字符串、布尔值*/

 // 多常量初始化 
const x, y int = 1, 2
// 类型推断 
const s = "Hello, World!" 
// 常量组 
const ( 
    a, b = 10, 100 
    c bool = false
)

/*在常量组中，如不提供类型和初始化值，那么视作与上一常量相同*/
const ( 
    s = "abc" 
    x // x = "abc" 
)

func main() {
    // 未使用局部常量不会引发编译错误
    const x = "xxx" 
}

/*枚举: 
关键字 iota 定义常量组中从 0 开始按行计数 的 自增枚举值*/
const ( 
    Sunday = iota // 0 
    Monday // 1
    Tuesday // 2 
    Wednesday // 3 
    Thursday // 4 
    Friday // 5 
    Saturday // 6 
)

/*通过自定义类型来实现枚举类型限制*/
type Color int 

const ( 
    Black Color = iota 
    Red 
    Blue 
) 

func test(c Color) {
} 

func main() { 
    c := Black 
    test(c) 
    x := 1 
    // Error: cannot use x (type int) as type Color in function argument 
    test(x) 
    // 常量会被编译器自动转换
    test(1) 
}
```

### 基本类型

| 类型          | ⻓度    | 默认值 | 说明                                            |
| ------------- | ------ | ------ | ----------------------------------------------- |
| bool          | 1      | false  |                                                 |
| byte          | 1      | 0      | uint8                                           |
| rune          | 4      | 0      | Unicode Code Point, int32                       |
| int uint      | 4 或 8 | 0      | 32 或 64 位 int8, uint8 1 0 -128 ~ 127, 0 ~ 255 |
| int16, uint16 | 2      | 0      | -32768 ~ 32767, 0 ~ 65535                       |
| int32, uint32 | 4      | 0      | -21亿 ~ 21 亿, 0 ~ 42 亿                        |
| int64, uint64 | 8      | 0      |                                                 |
| float32       | 4      | 0.0    |                                                 |
| float64       | 8      | 0.0    |                                                 |
| complex64     | 8      |        |                                                 |
| complex128    | 16     |        |                                                 |
| uintptr       | 4 或 8 |        | 足以存储指针的 uint32 或 uint64 整数            |
| array         |        |        | 值类型                                          |
| struct        |        |        | 值类型                                          |
| string        |        | ""     | UTF-8 字符串                                    |
| slice         |        | nil    | 引用类型                                        |
| map           |        | nil    | 引用类型                                        |
| channel       |        | nil    | 引用类型                                        |
| interface     |        | nil    | 接⼝                                             |
| function      |        | nil    | 函数                                            |

### 引用类型

**引用类型** 包括 ***slice、map 和 channel***。
它们有 **复杂** 的内部结构，除了申请内存外，还需要 **初始化** 相关属性 ；

***内置函数 new***  计算 **类型大小**，为其分配 **零值内存**，**返回指针** ；

而 **make** 会被编译器翻译成 **具体** 的 **创建函数**，由其 **分配内存** 和 **初始化成员结构**，返回 **对象** 而非 **指针** ；

```go
// 提供初始化表达式。
a := []int {0, 0, 0} 
a[1] = 10 
// makeslice 
b := make([]int, 3) 
b[1] = 10 
c := new([]int) 
// Error: invalid operation: c[1] (index of type *[]int)
c[1] = 10 
```

### 类型转换

golang 不支持 **隐式类型转换**，即便是从 **窄向宽** ***转换*** 也不行

```go
var b byte = 100 
// Error: cannot use b (type byte) as type int in assignment 
var n int = b 
// 显式转换
var n int = int(b) 

//不能将其他类型当 bool 值使用
a := 100 
// Error: non-bool a (type int) used as if condition 
if a { 
    println("true")
}
```

### 字符串

**字符串** 是 **不可变值类型**，内部用 **指针** 指向 ***UTF-8 字节数组*** :

* 默认值是空字符串 ""。
* 用索引号访问某字节，如 s[i]。
* 不能用序号获取字节元素指针，&s[i] 非法。
* 不可变类型，无法修改字节数组。
* 字节数组尾部不包含 NULL。

```go
/* 使用索引号访问字符 (byte) */
s := "abc"
println(s[0] == '\x61', s[1] == 'b', s[2] == 0x63)

/* 连接跨行字符串时，"+" 必须在上一行末尾，否则导致编译错误 */
s := "Hello, " + 
"World!"

/* 支持用两个索引号返回子串。子串依然指向原字节数组，仅修改了 指针和长度属性 */
s := "Hello, World!" 
s1 := s[:5] // Hello 
s2 := s[7:] // World! 
s3 := s[1:5] // ello

/* 修改字符串，可先将其转换成 []rune 或 []byte，完成后再转换为 string。无论哪种转换，都会重新分配内存，并复制字节数组 */
func main() { 
    s := "abcd" 
    bs := []byte(s) 
    bs[1] = 'B' 
    println(string(bs)) 
    
    u := "电脑"
    us := []rune(u) 
    us[1] = '话' 
    println(string(us))
}

/* 用 for 循环遍历字符串时，也有 byte 和 rune 两种方式。*/
func main() { 
    s := "abc汉字" 
    // byte
    for i := 0; i < len(s); i++ {
        fmt.Printf("%c,", s[i]) 
    } 
    
    fmt.Println() 
    
    // rune
    for _, r := range s {
        fmt.Printf("%c,", r) 
    }
}
```

### 指针

支持指针类型 **\*T**，指针的指针 **\*\*T**，以及包含包名前缀的 **\*.T**。

* 默认值 nil，没有 NULL 常量。
* 操作符 **"&"** 取变量地址，**"\*"** 透过 **指针** 访问目标 **对象**。
* 不支持指针运算，不支持 **"->"** 运算符，直接用 **"."** 访问标成员。

```go
func main() { 
    type data struct{ a int }
    var d = data{1234} 
    var p *data 
    p = &d
    // 直接用指针访问目标对象成员，无须转换。
    fmt.Printf("%p, %v\n", p, p.a) 
    
    /*可以在 unsafe.Pointer 和任意类型指针间进行转换*/
    x := 0x12345678
    // *int -> Pointer 
    p := unsafe.Pointer(&x) 
    // Pointer -> *[4]byte 
    n := (*[4]byte)(p) 
    
    for i := 0; i < len(n); i++ {
        fmt.Printf("%X ", n[i])
    }
    /*返回局部变量指针是安全的，编译器会根据需要将其分配在 GC Heap 上*/
    func test() *int { 
        x := 100
        // 在堆上分配 x 内存。但在内联时，也可能直接分配在目标栈。
        return &x 
    }
}
    
/*将 Pointer 转换成 uintptr，可变相实现指针运算*/
func main() { 
     d := struct { 
           s string 
           x int 
      }{"abc", 100} 
    // *struct -> Pointer -> uintptr 
    p := uintptr(unsafe.Pointer(&d)) 
    // uintptr + offset
    p += unsafe.Offsetof(d.x) 
    // uintptr -> Pointer 
    p2 := unsafe.Pointer(p) 
    // Pointer -> *int 
    px := (*int)(p2) 
    // d.x = 200 
    *px = 200 
    fmt.Printf("%#v\n", d)
}
```

### 自定义类型

可将类型分为 **命名** 和 **未命名** 两大类。

**命名类型** 包括 ***bool、int、string*** 等，而 ***array、 slice、map*** 等和具体元素类型、长度等有关，属于 **未命名类型** 。

**具有相同声明的未命名类型被视为同一类型** :

* 具有 **相同基类型** 的 **指针**;
* 具有 **相同元素类型** 和 **长度** 的 ***array*** ;
* 具有 **相同元素类型** 的 ***slice***  ;
* 具有 **相同键值类型** 的 ***map*** ;  
* 具有 **相同元素类型** 和 **传送方向** 的 ***channel*** ;  
* 具有 **相同字段序列** ( ***字段名、类型、标签、顺序*** ) 的 **匿名 struct** ;
* **签名相同** ( ***参数 和 返回值***，不包括参数名称) 的 ***function*** ;  
* **方法集相同** ( ***方法名、方法签名相同***，和 **次序无关**) 的 ***interface*** ;

```go
/*可用 type 在 全局 或 函数内定义 新类型, 新类型不是原类型的别名，除拥有相同数据存储结构外，它们之间没有任何关系，不会持有原类型任何信息。
除非目标类型是未命名类型，否则必须显式转换*/
func main() { 
    type bigint int64 
    var x bigint = 100 
    println(x) 
}

x := 1234 
// 必须显式转换，除非是常量。
var b bigint = bigint(x)  
var b2 int64 = int64(b)
// 未命名类型，隐式转换。
var s myslice = []int{1, 2, 3}  
var s2 []int = s
```

### 表达式

#### 保留字

| **程序保留字**                            |
| ----------------------------------------- |
| ***break default func interface select*** |
| ***case defer go map struct***            |
| ***chan else goto package switch***       |
| ***const fallthrough if range type***     |
| ***continue for import return var***      |

```go
/* 单位运算 */
// AND 都为 1; 
0110 & 1011 = 0010 
// OR ⾄少一个为 1;
0110 | 1011 = 1111 
// XOR 只能一个为 1; 
0110 ^ 1011 = 1101 
// AND NOT 清除标志位
0110 &^ 1011 = 0100 

/*标志位操作*/
a := 0
// 0000100: 在 bit2 设置标志位
a |= 1 << 2 
// 1000100: 在 bit6 设置标志位
a |= 1 << 6  
// 0000100: 清除 bit6 标志位
a = a  &^ (1 << 6) 

/*不⽀持运算符重载。尤其需要注意，"++"、"--" 是语句⽽非表达式*/
n := 0
p := &n

// b := n++ // syntax error 
// if n++ == 1 {} // syntax error 
// ++n // syntax error
n++ 
*p++ // (*p)++

/* 没有 "~"，取反运算也用 "^" */
x := 1 
x, ^x // 0001, -0010
```

#### 初始化

初始化 **复合对象**，必须使用 **类型标签**，且 **左⼤括号** 必须在 **类型尾部**

```go
// var a struct { x int } = { 100 } // syntax error
// var b []int = { 1, 2, 3 }  // syntax error 
// c := struct {x int; y string}  // syntax error: unexpected semicolon or newline 
// { 
// }
var a = struct{ x int }{100} 
var b = []int{1, 2, 3}

/*初始化值以 "," 分隔。可以分多行，但最后一行必须以 "," 或 "}" 结尾*/
a := []int{
    1,
    2
}

a := []int{ 
    1, 
    2,
}

b := []int{ 
    1,
    2
}
```

### 控制流

#### IF

* 可省略条件表达式括号。
* 支持初始化语句，可定义代码块局部变量。
* 代码块左大括号必须在条件表达式尾部。

```go
x := 0

// 不⽀持三元操作符 "a > b ? a : b"
if n := "abc"; x > 0 {
    println(n[2]) 
} else if x < 0 {
    println(n[1]) 
} else {
    println(n[0]) 
}
```

#### For

支持三种循环方式, 包括 ***类 while*** 语法

```go
s := "abc"
// 常见的 for 循环，⽀持初始化语句
for i, n := 0, len(s); i < n; i++ { 
    println(s[i]) 
}

n := len(s)
for n > 0 {
    // 常见的 for 循环，⽀持初始化语句。
    println(s[n])
    // 替代 while (n > 0) {} // 替代 for (; n > 0;) {}
    n--
 }

for { 
    println(s) 
}

func length(s string) int { 
    println("call length.") 
    return len(s) 
}

func main() {
    s := "abcd"
    for i, n := 0, length(s); i < n; i++ {
        // 避免多次调用 length 函数。
        println(i, s[i])
    }
}

```

#### Range

类似 ***迭代器*** 操作，返回 ***索引, 值*** 或 ***键, 值***

```go
s := "abc"

// 忽略 2nd value，⽀持 string/array/slice/map。
for i := range s { 
    println(s[i]) 
}

// 忽略 index。
for _, c := range s {
    println(c) 
}

// 忽略全部返回值，仅迭代。
for range s { 

}

m := map[string] int {"a": 1, "b": 2}

// 返回 (key, value)
for k, v := range m { 
    println(k, v) 
}

/* range 会复制对象 */
a := [3]int {0, 1, 2}

for i, v := range a {
    if i == 0 {
    // index、value 都是从复制品中取出。
    // 在修改前，我们先修改原数组。
    a[1], a[2] = 999, 999
    fmt.Println(a)
    // 确认修改有效，输出 [0, 999, 999]。
}

a[i] = v + 100
// 使用复制品中取出的 value 修改原数组。
}

fmt.Println(a)
// => [100, 101, 102]。

/*建议改用引用类型，其底层数据不会被复制*/
s := []int{1, 2, 3, 4, 5}
for i, v := range s {
    if i == 0 {
        // 复制 struct slice { pointer, len, cap }。
        s = s[:3]
        s[2] = 100
        // 对 slice 的修改，不会影响 range。 // 对底层数据的修改。
    }
    println(i, v)
}
// 另外两种引用类型 map、channel 是指针包装，⽽不像 slice 是 struct
```

#### Switch

分支表达式可以是 **任意类型** ，不限于 **常量**。
可省略 break，默认自动终止

```go
x := []int{1, 2, 3} 
i := 2

switch i {
    case x[1]:
        println("a") 
    case 1, 3:
        println("b") 
    default:
        println("c")
}

//如需要继续下一分⽀，可使用 fallthrough，但不再判断条件
x := 10
switch x { 
    case 10:
        println("a")
        fallthrough 
    case 0:
        println("b") 
}

switch {
    case x[1] > 0:
        println("a") 
    case x[1] < 0:
        println("b") 
    default:
        println("c")
}

// 带初始化语句
switch i := x[2]; {
    case i > 0:
        println("a")
    case i < 0:
        println("b")
    default:
        println("c")
}

```

#### Goto, Break, Continue

支持在函数内 goto 跳转。标签名区分大小写，未使用标签引发错误

```go
func main() {
    var i int
    for {
        println(i)
        i++
        if i > 2 { goto BREAK } 
    }
    BREAK:
        println("break")
    EXIT: // Error: label EXIT defined and not used
}

//配合标签，break 和 continue 可在多级嵌套循环中跳出
func main() { 
L1:
    for x := 0; x < 3; x++ {
L2:
    for y := 0; y < 5; y++ {
        if y > 2 { continue L2 } 
        if x > 1 { break L1 }
        print(x, ":", y, " ")
    }
        println()
   }
}
       
```

**break** 可用于 ***for、switch、select*** ⽽ **continue** 仅能用于 ***for循环***

```go
x := 100

switch { 
    case x >= 0:
        if x == 0 { break }
        println(x) 
}
```

### 函数

#### 函数定义

不支持 **嵌套(nested)、重载(overload)** 和 **默认参数(default parameter)**

* 无需声明原型。
* 支持不定长变参。
* 支持多返回值。
* 支持命名返回参数。
* 支持匿名函数和闭包。

使用关键字 **func** 定义函数，**左大括号** 依旧不能另起一行

```go
// 类型相同的相邻参数可合并。
func test(x, y int, s string) (int, string) {
    n := x + y
    return n, fmt.Sprintf(s, n)
// 多返回值必须用括号。
}

/*函数是第一类对象，可作为参数传递。建议将复杂签名定义为函数类型，以便于阅读*/
func test(fn func() int) int { 
    return fn()
}

// 定义函数类型。
type FormatFunc func(s string, x, y int) string

func format(fn FormatFunc, s string, x, y int) string { 
    return fn(s, x, y) 
}

func main() { 
    // 直接将匿名函数当参数。
    s1 := test(func() int { return 100 }) 
    
    s2 := format(func(s string, x, y int) string { 
        return fmt.Sprintf(s, x, y) 
    }, "%d, %d", 10, 20)
    
    println(s1, s2)
}
//有返回值的函数，必须有明确的终⽌语句，否则会引发编译错误

```

#### 变参

**变参** 本质上就是 ***slice***。只能有一个，且必须是最后一个

```go
func test(s string, n ...int) string { 
    var x int 
    for _, i := range n {
        x += i
    }
    return fmt.Sprintf(s, x)
}

func main() { 
    println(test("sum: %d", 1, 2, 3))
    // 使用 slice 对象做变参时，必须展开
    s := []int{1, 2, 3}
    println(test("sum: %d", s...))
}
```

#### 返回值

不能用 **容器对象** 接收 **多返回值**。只能用多个变量，或 "\_" 忽略。

```go
func test() (int, int) {
    return 1, 2
}

func main() {
    x, _ := test() 
    println(x)
}

/*多返回值可直接作为其他函数调用实参*/
func test() (int, int) { 
    return 1, 2 
}
func add(x, y int) int { 
    return x + y
}
func sum(n ...int) int {
    var x int 
    for _, i := range n {
        x += i 
    }
    return x
}
func main() { 
    println(add(test()))
    println(sum(test()))
}

/*命名返回参数 可看做与形参类似的局部变量，最后由 return 隐式返回*/
func add(x, y int) (z int) {
    z = x + y
    return
}

func main() { 
    println(add(1, 2))
}

/*命名返回参数可被同名局部变量遮蔽，此时需要显式返回*/
func add(x, y int) (z int) { 
    {
        var z = x + y
        return z
    }
}

/*命名返回参数允许 defer 延迟调用通过闭包读取和修改*/
func add(x, y int) (z int) { 
    defer func() {
        z += 100
    }()

    z = x + y
    return
}

func main() { 
    println(add(1, 2)) 
}

/*显式 return 返回前，会先修改命名返回参数*/
func add(x, y int) (z int) {
    defer func() {
        println(z) // 输出: 203 
    }()
    z = x + y
    return z + 200
    // 执行顺序: (z = z + 200) -> (call defer) -> (ret)
}

func main() { 
    println(add(1, 2)) 
}
```

#### 匿名函数

匿名函数可 **赋值给变量**，做为 **结构字段** ，或者在 **channel** 里传送。

```go
// function variable 
fn := func() { println("Hello, World!") } 

fn()

// function collection 
fns := [](func(x int) int){ 
    func(x int) int { return x + 1 },
    func(x int) int { return x + 2 },
}
println(fns[0](100))

// function as field  
d := struct {
    fn func() string
}{
    fn: func() string { return "Hello, World!" }, 
}

println(d.fn())

// --- channel of function --- 
fc := make(chan func() string, 2)
fc <- func() string { return "Hello, World!" } 
println((<-fc)()
```

#### 延迟调用

关键字 **defer** 用于 **注册延迟** 调用。这些调用直到  **return 前** 才被执行，通常用于释放资源或错误处理。

```go
func test() error { 
    f, err := os.Create("test.txt")
    if err != nil { return err } 
    
    // 注册调用，⽽不是注册函数。必须提供参数，哪怕为空。
    defer f.Close()  
    
    f.WriteString("Hello, World!")
    return nil 
}

/*多个 defer 注册，按 FILO 次序执行。哪怕函数或某个延迟调用发⽣错误，这些调用依旧会被执行*/
func test(x int) { 
    defer println("a") 
    defer println("b") 
    defer func() { 
        // div0 异常未被捕获，逐步往外传递，最终终⽌进程。 
        println(100 / x) 
    }()
    defer println("c") 
}

func main() {
    test(0) 
}

/*延迟调用参数在注册时求值或复制，可用指针或闭包 "延迟" 读取*/
func test() {
    x, y := 10, 20 
    defer func(i int) { 
        println("defer:", i, y) // y 闭包引用
    }(x) // x 被复制

    x += 10 
    y += 100
    println("x =", x, "y =", y)
}
```

#### 错误处理

没有结构化异常，使用 **panic** 抛出错误，**recover** 捕获错误

```go
func test() {
    defer func() { 
        if err := recover(); err != nil {
            println(err.(string)) // 将 interface{} 转型为具体类型。 
        } 
    }()
    
    panic("panic error!")
}
```

由于 **panic、recover** 参数类型为 **interface{}** ，因此可抛出 **任何类型** 对象

```go
func panic(v interface{}) 
func recover() interface{}

/*延迟调用中引发的错误，可被后续延迟调用捕获，但仅最后一个错误可被捕获*/
func test() { 
    defer func() { 
        fmt.Println(recover())
    }() 
    
    defer func() { 
        panic("defer panic")
    }()
    
    panic("test panic") 
}

func main() { 
    test() 
}
```

捕获函数 **recover** 只有在 **延迟调用** 内直接调用才会终止错误，否则总是返回 nil。任何 **未捕获的错误** 都会沿调用堆栈向外传递

```go
func test() { 
    defer recover() // 无效！ 
    defer fmt.Println(recover()) // 无效！ 
    defer func() { 
        func() { 
            println("defer inner") 
            recover() // 无效！ 
        }()
    }() 
    panic("test panic")
}

func main() { 
    test() 
}

/*使用 延迟匿名函数 或下⾯这样都是有效的*/
func except() { 
    recover()
} 

func test() { 
    defer except()
    panic("test panic")
}
```

除用 **panic** 引发 **中断性错误** 外，还可返回 **error 类型** 错误对象来示函数调用状态

```go
type error interface { 
    Error() string
}
```

标准库 ***errors.New***  和 ***fmt.Errorf***  函数用于创建实现 **error** 接⼝的错误对象。通过判断错误 **对象实例** 来确定具体 **错误类型**

```go
var ErrDivByZero = errors.New("division by zero") 

func div(x, y int) (int, error) {
    if y == 0 { 
        return 0, ErrDivByZero 
    } 
    return x / y, nil 
}

func main() { 
    switch z, err := div(10, 0); err { 
    case nil: 
        println(z) 
    case ErrDivByZero:
        panic(err) 
    }
}
```

一般情况下导致关键流程出现 **不可修复性错误** 的使用 **panic**，其他使用 **error**

#### 初始化函数

* 每个源文件都可以定义一个或多个 **初始化函数**。
* 编译器不保证多个初始化函数 **执行次序**。
* 初始化函数在单一线程被调用，**仅执行一次**。
* 初始化函数在包所有全局变量初始化后执行。
* 在所有初始化函数结束后才执行 **main.main**。
* 无法调用 **初始化函数**

因为无法保证初始化函数执行顺序，因此全局变量应该直接用 var 初始化

```go
var now = time.Now() 

func main() { 
    fmt.Println("main:", int(time.Now().Sub(now).Seconds())) 
}

func init() { 
    fmt.Println("init:", int(time.Now().Sub(now).Seconds())) 
    w := make(chan bool) 
    
    go func() { 
        time.Sleep(time.Second * 3) 
        w <- true
    }()   
    <-w 
 }
```

#### 数据

#### Array

* **数组** 是 **值类型**，赋值 和 传参会 **复制** 整个 **数组**，而不是 **指针**;
* **数组长度** 必须是 **常量**，且是 **类型的组成** 部分 **[2]int 和 [3]int** 是不同类型;
* ⽀持 "=="、"!=" 操作符，因为 **内存** 总是被 **初始化** 过的;
* 指针数组 ***[n]\*T***，数组指针 ***\*[n]T***  ;

```go
/*未初始化元素值为 0*/
a := [3]int{1, 2} 
/*通过初始化值确定数组⻓度*/ 
b := [...]int{1, 2, 3, 4}
/*使用索引号初始化元素*/
c := [5]int{2: 100, 4:200} 

d := [...]struct { 
    name string 
    age uint8 
}{ 
    {"user1", 10}, // 可省略元素类型
    {"user2", 20}, // 最后一行的逗号
}

/*多维数组*/
a := [2][3]int{{1, 2, 3}, {4, 5, 6}} 
b := [...][2]int{{1, 1}, {2, 2}, {3, 3}} // 第 2 纬度不能用 "..."。

/*内置函数 len[数组⻓度] 和 cap[元素数量]*/
a := [2]int{} 
println(len(a), cap(a)) // 2, 2
```

#### Slice

***slice*** 不是 **数组** 或 **数组指针**。它通过 **内部指针** 和 **相关属性** 引用 **数组片段**，以实现 **变⻓** 方案 ;

* 引用类型。但自身是结构体，值拷贝传递。
* 属性 len 示可用元素数量，读写操作不能超过该限制。
* 属性 cap 表示最大扩张容量，不能超出数组限制。
* 如果 slice == nil，那么 len、cap 结果都等于 0。

```C
struct Slice {
    /* actual data */ 
    byte* array; 
    /* number of elements */ 
    uintgo len; 
     /* allocated number of elements */
    uintgo cap; 
};
```

```go
data := [...]int{0, 1, 2, 3, 4, 5, 6}
// [low : high : max]
slice := data[1:4:5] 

data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
//=> data[:6:8] ----+[0 1 2 3 4 5] ----+ 6 8 ----+省略 low
//=> data[5:] ----+[5 6 7 8 9]----+ 5 5 ----+省略 high、max
//=> data[:3] ----+[0 1 2]----+3 10 ----+省略 low、max
//=> data[:] ----+[0 1 2 3 4 5 6 7 8 9] ----+10 10 ----+全部省略

/*读写操作实际目标是底层数组，只需注意索引号的差别*/
data := [...]int{0, 1, 2, 3, 4, 5} 
s := data[2:4] 
s[0] += 100 
s[1] += 200 

fmt.Println(s)
//=> [102 203]
fmt.Println(data)
//=> [0 1 102 203 4 5]

/*可直接创建 slice 对象，自动分配底层数组*/
s1 := []int{0, 1, 2, 3, 8: 100} // 通过初始化表达式构造，可使用索引号。 
fmt.Println(s1, len(s1), cap(s1)) 
//=>  [0 1 2 3 0 0 0 0 100] 9 9

s2 := make([]int, 6, 8) // 使用 make 创建，指定 len 和 cap 值。 
fmt.Println(s2, len(s2), cap(s2))
//=> [0 0 0 0 0 0] 6 8

s3 := make([]int, 6) // 省略 cap，相当于 cap = len。 
fmt.Println(s3, len(s3), cap(s3))
//=> [0 0 0 0 0 0] 6 6

/*使用 make 动态创建 slice，避免了数组必须用常量做⻓度的⿇烦。还可用指针直接访问 底层数组，退化成 普通数组 操作*/
s := []int{0, 1, 2, 3}

// *int, 获取底层数组元素指针。
p := &s[2]  
*p += 100 
fmt.Println(s)
//=> [0 1 102 3]

/* [][]T 是指元素类型为 []T*/
data := [][]int{ 
    []int{1, 2, 3},
    []int{100, 200},
    []int{11, 22, 33, 44},
}

/*可直接修改 struct array/slice 成员*/
d := [5]struct { 
    x int
}{}

s := d[:] 
d[1].x = 10 
s[2].x = 20

fmt.Println(d) 
//=> [{0} {10} {20} {0} {0}]
fmt.Printf("%p, %p\n", &d, &d[0])
//=> 0x20819c180, 0x20819c180

```

#### reslice

**reslice** 是基于已有 slice 创建新 slice 对象，以便在 **cap** 允许范围内调整属性;

```go
s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

s1 := s[2:5] 
// => [2 3 4] 
s2 := s1[2:6:7]
// => [4 5 6 7] 
s3 := s2[3:6]
// => Error

/*新对象 依旧指向 原底层数组*/
s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

s1 := s[2:5]
// => [2 3 4] 
s1[2] = 100 
s2 := s1[2:6] 
//=> [100 5 6 7] 
s2[3] = 200

fmt.Println(s)
//=> [0 1 2 3 100 5 6 200 8 9]
```

#### append

向 slice **尾部** 添加数据，返回新的 slice 对象

```go
s := make([]int, 0, 5) 
fmt.Printf("%p\n", &s) 
//=> 0x210230000
s2 := append(s, 1) 
fmt.Printf("%p\n", &s2) 
//=> 0x210230040
fmt.Println(s, s2)
//=> [] [1]

/*底层是在 array[slice.high] 写数据*/
data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
s := data[:3] 
// 添加多个值
s2 := append(s, 100, 200)

fmt.Println(data)
//=> [0 1 2 100 200 5 6 7 8 9]
fmt.Println(s)
//=> [0 1 2]
fmt.Println(s2)
//=> [0 1 2 100 200]

/*一旦超出原 slice.cap 限制，就会重新分配底层数组，即便原数组并未填满*/
data := [...]int{0, 1, 2, 3, 4, 10: 0} 

s := data[:2:3]
// 一次 append 两个值，超出 s.cap 限制。 
s = append(s, 100, 200)

 // 重新分配底层数组，与原数组无关。
fmt.Println(s, data) 
//=> [0 1 100 200] [0 1 2 3 4 0 0 0 0 0 0]
 // 比对底层数组起始指针
fmt.Println(&s[0], &data[0])
//=>0x20819c180  0x20817c0c0
```

通常以 **2倍容量** 重新分配 **底层数组**。在大批量添加数据时，建议一次性分配足够大的空间，以减少内存分配和数据复制开销。或初始化足够长的 len 属性，改用索引号进行操作。及时释放不再使用的 slice 对象，避免持有过期数组，造成 GC 无法回收;

#### copy

***函数 copy*** 在两个 slice 间复制数据，**复制长度** 以 **len 小** 的为准。两个 slice 可指向同一底层数组，允许元素区间重叠

```go
data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

s := data[8:]
s2 := data[:5] 
// dst:s2, src:s
copy(s2, s) 

fmt.Println(s2)
//=> [8 9 2 3 4]
fmt.Println(data)
//=> [8 9 2 3 4 5 6 7 8 9]
```

应及时将所需数据 copy 到较小的 slice，以便释放 **超大号** 底层数组内存

#### Map

引用类型，哈希表。***键*** 必须是支持 **相等运算符** (==、!=) 类型，比如 ***number、string、 pointer、array、struct***，以及对应的 ***interface***。**值** 可以是 **任意类型**，没有限制 ;

```go
m := map[int]struct { 
    name string 
    age int 
}{ 
    1: {"user1", 10}, 
    2: {"user2", 20},
} 
println(m[1].name)

/*预先给 make 函数一个合理元素数量参数, 有助于提升性能 [避免后续操作时频繁扩张]*/
m := make(map[string]int, 1000)

m := map[string]int{
    "a": 1, 
} 
// 判断 key 是否存在
if v, ok := m["a"]; ok { 
    println(v)
} 
// 对于不存在的 key，直接返回 \0 , 不会出错
println(m["c"])  
// 新增或修改
m["b"] = 2  
// 删除 如果 key 不存在,不会出错
delete(m, "c")  
// 获取键值对数量 cap 无效
println(len(m)) 
// 迭代，可仅返回 key。随机顺序返回，每次都不相同 [和 golang 版本实现相关]
for k, v := range m { 
    println(k, v)
}

/*从 map 中取回的是一个 value 临时复制品，对其成员的修改是没有任何意义的*/
type user struct{ name string } 
//当 map 因扩张⽽重新哈希时，各键值项存储位置都会发⽣改变
m := map[int]user{  
    1: {"user1"},
}
m[1].name = "Tom"

 /*正确赋值方式*/
u := m[1]
u.name = "Tom" 
m[1] = u

m2 := map[int]*user{ 
    1: &user{"user1"},
} 
// 返回的是指针复制品。透过指针修改原对象是允许的
m2[1].name = "Jack" 
```

#### Struct

**值类型**，赋值和传参会复制全部内容。可用 "\_" 定义补位字段，支持指向 **自身类型** 的 **指针成员** ;

```go
type Node struct {
        _  int 
        id  int 
        data  *byte 
        next  *Node 
}

func main() { 
     n1 := Node{
            id: 1,
            data: nil,
     }

     n2 := Node{
         id: 2,
         data: nil,
         next: &n1,
     }
     
     /*顺序初始化必须包含全部字段，否则会出错*/
    type User struct { 
        name string
        age int
    } 
    u1 := User{"Tom", 20} 
    // Error: too few values in struct initializer
    u2 := User{"Tom"} 

    /*⽀持匿名结构，可用作结构成员或定义变量*/
    type File struct { 
        name string 
        size int 
        attr struct { 
            perm int 
            owner int
        } 
    }
    
    f := File{ 
        name: "test.txt", 
        size: 1025, 
        // attr: {0755, 1}, // Error: missing type in composite literal 
    }
    f.attr.owner = 1
    f.attr.perm = 0755
    
    var attr = struct { 
        perm int 
        owner int 
    }{2,0755}
    
    f.attr = attr
    
    /*⽀持 "=="、"!=" 相等操作符，可用作 map 键类型*/
    type User struct { 
        id int
        name string 
    } 
    
    m := map[User]int{ 
        User{1, "Tom"}: 100,
    }
    
    /*可定义字段标签，用反射读取。标签是类型的组成部分*/
    var u1 struct { name string "username" }
    var u2 struct { name string }
    u2 = u1
    //=> Error: cannot use u1 (type struct { name string "username" }) as type struct { name string } in assignment
    
    /*空结构 "节省" 内存，比如用来实现 set 数据结构，或者实现没有 "状态" 只有方法的 "静态类"*/
    var null struct{}
    set := make(map[string]struct{ }) 
    set["a"] = null
    
    /*匿名字段 不过是一种语法糖，从根本上说，就是一个与成员类型同名 (不含包名) 的字段。 被匿名嵌⼊的可以是任何类型，当然也包括 指针*/
    type User struct { 
        name string
    } 
    type Manager struct {
        User 
        title string
    }
    
    m := Manager{ 
        User: User{"Tom"}, // 匿名字段的显式字段名，和类型名相同
        title: "Administrator",
    }
 
    /*可以像普通字段那样访问匿名字段成员，编译器从外向内逐级查找所有层次的匿名字段， 直到发现目标或出错*/
    type Resource struct { 
        id int
    } 
    type User struct { 
        Resource
        name string
    } 
    type Manager struct {
        User 
        title string 
    }
    
    var m Manager
    m.id = 1 
    m.name = "Jack"
    m.title = "Administrator"
    
    /*外层同名字段会遮蔽嵌⼊字段成员，相同层次的同名字段也会让编译器报错，需要显示命名*/
    type Resource struct { 
        id int 
        name string 
    }
    
    type Classify struct { 
        id int 
    } 
    
    type User struct {
        // Resource.id 与 Classify.id 处于同一层次
        Resource 
        Classify
        // 遮蔽 Resource.name
        name string 
     }
     
    u := User{
        Resource{1, "people"},
        Classify{100},
        "Jack",
     }
    
    println(u.name) 
    //=> User.name: Jack 
    println(u.Resource.name) 
    //=> people 
    println(u.id) 
    //=> Error: ambiguous selector u.id 
    println(u.Classify.id) 
    //=> 100
    
    /* 不能同时嵌⼊某一类型和其指针类型，因为它们名字相同 */
    type Resource struct { 
        id int 
    } 
    
    type User struct { 
        *Resource
        // Resource  Error: duplicate field Resource 
        name string 
    }
    
    u := User{ 
        &Resource{1},
        "Administrator",
    }
    
    println(u.id)
    println(u.Resource.id)
    
    /*⾯向对象三⼤特征⾥，Go 仅支持封装，尽管匿名字段的内存布局和行为类似继承。没有 class 关键字，没有继承、多态等等*/
    type User struct { 
        id int 
        name string 
    } 
    
    type Manager struct { 
        User 
        title string 
    }
    
    m := Manager{User{1, "Tom"}, "Administrator"}
}
```

### 方法

#### 方法定义

**方法** 总是 ***绑定*** **对象实例**，并 **隐式** 将 **实例** 作为 ***第一实参 (receiver)***

* 只能为当前 **包内** 命名类型 定义方法
* 参数 receiver 可任意命名。如方法中未曾使用，可省略参数名。
* 参数 receiver 类型可以是 T 或  \*T。**基类型 T** 不能是 ***接口 或 指针***。
* 不支持方法重载，receiver 只是参数签名的组成部分。
* 可用实例 value 或 pointer 调用全部方法，编译器自动转换。

```go
/*没有构造和析构方法，通常用简单工⼚模式返回对象实例*/
type Queue struct { 
    elements []interface{} 
} 
// 创建对象实例
func NewQueue() *Queue {  
    return &Queue{make([]interface{}, 10)}
} 
// 省略 receiver 参数名
func (*Queue) Push(e interface{}) error { 
    panic("not implemented")
}
// receiver 参数名可以是 self、this 或其他
func (self *Queue) length() int { 
    return len(self.elements) 
}

/*receiver T 和 *T 的差别*/
type Data struct{
    x int
} 
// func ValueTest(self Data); 
func (self Data) ValueTest() {
    fmt.Printf("Value: %p\n", &self)
}
// func PointerTest(self *Data);
func (self *Data) PointerTest() {  
    fmt.Printf("Pointer: %p\n", self)
}

func main() { 
    d := Data{} 
    p := &d 
    fmt.Printf("Data: %p\n", p)
    
    d.ValueTest() 
    //=> ValueTest(d)
    d.PointerTest()
    //=> PointerTest(&d) 
    p.ValueTest()
    //=> ValueTest(*p) 
    p.PointerTest()
    //=> PointerTest(p) 
}

/*可以像字段成员那样访问匿名字段方法，编译器负责查找*/
type User struct {
    id int
    name string
}

type Manager struct { 
    User
}

// receiver = &(Manager.User) 
func (self *User) ToString() string { 
    return fmt.Sprintf("User: %p, %v", self, self) 
}

func main() {
    m := Manager{User{1, "Tom"}}
    
    fmt.Printf("Manager: %p\n", &m) 
    fmt.Println(m.ToString()) 
}
```

#### 匿名字段

```go
/*通过匿名字段，可获得和继承类似的复用能⼒。依据编译器查找次序，只需在外层定义同 名方法，就可以实现 override */
type User struct { 
    id int
    name string 
}

type Manager struct { 
    User 
    title string
}

func (self *User) ToString() string { 
    return fmt.Sprintf("User: %p, %v", self, self)
}

func (self *Manager) ToString() string { 
    return fmt.Sprintf("Manager: %p, %v", self, self)
}

func main() { 
    m := Manager{User{1, "Tom"}, "Administrator"}
    
    fmt.Println(m.ToString())
    fmt.Println(m.User.ToString())
}
```

#### 方法集

每个 **类型** 都有与之关联的 **方法集**，这会影响到 **接⼝** 实现规则

* 类型 T 方法集包含全部 receiver T 方法。
* 类型 \*T 方法集包含全部 receiver T + \*T 方法。
* 如类型 S 包含匿名字段 T，则 S 方法集包含 T 方法。
* 如类型 S 包含匿名字段 \*T，则 S 方法集包含 T + \*T 方法。
* 不管嵌⼊ T 或 \*T，\*S **方法集** 总是包含 T + \*T 方法

用实例 value 和 pointer 调用方法 (含匿名字段) 不受方法集约束，编译器总是查找全部方法，并自动转换 receiver 实参;

#### Method 表达式

根据调用者不同，方法分为两种表现形式：

```go
/* 前者称为 method value，后者 method expression 两者都可像普通函数那样赋值和传参

区别在于 method value 绑定实例，而 method expression 则须显式传参 */
instance.method(args...) ---> <type>.func(instance, args...)

type User struct { 
    id int 
    name string 
}

func (self *User) Test() {
    fmt.Printf("%p, %v\n", self, self)
}

func main() { 
    u := User{1, "Tom"} 
    u.Test() 
    mValue := u.Test 
    // 隐式传递 receiver 
    mValue() 
    // 显式传递 receiver
    mExpression := (*User).Test mExpression(&u) 
}


/* method value 会复制 receiver */
type User struct { 
    id int 
    name string 
}

func (self User) Test() { 
    fmt.Println(self)
}

func main() { 
    u := User{1, "Tom"} 
    // ⽴即复制 receiver，因为不是指针类型，不受后续修改影响
    mValue := u.Test 
    
    u.id, u.name = 2, "Jack"
    
    u.Test()
    //=> {2 Jack}
    mValue()
    //=> {1 Tom}
}

/*可依据方法集转换 method expression*/
type User struct { 
    id int
    name string 
} 

func (self *User) TestPointer() {
    fmt.Printf("TestPointer: %p, %v\n", self, self)
}

func (self User) TestValue() { 
    fmt.Printf("TestValue: %p, %v\n", &self, self) 
}

func main() { 
    u := User{1, "Tom"}
    fmt.Printf("User: %p, %v\n", &u, u)

    mv := User.TestValue 
    mv(u) 

    mp := (*User).TestPointer
    mp(&u) 
    // *User 方法集包含 TestValue。 
    mp2 := (*User).TestValue 
    // 签名变为 func TestValue(self *User) 实际依然是 receiver value copy 
    mp2(&u) 
} 
```

### 接口

**接口** 是一个或多个 **方法签名的集合**， **任何类型的方法集** 中只要拥有与之对应的 **全部方法**， 就表示它 "实现" 了该接口，无须在该类型上 **显式添加** 接口声明。
所谓 **对应方法**，是指有 **相同名称、参数列表 (不包括参数名) 以及返回值** 。当然，该类型还可以有其他方法；

#### 接口定义

* 接口命名习惯以 er 结尾，结构体。
* 接口只有方法签名，没有实现。
* 接口没有数据字段。
* 可在接口中嵌入其他接口。
* 类型可实现多个接口

```go
type Stringer interface { 
    String() string 
}

type Printer interface { 
    // 接口嵌⼊
    Stringer  
    Print()
}

type User struct {
    id int 
    name string 
}

func (self *User) String() string { 
    return fmt.Sprintf("user %d, %s", self.id, self.name) 
}

func (self *User) Print() { 
    fmt.Println(self.String()) 
} 

func main() { 
    var t Printer = &User{1, "Tom"} 
    // *User 方法集包含 String、Print
    t.Print() 
}
```

**空接⼝** ***interface{}*** 没有任何方法签名，也就意味着 **任何类型** 都实现了空接口。
其作用类似面向对象语言中的 **根对象 object**

```go
func Print(v interface{}) { 
    fmt.Printf("%T: %v\n", v, v)
} 

func main() {
    Print(1) 
    Print("Hello, World!")
}
```

**匿名接口** 可用作变量类型，或结构成员

```go
type Tester struct { 
    s interface {
        String() string 
    } 
} 

type User struct { 
    id int 
    name string 
} 

func (self *User) String() string { 
    return fmt.Sprintf("user %d, %s", self.id, self.name) 
} 

func main() { 
    t := Tester{&User{1, "Tom"}}
    fmt.Println(t.s.String()) 
}
```

#### 执行机制

**接口对象** 由 ***接口表 (interface table)*** **指针** 和 **数据指针** 组成

```C
struct Iface { 
    Itab* tab;
    void* data;
};

struct Itab { 
    InterfaceType* inter;
    Type* type;
    void (*fun[])(void);
};
```

**接口表** 存储 **元数据信息**，包括 ***接口类型、动态类型，以及实现接口的方法指针*** 无论是 **反射** 还是通过 **接口调用** 方法，都会用到这些信息。 **数据指针** 持有的是目标对象的 **只读** 复制品，复制完整对象或指针 ;

```go
type User struct { 
    id int 
    name string 
} 

func main() { 
    u := User{1, "Tom"} 
    var i interface{} = u

    u.id = 2 
    u.name = "Jack" 

    fmt.Printf("%v\n", u)
    //=> {2 Jack} 
    fmt.Printf("%v\n", i.(User))
    //=> {1 Tom}
}
```

**接口** 转型返回 ***临时对象***，只有使用 **指针** 才能修改其 **状态** ;

```go

type User struct { 
    id int 
    name string 
} 

func main() { 
    u := User{1, "Tom"}
    var vi, pi interface{} = u, &u 
 
    vi.(User).name = "Jack"
     //=> Error: cannot assign to vi.(User).name
    pi.(*User).name = "Jack"
    
    fmt.Printf("%v\n", vi.(User)) 
    //=> {1 Tom}
    fmt.Printf("%v\n", pi.(*User))
    //=> &{1 Jack}
}
```

##### 接口转换

利用 **类型推断**，可判断 **接口对象** 是否是某个具体的 **接口 或 类型**

```go
type User struct { 
    id int 
    name string
}

func (self *User) String() string { 
    return fmt.Sprintf("%d, %s", self.id, self.name)
}

func main() { 
    var o interface{} = &User{1, "Tom"}
    
    if i, ok := o.(fmt.Stringer); ok {
        fmt.Println(i) 
    } 
    
    u := o.(*User) 
    // u := o.(User) // panic: interface is *main.User, not main.User 
    fmt.Println(u)
}
```

还可用 switch 做 **批量类型** 判断，不支持 **fallthrough**

```go
func main() { 
    var o interface{} = &User{1, "Tom"}

    switch v := o.(type) { 
        case nil: // o == nil 
            fmt.Println("nil")
        case fmt.Stringer: // interface
            fmt.Println(v) 
        case func() string: // func 
            fmt.Println(v()) 
        case *User: // *struct 
            fmt.Printf("%d, %s\n", v.id, v.name) 
        default: 
            fmt.Println("unknown") 
    } 
}

/*超集接口对象可转换为子集接口，反之出错*/
type Stringer interface { 
    String() string 
}

type Printer interface { 
    String()  string 
    Print()
}

type User struct { 
    id int 
    name string 
}

func (self *User) String() string { 
    return fmt.Sprintf("%d, %v", self.id, self.name)
} 

func (self *User) Print() { 
    fmt.Println(self.String())
}

func main() { 
    var o Printer = &User{1, "Tom"}
    var s Stringer = o 
    
    fmt.Println(s.String()) 
}
```

#### 接口技巧

```go
    /*让编译器检查，以确保某个类型实现接口*/
    var _ fmt.Stringer = (*Data)(nil)
  
    /*函数直接 "实现" 接口能省不少事 */
    type Tester interface { 
        Do() 
    }

    type FuncDo func() 
    
    func (self FuncDo) Do() {
        self() 
    }
    
    func main() {
        var t Tester = FuncDo(func() { println("Hello, World!") }) 
        t.Do()
    }
```

### 并发

#### Goroutine

Go 在 **语言层面** 对 **并发编程** 提供支持，一种类似 **协程**，称作 **goroutine 机制** ；
只需在函数调用语句前添加 ***go*** 关键字，就可创建 **并发执行单元** ；
调度器会自动动将其安排到合适的系统线程上执行；
***goroutine*** 是一种非常 **轻量级** 的实现，可在单个进程里执行成千上万的并发任务 ；

```go
    go func() {
        println("Hello, World!")
    }()
```

**调度器** 不能保证多个 ***goroutine*** 执行 **次序**，且进程退出时不会等待它们 **结束** 。
**进程** 启动后仅允许 **一个系统线程** 服务于 ***goroutine*** 。
可使用 **环境变量** 或 **标准库** 函数 ***runtime.GOMAXPROCS*** 修改，让 **调度器** 用多个线程实现 **多核并行**，而不仅仅是 **并发** 。

```go
func sum(id int) { 
    var x int64
    for i := 0; i < math.MaxUint32; i++ { 
        x += int64(i) 
    } 
    println(id, x)
}

func main() {
    wg := new(sync.WaitGroup)
    wg.Add(2) 
    
    for i := 0; i < 2; i++ {
        go func(id int) { 
            defer wg.Done()
            sum(id)
         }(i)
    wg.Wait()
   }
}
```

调用 ***runtime.Goexit***  将立即终止当前 ***goroutine*** 执行，**调度器** 确保所有已注册 ***defer 延迟调用*** 被执⾏

```go
func main() { 
    wg := new(sync.WaitGroup) 
    wg.Add(1) 
    go func() {
        defer wg.Done() 
        defer println("A.defer")
        func() { 
            defer println("B.defer")
            // 终⽌当前 goroutine 
            runtime.Goexit() 
            // 不会执⾏ 
            println("B") 
        }() 
        // 不会执⾏
        println("A") 
    }() 
    wg.Wait()
}
```

和协程 **yield** 作用类似，**Gosched** 让出底层线程，将当前 ***goroutine*** 暂停，放回队列等待下次被调度执行

```go
func main() {
    wg := new(sync.WaitGroup) 
    wg.Add(2)
    
    go func() { 
        defer wg.Done() 
        for i := 0; i < 6; i++ { 
            println(i) 
            if i == 3 { 
                runtime.Gosched() 
            } 
        }
    }()
    
   go func() { 
        defer wg.Done()
        println("Hello, World!")
    }() 
    wg.Wait()
 }
```

#### Channel

引用类型 ***channel*** 是 **CSP 模式** 的具体实现，用于多个 ***goroutine*** 通讯。其内部实现了同步，确保并发安全。

```go
func main() {
    // 数据交换队列 
    data := make(chan int) 
    // 退出通知 
    exit := make(chan bool) 
    // 从队列迭代接收数据; 直到 close 
    go func() { 
        for d := range data {
            fmt.Println(d)
        } 
        fmt.Println("recv over.")
        // 发出退出通知 
        exit <- true 
    }() 
    // 发送数据
    data <- 1  
    data <- 2 
    data <- 3 
    // 关闭队列
    close(data) 
    
    fmt.Println("send over.") 
     // 等待退出通知
    <-exit
}

```

异步方式通过判断 **缓冲区** 来决定是否阻塞。如果缓冲区已满，发送被阻塞；缓冲区为空， 接收被阻塞。
通常情况下，异步 ***channel*** 可减少排队阻塞，具备更高的效率。
但应该考虑 **使用指针规避大对象拷贝**，将多个元素打包，减小缓冲区大小等

```go
func main() { 
    // 缓冲区可以存储 3 个元素 
    data := make(chan int, 3) 
    exit := make(chan bool) 
    
    // 在缓冲区未满前，不会阻塞
    data <- 1  
    data <- 2 
    data <- 3
  
    go func() { 
         // 在缓冲区未空前，不会阻塞
         for d := range data {
            fmt.Println(d)
         }
        exit <- true 
    }()
    
    // 如果缓冲区已满，阻塞
    data <- 4  
    data <- 5
    
    close(data)
    <-exit
}
```

缓冲区是内部属性，并非类型构成要素

```go
var a, b chan int = make(chan int), make(chan int, 3)

/*可以将 channel 隐式转换为单向队列，只收或只发*/
c := make(chan int, 3) 
// send-only 
var send chan<- int = c 
 // receive-only
var recv <-chan int = c

/*如果需要同时处理多个 channel，可使用 select 语句。它随机选择一个可用 channel 做 收发操作，或执⾏ default case*/
func main() { 
    a, b := make(chan int, 3), make(chan int)

    go func() { 
        v, ok, s := 0, false, "" 

        for { 
            select { 
                case v, ok = <-a: s = "a" 
                case v, ok = <-b: s = "b" 
            } 

            if ok { 
                fmt.Println(s, v)
            } else { 
                os.Exit(0)
            } 

        }  
    }()


    for i := 0; i < 5; i++ {
        // 随机选择可用 channel 发送数据
        select { 
           case a <- i: 
           case b <- i: 
        } 
    } 

    close(a)
    //没有可用 channel，阻塞 main goroutine
    select {}
}

/*用简单 工⼚模式 打包并发任务和 channel*/
func NewTest() chan int {
    c := make(chan int) 
    rand.Seed(time.Now().UnixNano()) 
    
    go func() { 
        time.Sleep(time.Second) 
        c <- rand.Int() 
    }()
    return c 
} 

func main() {
    t := NewTest()
    // 等待 goroutine 结束返回
    println(<-t) 
}

/*用 channel 实现信号量 (semaphore) */
func main() { 
    wg := sync.WaitGroup{} 
    wg.Add(3)
    
    sem := make(chan int, 1) 
    
    for i := 0; i < 3; i++ {
        go func(id int) {
            defer wg.Done() 
            // 向 sem 发送数据 阻塞或者成功
            sem <- 1  

            for x := 0; x < 3; x++ {
                fmt.Println(id, x)
            } 
            // 接收数据，使得其他阻塞 goroutine 可以发送数据
            <-sem 
        }(i) 
    } 
    wg.Wait()
 }
 
 /*用 closed channel 发出退出通知 */
func main() {
    var wg sync.WaitGroup
    quit := make(chan bool)

    for i := 0; i < 2; i++ { 
        wg.Add(1) 
        go func(id int) { 
            defer wg.Done() 

            task := func() { 
                println(id, time.Now().Nanosecond())
                time.Sleep(time.Second) 
            }


            for { 
                select { 
                    // closed channel 不会阻塞，因此可用作退出通知
                    case <-quit: 
                        // 执⾏正常任务
                        return 
                    default: 
                        task()
                } 
            }
         }(i)
     }

    // 让测试 goroutine 运⾏一会
    time.Sleep(time.Second * 5) 
    // 发出退出通知
    close(quit) 
    wg.Wait()
}

/*用 select 实现超时 (timeout) */
func main() { 
    w := make(chan bool)
    c := make(chan int, 2) 
    
    go func() { 
        select { 
            case v := <-c: fmt.Println(v) 
            case <-time.After(time.Second * 3): fmt.Println("timeout.")
        }
        w <- true 
    }() 
    
    // c <- 1 注释掉，引发 timeout。
    <-w 
}

/*channel 是第一类对象，可传参 (内部实现为指针) 或者作为结构成员*/
type Request struct {
     data []int
     ret chan int
}

func NewRequest(data ...int) *Request {
    return &Request{ data, make(chan int, 1) }
}

func Process(req *Request) { 
    x := 0 
    
    for _, i := range req.data { 
        x += i
    }
    
    req.ret <- x 
 }
 
 func main() { 
    req := NewRequest(10, 20, 30) 
    Process(req) 
    fmt.Println(<-req.ret)
 }
```
