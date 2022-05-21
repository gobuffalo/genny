package safe

/* This is a copy of `gobuffalo/events/internal/safe/safe.go`, which was
   originally `markbates/safe/safe.go`. If you found any bug on this file,
   please fix the copies too.
*/

import (
	"errors"
	"fmt"
)

// Run the function safely knowing that if it panics
// the panic will be caught and returned as an error
func Run(fn func()) (err error) {
	return RunE(func() error {
		fn()
		return nil
	})
}

// Run the function safely knowing that if it panics
// the panic will be caught and returned as an error
func RunE(fn func() error) (err error) {
	defer func() {
		if err != nil {
			return
		}
		if ex := recover(); ex != nil {
			if e, ok := ex.(error); ok {
				err = e
				return
			}
			err = errors.New(fmt.Sprint(ex))
		}
	}()
	return fn()
}
