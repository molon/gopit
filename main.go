package main

import (
	"log"
)

type user struct {
	name  string
	email string
}

// //go:noinline
// func test(u interface{}) {
// 	t := &u
// 	println(t)
// }

// func GbkToUtf8(s []byte) ([]byte, error) {
// 	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
// 	d, e := ioutil.ReadAll(reader)
// 	if e != nil {
// 		return nil, e
// 	}
// 	return d, nil
// }

// func haha(m []string) []string {
// 	return m
// }

type tmp string

type tm struct {
	A int
	B int
}

func (t *tm) Stop() {
	log.Println(t.A, " ", t.B)
}

func test() (interface{}, bool) {
	return nil, true
}

func checkInclusion(s1 string, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}

	// 找到两个字符串以s1长度的字符计数
	s1c2c := map[byte]int{}
	s2c2c := map[byte]int{}
	for index := 0; index < len(s1); index++ {
		s1c2c[s1[index]]++
		s2c2c[s2[index]]++
	}

	left := 0

	equal := func() bool {
		// log.Println("s1c2c", s1c2c)
		// log.Println("s2c2c", s2c2c)
		// log.Println("s2", left, s2[left:left+len(s1)])
		for k := range s1c2c {
			if s1c2c[k] != s2c2c[k] {
				return false
			}
		}
		return true
	}

	// 移动s2修正s2c2c直到s1c2c和s2c2c相等
	// 从s1长度的位置开始移动检测
	// 例如 abc badbac
	for idx := len(s1) - 1; idx < len(s2); idx++ {
		if equal() {
			return true
		}

		// 如果已经是最后一个，无需再移动
		if idx >= len(s2)-1 {
			break
		}

		// 检测失败，则往右移动一个位置，把s2当前位置的字符计数+1，之前窗口的首位置字符计数-1
		s2c2c[s2[idx+1]]++
		s2c2c[s2[left]]--
		left++
	}

	return false
}

func main() {
	log.Println(checkInclusion("ab", "eidbaooo"))
	// a := make(chan bool, 1)
	// b := make(chan bool, 1)

	// a <- true
	// go func() {
	// 	for index := 0; index < 10; index += 2 {
	// 		<-a
	// 		log.Println(index)
	// 		b <- true
	// 	}
	// }()

	// go func() {
	// 	for index := 1; index < 10; index += 2 {
	// 		<-b
	// 		log.Println(index)
	// 		a <- true
	// 	}
	// }()

	// time.Sleep(1 * time.Second)
	// runtime.GOMAXPROCS(1)

	// go func() {
	// 	for index := 0; index < 10; index += 2 {
	// 		log.Println(index)
	// 		// runtime.Gosched()
	// 	}

	// 	time.Sleep(3 * time.Second)
	// }()

	// go func() {
	// 	runtime.Gosched()
	// 	for index := 1; index < 10; index += 2 {
	// 		log.Println(index)
	// 		// runtime.Gosched()
	// 	}
	// }()
	// time.Sleep(5 * time.Second)
	// doneC := make(chan struct{}, 1)
	// <-doneC
	// s := []string{"a", "b", "c", "d"}
	// for i, v := range s {
	// 	log.Println(i, v)
	// 	if "b" == v {
	// 		s = append(s[:i], s[i+1:]...)
	// 	}
	// }

	// log.Println("Begin")
	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	<-time.After(1 * time.Second)
	// }()

	// go func() {
	// 	wg.Wait()
	// 	log.Println("Wait1")
	// }()

	// wg.Wait()
	// log.Println("Wait2")

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	<-time.After(1 * time.Second)
	// }()

	// wg.Wait()
	// log.Println("Wait3")

	// a4 := 1009
	// var a5 uint64 = 0xFFFFFFFFFFFFFFF0
	// var aa uint64 = uint64(a4+15) & a5
	// print(aa)
	// a := tm{
	// 	A: 1,
	// 	B: 2,
	// }
	// b := tm{
	// 	A: 1,
	// 	B: 2,
	// }
	// b.A = 3
	// log.Println(a == b)
	// log.Printf("%p,%p", &a, &b)
	// var callbacks []func()
	// tms := []tm{
	// 	tm{A: 1},
	// 	tm{A: 2},
	// }
	// tms2 := []tm{}
	// for idx := 0; idx < 2; idx++ {
	// 	callbacks = append(callbacks, func() {
	// 		log.Println(idx)
	// 	})
	// }

	// for _, tm := range tms {
	// 	tms2 = append(tms2, tm)

	// 	callbacks = append(callbacks, func() {
	// 		log.Println(tm)
	// 	})
	// }

	// for _, callback := range callbacks {
	// 	callback()
	// }

	// jsn, _ := json.Marshal(tms2)
	// log.Println(string(jsn))

	// fmt.Println(1 << 0)
	// fmt.Println(1 << 1)
	// fmt.Println(1 << 2)
	// fmt.Println(1 << 3)
	// a, b := test()
	// log.Printf("%v:%v", a.(*tm), b)

	// as := map[int]int{}
	// for index := 0; index < 20; index++ {
	// 	as[index] = index
	// }

	// cnt := 0
	// for k, _ := range as {
	// 	log.Println(k)
	// 	delete(as, k)
	// 	cnt++
	// }

	// log.Println("count: ", cnt)

	// c := make(chan int)

	// for {
	// 	select {
	// 	case n := <-c:
	// 		if n == 10 {
	// 			break
	// 		}
	// 	}
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	// defer cancel()

	// go func() {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Printf("done err: %v", ctx.Err())
	// 	}
	// }()

	// go func() {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Printf("done err: %v", ctx.Err())
	// 	}
	// }()

	// time.Sleep(time.Second)

	// cancel()

	// time.Sleep(time.Second)

	// doneC := make(chan int, 3)
	// doneC <- 1
	// doneC <- 2
	// doneC <- 3
	// close(doneC)
	// doneC <- 4
	// fmt.Println(<-doneC)
	// fmt.Println(<-doneC)
	// fmt.Println(<-doneC)
	// fmt.Println(<-doneC)
	// fmt.Println(<-doneC)
	// var a = "1"
	// var b = "1"
	// var c tmp = "1"

	// log.Printf("%#v %#v %#v ", a, b, c)
	// if a == b {
	// 	log.Println("a==b")
	// }
	// if tmp(a) == string(c) {
	// 	log.Println("a==c")
	// }
	// m1 := []string{"1", "2"}
	// m2 := haha(m1)
	// m2[1] = "3"
	// fmt.Println(m1)

	// rand.Seed(time.Now().UnixNano())
	// index := rand.Intn(10)
	// fmt.Println(index)
	// file, err := os.Open("template.csv")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// defer file.Close()
	// reader := csv.NewReader(file)
	// for {
	// 	cols, err := reader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		fmt.Println("Error:", err)
	// 		return
	// 	}

	// 	fmt.Println(cols)

	// 	utf8Cols := []string{}
	// 	for _, col := range cols {
	// 		str, err := GbkToUtf8([]byte(col))
	// 		if err != nil {
	// 			fmt.Println("Error:", err)
	// 			return
	// 		}
	// 		utf8Cols = append(utf8Cols, string(str))
	// 	}
	// 	fmt.Println(utf8Cols)
	// }

	// rand.Seed(time.Now().UnixNano())
	// fmt.Print(rand.Intn(100), ",")

	// u1 := createUserV1()
	// u2 := createUserV2()
	// println("u1", &u1, "u2", &u2)

	// u := user{
	// 	name: "molon",
	// }
	// test(u)

	// u := 1
	// fmt.Println(u)
	// fmt.Println(&u)
	// print(u)
}

// //go:noinline
// func createUserV1() user {
// 	u := user{
// 		name:  "Bill",
// 		email: "bill@ardanlabs.com",
// 	}
// 	println("V1", &u)
// 	return u
// }

// //go:noinline
// func createUserV2() *user {
// 	u := user{
// 		name:  "Bill",
// 		email: "bill@ardanlabs.com",
// 	}
// 	println("V2", &u)
// 	return &u
// }
