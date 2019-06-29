package main

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type Sample struct {
	Name string
}

func (*Sample) Close() error {
	return nil
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
	fmt.Printf("%d\n", 100.0/a)
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
		上述描述不准确，主要还是得看初始化时候的层级
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

func isSliceOrArrayWithKind(args interface{}, kind reflect.Kind) (b bool) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			} else if s, ok := r.(string); ok {
				panic(s)
			} else {
				b = false
			}
		}
	}()

	t := reflect.TypeOf(args)
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return false
	}
	return t.Elem().Kind() == kind
}

func TestIsSliceOrArrayWithKind(t *testing.T) {
	a := []string{"h", "e"}
	if !isSliceOrArrayWithKind(a, reflect.String) {
		t.Fatal("is not")
	}

	//上面的很傻比，语法本身就支持，开始没想到
	var b interface{} = a
	if _, ok := b.([]string); ok {
		fmt.Println("is string array")
	}
}

func TestTime(*testing.T) {
	t := time.Now() //这个里面默认给到的是当前时区
	label := "2006-01-02 15:04:05"
	ts := t.Local().Format(label) //这里即使不给Local()，出来的也是当前时区的format，因为Now()里面已经设置了时区

	// t2, _ := time.Parse(label, ts) //这里如果不使用ParseInLocation，转出来的是UTC时区
	//ts2 := t2.Format(label) //这里的结果长得是和上面的ts一样，那也只是在UTC时区的长的一样的时间，t和t2完全不是一码事

	//只有这样才能Parse到目的时区，最终即使不给Local()，出来的也是相同的时间
	t2, _ := time.ParseInLocation(label, ts, time.Local)
	ts2 := t2.Local().Format(label)

	fmt.Println(ts, ts2)

	//综上所述
	//在我们一定要基于当前时区转换的情况下，为了保险起见，请Format之前一定要加Local()，Parse时候一定要用ParseInLocation(...time.Local)
}

func doTransferInterface(e io.Closer) interface{} {
	if e == nil {
		fmt.Println("e is nil")
	} else {
		fmt.Println("e is not nil")
	}
	var s *Sample
	return s
}

//测试接口作为参数和返回值nil最终判断
func TestTransferNilInterface(*testing.T) {
	// var x *Sample
	// a := doTransferInterface(x) //上面两行的话输出e is not nil
	a := doTransferInterface(nil) //这行的话输出e is nil
	if a == nil {
		fmt.Println("a is nil")
	} else {
		fmt.Println("a is not nil")
	}
}

func TestContinue(*testing.T) {
	path := []string{"a", "b", "c", "d"}
	fType := "c"

fields:
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		for j, e := range path {
			fmt.Println("   ", j, " ", e)
			if fType == e {
				fmt.Println("continue fields")
				continue fields
			}
		}
	}
}

func collectIntCombines(groups *[][]int, gIndex int, combine []int, combines *[][]int) {
	if gIndex < 0 || gIndex >= len(*groups) {
		return
	}
	for _, op := range (*groups)[gIndex] {
		cb := append(combine, op)
		/*
				根据下面这句的输出可以发现
				1->0xc000082240,0x124fbc8,[1],[]
				1->0xc000082250,0xc000082240,[1 11],[1]
				1->0xc0000bc0c0,0xc000082250,[1 11 21],[1 11]
				1->0xc0000bc0c0,0xc0000bc0c0,[1 11 21 31],[1 11 21]
				2->collect:0xc0000bc0c0,[1 11 21 31]
				1->0xc0000bc0c0,0xc0000bc0c0,[1 11 21 32],[1 11 21]
				2->collect:0xc0000bc0c0,[1 11 21 32]
				final resu.t:[[1 11 21 32] [1 11 21 32]] // 注意结果却是两个重复的数据
				注意第四句，打印出来的地址一样，但是内容不一样。
				这说明一个问题，俩数组存储内容的区域用的一块地，但是数组对象基本信息用的地不是一个地
				cb的基本信息认为有4个元素，combine认为有3个，这样以来，修改了两者公用的那块区域都会影响到另外一方。

			最特么奇怪的是这不是必然情况，偏偏到第三层递归时候append出现这种情况，注释掉下面的41 42即可看效果
			解释：因为原数组的cap还够用就不会开辟新内存，https://stackoverflow.com/a/28143457，需要开辟新内存的时机也讲究
			2,4,8,16,32如此类推，所以上面第三层出现这情况了
		*/
		fmt.Printf("1->%p,%p,%v,%v\n", cb, combine, cb, combine)
		if gIndex == len(*groups)-1 { //如果是最后一个数组，则可以完成一个组合且收集到combines里
			fmt.Printf("2->collect:%p,%v\n", cb, cb)
			cp := make([]int, len(cb))
			copy(cp, cb) //必须copy一份，否则之后的循环可能会影响已记录部分
			*combines = append(*combines, cp)
			//下面是错误示范
			// *combines = append(*combines, cb)
			continue
		}
		collectIntCombines(groups, gIndex+1, cb, combines)
	}
}

func TestCollectIntCombines(t *testing.T) {
	groups := [][]int{
		{
			1,
		},
		{
			11,
		},
		{
			21,
		},
		{
			31, 32,
		},
	}

	combines := [][]int{}
	collectIntCombines(&groups, 0, []int{}, &combines)

	t.Log(combines)
	t.Log("一共", len(combines), "个结果")
}

type ShopORM struct {
	AccountID       string
	AddressCity     string
	AddressProvince string
}

func (s ShopORM) TableName() string {
	s.AccountID = "100"
	return "card_shops"
}

func (s *ShopORM) TableName2() string {
	s.AccountID = "200"
	return "card_shops2"
}

func TestValuePointerMethod(t *testing.T) {
	s := ShopORM{}
	fmt.Println(s.TableName())
	fmt.Println(s.AccountID)
	fmt.Println(s.TableName2())
	fmt.Println(s.AccountID)
}

func TestTest(t *testing.T) {
	ss := []ShopORM{
		ShopORM{
			AccountID: "1",
		},
		ShopORM{
			AccountID: "2",
		},
	}
	ss2 := []*ShopORM{}
	// for _, s := range ss {
	// 	ss2 = append(ss2, &s)
	// }
	var s ShopORM
	for i := 0; i < len(ss); i++ {
		s = ss[i]
		ss2 = append(ss2, &s)
	}

	fmt.Println(ss2)
}

func TestAppend(t *testing.T) {
	a := []int{1, 2}

	b := append(a, 3)
	c := append(a, 4)
	fmt.Printf("a:%p->%v\nb:%p->%v\nc:%p->%v\n", a, a, b, b, c, c)

	a = append(a, 3)
	d := append(a, 4)
	e := append(a, 5)
	fmt.Printf("a:%p->%v\nd:%p->%v\ne:%p->%v\n", a, a, d, d, e, e)

}
