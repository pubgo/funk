package errors

// RootCause returns the deepest error in the unwrap chain.
func RootCause(err error) error {
	if err == nil {
		return nil
	}

	for {
		next := Unwrap(err)
		if next == nil {
			return err
		}
		err = next
	}
}

// Walk traverses the error chain from outer to inner.
// Returning false from fn stops traversal.
func Walk(err error, fn func(error) bool) {
	for err != nil {
		if !fn(err) {
			return
		}
		err = Unwrap(err)
	}
}
