package safe

/* This is a copy of `gobuffalo/events/internal/safe/safe.go`, which was
   originally `markbates/safe/safe.go`. If you found any bug on this file,
   please fix the copies too.

   This safeguard functions will be deprecated soon.
   see https://github.com/gobuffalo/genny/pull/47 for more details.
*/

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// Run the function safely knowing that if it panics
// the panic will be caught and returned as an error
func Run(fn func()) (err error) {
	return RunE(func() error {
		fn()
		return nil
	})
}

var warningOnce sync.Once

// see discussion on PR #47
const warning = "PANIC RECOVERED! currently, genny recovers panic on the runner functions but the behavior will be dropped. please prepare your own recovery."

// Run the function safely knowing that if it panics
// the panic will be caught and returned as an error
func RunE(fn func() error) (err error) {
	defer func() {
		if err != nil {
			return
		}
		if ex := recover(); ex != nil {
			warningOnce.Do(func() {
				fmt.Fprintln(os.Stderr, warning)
			})
			if e, ok := ex.(error); ok {
				err = e
				return
			}
			err = errors.New(fmt.Sprint(ex))
		}
	}()
	return fn()
}
