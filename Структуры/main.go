package main

import (
	"fmt"
	"math"
)

type User struct {
	Name string

	Rating float64
}

func (u User) Greeting ()  {
	fmt.Println("Здарова, я", u.Name)
	fmt.Println("Мой рейтинг:", u.Rating)
}

func (u User) Goodbye()  {
	fmt.Println("Пока, с вами был ваш", u.Name)
	fmt.Println("Мой рейтинг бывал:", math.Float64bits(u.Rating))
}

// Это метод, здесь мы как бы расширяем структуру User
func (u *User) upRating(rating float64) {
	if u.Rating + rating <= 10 {
		u.Rating += rating
		fmt.Println("Я добавил рейтинг, теперь он:", u.Rating)
	} else {
		fmt.Println("Рейтинг не может быть больше 10")
	}
}

//То же самое, но с указателем на структуру User, из-за чего будет меняться исходная структура, а не копия
//func upRating(u *User ,rating float64) {
//	if u.Rating + rating <= 10 {
//		u.Rating += rating
//		fmt.Printf("Я добавил рейтинг, теперь он: %.2f\n", u.Rating)
//	} else {
//		fmt.Println("Рейтинг не может быть больше 10")
//	}
//}

func main () {
	user := User{
		Name: "Серега",
		Rating: 4.5,
	}

	ptr := &user

	user.Greeting()
	ptr.upRating(2.28)
	user.Goodbye()
}