package errfmt

import (
	"fmt"
	"strings"

	"github.com/defany/goblin/rt"
)

func WithSource(err error, comments ...string) error {
	if err == nil {
		return nil
	}

	loc := rt.CallerShortLocation(1)

	if len(comments) > 0 {
		return fmt.Errorf("%s [%s] -> %w", loc, strings.Join(comments, " "), err)
	}

	return fmt.Errorf("%s -> %w", loc, err)
}
