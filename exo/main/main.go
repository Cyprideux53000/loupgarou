package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func makeCounter(init int) func() int {
	offset := init + 2
	return func() int {
		offset++
		return offset
	}
}

func swap(a *int, b *int) {
	var tmp int
	tmp = *a
	*a = *b
	*b = tmp
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if u.Name == "" || u.Email == "" || u.Password == "" {
		http.Error(w, "name, email and password are required", http.StatusBadRequest)
		return
	}

	sendJSON(w, map[string]string{
		"message": "Account created successfully",
		"name":    u.Name,
		"email":   u.Email,
	})
}

func sendJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	// fmt.Println("Hello, World!")

	// p := person.NewPerson("Alice", 30)
	// e := employe.NewEmployee(p, 5000)
	// fmt.Println(p)
	// fmt.Println(e)

	// x := 5
	// y := 5
	// if x > 0 {
	// 	y := 10
	// 	fmt.Println("x is positive:", x, "y is:", y)
	// }
	// fmt.Println("x is:", x, "y is:", y)

	// add := func(a, b int) int {
	// 	return a + b
	// }
	// result := add(3, 5)
	// fmt.Println(result)
	// fmt.Println(makeCounter(0))
	// init := 1
	// fmt.Println(makeCounter(init)())
	// fmt.Println(init)

	// a, b := 1, 2
	// fmt.Printf("Before swap, value of a=%d and b==%d \n", a, b)
	// swap(&a, &b)
	// fmt.Printf("After swap, value of a=%d and b=%d \n", a, b)

	// http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Hello World!")
	// })
	// http.HandleFunc("/account", createAccount)
	// log.Println("Server listening on :8080")
	// http.ListenAndServe(":8080", nil)

	ages := map[string]int{
		"Alice": 25,
	}
	fmt.Println(ages["Alice"])
	if age, ok := ages["Alice"]; ok {
		fmt.Println(age)
	}
	fmt.Println(ages["John"])
	if age, ok := ages["John"]; ok {
		fmt.Println(age)
	} else {
		fmt.Println("Not found")
	}
}
