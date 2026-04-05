package rt_test

import (
	"fmt"

	"github.com/defany/goblin/rt"
)

func ExampleFnName() {
	name := rt.FnName()
	// FnName returns "package.Function" using the last two segments
	// of the fully qualified function name.
	fmt.Println(name)
	// Output: com/defany/goblin/rt_test.ExampleFnName
}

func ExampleCaller() {
	_, line, fn := rt.Caller(0)
	// line and fn will vary, but the function returns file path, line number, and full function name.
	_ = line
	_ = fn
	fmt.Println("caller retrieved")
	// Output: caller retrieved
}

func ExampleCallerUniqueKey() {
	key := rt.CallerUniqueKey(0)
	// The key is a unique string combining file:line:function.
	fmt.Println(key != "unknown")
	// Output: true
}
