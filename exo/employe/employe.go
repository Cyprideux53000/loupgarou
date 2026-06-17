package employe

import (
	"loupgarou/person"
	"strconv"
)

type Employee struct {
	salary int
	person.Person
}

func NewEmployee(p person.Person, salary int) Employee {
	return Employee{
		salary: salary,
		Person: p,
	}
}

func (e Employee) String() string {
	return e.Person.String() + " My salary is " + strconv.Itoa(e.salary) + "."
}
