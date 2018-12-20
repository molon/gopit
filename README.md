# gopit
golang坑的一些记录和测试

#### 只能文字记录的一些
- `vendor`文件夹在`go1.8`版本的多`GOPATH`环境下可能会被忽略，最终被使用的还是外部库。
- 如果使用`vendor`文件夹，请尽量保证不会引用文件夹以外的库，也不要让他们来引用，否则容易产生各种各样的问题。最明显的例子就是`若双方引用了同名的其他pkg，因为上述使用不当最终编译器会认为类型错误，而且也遇到过不提示编译成功，这更可怕，查错很难(那次导致了两个全局变量被生成，pkg1写的是这个，pkg2读的却是那个)。`
- `vscode`的`debug`配置里要加上`env:{}`后，它才能找到`GOPATH`
- `gorm`的`Count`方法会转成`COUNT(*)`，所以如果`Join`了两张表有相同名称字段的话，`sql`会提示命名模糊的错误。而手动指定`Select("COUNT(1) AS cnt")`语句的话，`COUNT(1) AS`部分会被忽略

#### main_test.go 概览
- `TestDivision` 整型浮点数运算
- `TestStmt` `:=`在不同作用域下的行为
- `TestRecover` `recover`的作用域
- `TestMap` `map`里`nil`元素相关
- `TestFmt` `fmt`相关
- `TestDefer` 测试`defer`的作用域和执行顺序啥的
- `TestSet` 测试赋值语句相关
- `TestSlice` 测试越界问题
- `TestIsSliceOrArrayWithKind` 通过反射判断是否是某类型的数组或者slice
- `TestTime` `time`时区相关
- `TestTransferNilInterface` 使用接口类型作为参数和返回值传递nil的判断

#### 其他未整理
```
大坑集锦

下面出现的S都表示为一个结构体类型

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
//所以切记，对e进行&取地址然后存储其地址毫无意义，因为其内容在下次循环必定会被改变，必须存储的话，请自行copy一份
//例如ee:=e都可
// 如果a是指针元素的数组，则一般无忧，因为一般不会对指针元素再去取地址
}

---
https://stackoverflow.com/a/28143457
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
	}

	combines := [][]int{}
	collectIntCombines(&groups, 0, []int{}, &combines)

	t.Log(combines)
	t.Log("一共", len(combines), "个结果")
}

-------------
```
