package cond_test

import (
	"fmt"

	"github.com/defany/goblin/cond"
)

func ExampleTernary() {
	result := cond.Ternary(true, "yes", "no")
	fmt.Println(result)
	// Output: yes
}

func ExampleTernary_false() {
	result := cond.Ternary(false, "yes", "no")
	fmt.Println(result)
	// Output: no
}

func ExampleTernary_int() {
	age := 20
	label := cond.Ternary(age >= 18, "adult", "minor")
	fmt.Println(label)
	// Output: adult
}
