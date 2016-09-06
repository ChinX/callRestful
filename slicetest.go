package main

import "log"

func main() {
	// 创建一个长度为0，容量为1的切片
	// 底层隐式的创建了一个数组 arr=[1]int{0}
	// a引用了arr的地址，并引用了arr从位置0的0个长度
	a := make([]int, 0, 1)

	// 将a的引用传递给b，b引用了arr从位置0的0个长度
	b := a

	//打印a、b的地址（相同）
	log.Printf("a 地址 ：%p", a)
	log.Printf("b 地址 ：%p", b)
	//打印a、b的值（相同）
	log.Println("a is", a)
	log.Println("b is", b)

	// 切片a增加一个元素
	// 底层隐式执行了判断：arr长度 >= a增加元素后的长度
	// a引用了arr从位置0的1个长度，并将arr该长度位置的值赋值为1
	// b引用了arr从位置0的0个长度，所以b的值不变
	a = append(a,1)

	//打印a、b的地址（相同）
	log.Printf("a 地址 ：%p", a)
	log.Printf("b 地址 ：%p", b)
	//打印a、b的值（不同）
	log.Println("a is", a)
	log.Println("b is", b)

	// 切片b增加一个元素
	// 底层隐式执行了判断：arr长度 >= b增加元素后的长度
	// b引用了arr从位置0的1个长度，并将arr该长度位置的值赋值为2
	// a引用了arr从位置0的1个长度，所以a的值改变了
	b = append(b,2)

	//打印a、b的地址（相同）
	log.Printf("a 地址 ：%p", a)
	log.Printf("b 地址 ：%p", b)
	//打印a、b的值（相同）
	log.Println("a is", a)
	log.Println("b is", b)

	// 将a的引用传递给c，c引用了arr从位置0的1个长度
	// 增加c仅供后面的对比参考
	c := a

	// 切片a再增加一个元素
	// 底层隐式执行了判断：arr长度 < a增加元素后的长度
	// 底层隐式创建了一个长度为2*arr长度的数组并copy了arr的值，arr1 := [2]int{2,0}
	// a引用变为arr1从位置0的2个长度，并将arr该长度位置的值赋值为3
	// b引用了arr从位置0的1个长度，所以b的值不变
	// c引用了arr从位置0的1个长度，所以c的值不变
	a = append(a,3)

	//打印a、b、c的地址（b、c相同，a不同）
	log.Printf("a 地址 ：%p", a)
	log.Printf("b 地址 ：%p", b)
	log.Printf("c 地址 ：%p", c)
	//打印a、b、c的值（b、c相同，a不同）
	log.Println("a is", a)
	log.Println("b is", b)
	log.Println("c is", c)

	// 切片b再增加一个元素
	// 底层隐式执行了判断：arr长度 < b增加元素后的长度
	// 底层隐式创建了一个长度为2*arr长度的数组并copy了arr的值，arr2 := [2]int{2,0}
	// b引用变为arr2从位置0的2个长度，并将arr该长度位置的值赋值为4
	// a引用了arr1从位置0的2个长度，所以a的值与上次相同
	// c引用了arr从位置0的1个长度，所以c的值不变
	b = append(b,4)

	//打印a、b、c的地址（a、b、c不同，c地址还是原来的地址）
	log.Printf("a 地址 ：%p", a)
	log.Printf("b 地址 ：%p", b)
	log.Printf("c 地址 ：%p", c)
	//打印a、b、c的值（a、b、c不同，c值还是原来的值）
	log.Println("a is", a)
	log.Println("b is", b)
	log.Println("c is", c)
}
