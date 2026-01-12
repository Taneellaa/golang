package main

import "fmt"

func main(){
	//name := "Васечка"
	//fmt.Println("До изменения:", name)

	//ptr := &name
	//changeName(ptr)
	//fmt.Println("После изменения:", name)

	//Пример nil указателя:
	number := 15
	fmt.Println("number:", number)

	var ptr *int = &number
		fmt.Println("ptr:", ptr)


	//if ptr != nil {
	//	fmt.Println("Разыменование:", *ptr)
	//} else {
	//	fmt.Println("Получен nil-указатель")
	//}
}

