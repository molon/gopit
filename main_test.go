package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

type Sample struct {
	Name string
}

func TestDivision(t *testing.T) {
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

func doForTestStmt(o int) (int, error) {
	return o, nil
}

func TestStmt(t *testing.T) {
	r := 1

	fmt.Printf("1-- %#v,%p\n", r, &r)
	r, err := doForTestStmt(2) //虽说:=了，但是不是一个新的r
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("2-- %#v,%p\n", r, &r)
	if true {
		r, err = doForTestStmt(5)   //这个修改的是上面的
		r, err2 := doForTestStmt(7) //这个不像上面那个:=，这里会局部初始化一个新的同步的r变量
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

func doForTestRecover1() (sp *Sample, err error) {
	var s *Sample
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered doForTestRecover1: ", r)
		}
		s = &Sample{ //这时候的s再设置也没屌用，下面的return s, nil又不会执行到
			Name: "hello2",
		}
		err = fmt.Errorf("panic err in defer") //这时候设置的err还是会认的，因为返回参数设置了同名
	}()

	//这里的panic会被捕获，然后此函数的返回值会是nil，基本可以认为单纯var x xxtype的值
	panic(fmt.Errorf("Panic doForTestRecover1"))

	fmt.Println("will return doForTestRecover1")
	s = &Sample{
		Name: "hello",
	}
	return s, nil
}

func justDoRecover() {
	if r := recover(); r != nil {
		fmt.Println("Recovered justDoRecover: ", r)
	}
	fmt.Println("End justDoRecover")
}

func doForTestRecover2() {
	f := func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered f: ", r)
		}
		fmt.Println("End f")
	}
	defer func() {
		//注意这俩里面即使panic了，defer的内容也一定会执行的
		//注意，这个方法里面的recover即使执行了，也不会捕获下面的panic
		//所以recover只能捕获同级的或者子级的panic
		f()
		justDoRecover()
	}()
	panic(fmt.Errorf("Panic doForTestRecover2"))
}

func TestRecover(t *testing.T) {
	fmt.Println(doForTestRecover1())
	doForTestRecover2()
}

//just test 不是坑
func TestFilePathGlob(t *testing.T) {
	m, err := filepath.Glob("./*.go")
	if err != nil {
		t.Fatal(m)
	}
	fmt.Println(m)
}

func editMap(m map[string]interface{}) {
	m["edit"] = "123"
}

//Map相关
func TestMap(t *testing.T) {
	m := make(map[string]interface{})
	m["1"] = "haha1"
	m["2"] = "haha2"
	fmt.Println(m)
	m["1"] = nil
	fmt.Println(m)
	if s, ok := m["1"]; ok {
		fmt.Println("1 exists: ", s)
	}
	delete(m, "2")
	fmt.Println(m)
	//上述结论是 map 里是可存储nil的，删除必须使用delete

	editMap(m)
	fmt.Println(m)

	val, ok := m["edit"].(string)
	if ok {
		fmt.Println("exists and is string:", val)
	} else {
		fmt.Println("not exists or not string")
	}
}

func TestFmt(t *testing.T) {
	//中间默认就会产生空格
	fmt.Println("1", "2")
}

func inc(a *int, delta int) int {
	*a = *a + delta
	return *a
}

func doDefer() int {
	i := 1
	defer func(a int) {
		i = i + a
		fmt.Println("inside:", i)
	}(inc(&i, 1)) //作为参数的inc(&i, 1)已经提前执行了，所以这句后i已经变为2了

	return inc(&i, 5) //此句执行后结果为7，defer里的打印出来是 inside:9
	return i + 5      //此句执行后结果为7，defer里的打印出来是 inside:4

	//上面两句，前者是在return这句执行后对i进行了修改，后者没有。
	//根据结果可以判断出，return后的语句会比defer里的玩意先执行。
}

func TestDefer(*testing.T) {
	fmt.Println("doDefer:", doDefer())
}

func TestSet(*testing.T) {
	a := Sample{Name: "molon"}
	b := a
	c := b
	fmt.Printf("a:%p\n", &a)
	fmt.Printf("b:%p\n", &b)
	fmt.Printf("c:%p\n", &c)
	//上述打印出来的结果不一样，说明只要赋值非指针就是copy
}

func TestSlice(*testing.T) {
	a := []string{"a", "b", "c", "d", "e"}

	fmt.Println(a[0])
	//fmt.Println(a[3]) //越界了会panic
	fmt.Println(a[0:2])
	fmt.Println(a[0:5]) //注意5虽然看起来越界了，但是因为最终不会拿取这个位置，所以没事
	//fmt.Println(a[0:6])
}
