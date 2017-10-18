package main

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	var a uint64 = 250

	/*
		输出
			0x0
			%!f(uint64=00)
			2.5
			2.50
		这说明，前两者的运算的结果类型最终都是依着a的类型来，100.0并未像其他语言一样变成浮点数运算
	*/
	fmt.Printf("%#v\n", 100.0/a)
	fmt.Printf("%.2f\n", 100.0/a)
	fmt.Printf("%#v\n", float64(a)/100)
	fmt.Printf("%.2f\n", float64(a)/100)

	//var b float64 = 750.0
	//下面两句会编译摆错，浮点数类型和整型不能直接做做运算
	//fmt.Printf("%#v\n", b/a)
	//fmt.Printf("%.2f\n", b/a)
}

func do(o int) (int, error) {
	return o, nil
}

func TestStmt(t *testing.T) {
	r := 1

	fmt.Printf("1-- %#v,%p\n", r, &r)
	r, err := do(2) //虽说:=了，但是不是一个新的r
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("2-- %#v,%p\n", r, &r)
	if true {
		r, err = do(5)   //这个修改的是上面的
		r, err2 := do(7) //这个不像上面那个:=，这里会局部初始化一个新的同步的r变量
		if err2 != nil {
			t.Fatal(err2)
		}
		fmt.Printf("3-- %#v,%p\n", r, &r)
	}
	fmt.Printf("4-- %#v,%p\n", r, &r)

	/*
		输出
		1-- 1,0xc420016f88
		2-- 2,0xc420016f88
		3-- 7,0xc420016ff0
		4-- 5,0xc420016f88
		结论就是
		因为多返回值而必须使用:=的情况下，若是同级语句块的话，则不是新的变量生成，否则就是新变量生成
	*/
}
