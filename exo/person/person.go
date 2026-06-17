package person

import "strconv"

type Person struct {
	name string
	age  int
}

func NewPerson(name string, age int) Person {
	return Person{
		name: name,
		age:  age,
	}
}

func (p Person) String() string {
	return "My name is " + p.name + " and I am " + strconv.Itoa(p.age) + " years old."
}
