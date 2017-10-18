# gopit
golang坑的一些测试

大坑集锦

下面出现的S都表示为一个结构体类型

---
vscode的debug配置里要加上env:{}它才能找到GOPATH

---
var a *string --> nil
var a string --> 非nil
var a *S --> nil
var a S -> 非nil
上面的很好理解，但是下面的注意
var a []string --> nil
a := []string{} --> 非nil
var a map[int]string --> nil
a := map[int]string{} --> 非nil
这个就是为什么json.Marshal出来的结果前者是null，后者是[]。
最终判断是否为nil的方式
func isInvalidOrNil(e interface{}) (b bool) {
	v := reflect.ValueOf(e)
	valid := v.IsValid()
	if !valid {
		return true
	}

	if e == nil {
		return true
	}

	//支持IsNil的话就输出结果，否则就是false
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
	return v.IsNil()
}

---
a := S{}
b := a
这样的话实际上b是新开辟了内存空间的

---
a := map[int]S{}
a[0] = S{"str"}
a[0].xxx = "str2"
//上面这样不行，因为定义的map value不是指针，不能直接操作，很奇葩，原因暂时不知道
//value改为指针即可

---
a:= []S{........}
for index,e := range a {
	//这里的e是开辟了一份内存空间，例如0x4a4a4c这块
	//然后每次循环都是从a的对应index哪里copy进来到这块内存
//所以切记，存储e的指针地址一般毫无意义，因为其内容在下次循环必定会被改变，必须存储的话，请自行copy一份
//例如ee:=e都可
}

---
https://stackoverflow.com/a/28143457
func collectIntCombines(groups *[][]int, gIndex int, combine []int, combines *[][]int) {
	if gIndex < 0 || gIndex >= len(*groups) {
		return
	}
	for index, op := range (*groups)[gIndex] {
		if gIndex == len(*groups)-1 && index == len((*groups)[gIndex])-1 {
			fmt.Printf("%p,%v\n", &(*combines)[0], (*combines)[0])
		}
		cb := append(combine, op)
		if gIndex == len(*groups)-1 && index == len((*groups)[gIndex])-1 {
			fmt.Printf("%p,%v\n", &(*combines)[0], (*combines)[0])
		}
		/*
				根据下面这句的输出可以发现
				0xc420354340,0x1d8f318,[1],[]
				0xc420354350,0xc420354340,[1 11],[1]
				0xc4202142c0,0xc420354350,[1 11 21],[1 11]
				0xc4202142c0,0xc4202142c0,[1 11 21 31],[1 11 21]
				0xc4202143e0,[1 11 21 31]
				0xc4202143e0,[1 11 21 32]
				0xc4202142c0,0xc4202142c0,[1 11 21 32],[1 11 21]
				注意第四句，打印出来的地址一样，但是内容不一样。
				这说明一个问题，俩数组存储内容的区域用的一块地，但是数组对象基本信息用的地不是一个地
				cb的基本信息认为有4个元素，combine认为有3个，这样以来，修改了两者公用的那块区域都会影响到另外一方。

			最特么奇怪的是这不是必然情况，偏偏到第三层递归时候append出现这种情况，注释掉下面的41 42即可看效果
			解释：因为原数组的cap还够用就不会开辟新内存，https://stackoverflow.com/a/28143457，需要开辟新内存的时机也讲究
			2,4,8,16,32如此类推，所以上面第三层出现这情况了
		*/
		fmt.Printf("%p,%p,%v,%v\n", cb, combine, cb, combine)
		if gIndex == len(*groups)-1 { //如果是最后一个数组，则可以完成一个组合且收集到combines里
			//cp := make([]int, len(cb))
			//copy(cp, cb) //必须copy一份，否则之后的循环可能会影响已记录部分
			//*combines = append(*combines, cp)
			//下面是错误示范
			*combines = append(*combines, cb)
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
		//{
		//	41, 42,
		//},
	}

	combines := [][]int{}
	collectIntCombines(&groups, 0, []int{}, &combines)

	t.Log(combines)
	t.Log("一共", len(combines), "个结果")
}

-------------
func TestFormat(t *testing.T) {
	var a uint64 = 250
	fmt.Printf("%#v\n", a/100.0)
	fmt.Printf("%.2f\n", a/100.0)
	fmt.Printf("%#v\n", float64(a)/100)
	fmt.Printf("%.2f\n", float64(a)/100)
}
输出
0x2
%!f(uint64=02)
2.5
2.50
所以不像其他语言，a/100.0的方式结果基本只会是a的类型。
测试
	var a uint64 = 250
	var b float64 = 750.0
	fmt.Printf("%#v\n", b/a) //发现不一样的类型，压根不能做运算，很搞笑
	fmt.Printf("%.2f\n", b/a)

-------------
go vendor机制在多个GOPATH的环境下偶尔不好用，已知是1.8出问题，其他版本不确定

-------------
ok:=true
if xxx {
	str,ok:=ttt.(string)
//这里的ok就和上面的ok没关系了
//但是如果没有if 这块呢？需要测试下
}