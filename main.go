package main

import "fmt"

func main(){
	name := "Васечка"
	fmt.Println("До изменения:", name)

	ptr := &name
	changeName(ptr)
	fmt.Println("После изменения:", name)
}


func changeName (s *string){
	*s = "Some Name"
}