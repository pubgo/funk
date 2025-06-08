package anyhow

// === Core State Interfaces ===

type Catchable interface {
	Catch(err *error) bool
	CatchErr(err *Error) bool
}

// Checkable defines types that can be checked for Ok/Error state
type Checkable interface {
	IsOk() bool
	IsError() bool
	String() string
}

// Unwrappable defines types that can provide their contained value
type Unwrappable[T any] interface {
	Unwrap() T
	UnwrapOr(defaultValue T) T
}

// ErrorAccessible defines types that can provide error information
type ErrorAccessible interface {
	Err() error
}

// === Functional Programming Interfaces ===

// Mappable defines types that support value transformation
type Mappable[T any] interface {
	Map(fn func(T) T) Mappable[T]
}

// TypeMappable defines types that support type transformation
type TypeMappable[T any] interface {
	// Note: Go generics limitation - we can't define this perfectly in interface
	// But we can document the expected signature: MapTo[U any](fn func(T) U) TypeMappable[U]
}

// Flatmappable defines types that support monadic bind operations
type Flatmappable[T any] interface {
	FlatMap(fn func(T) Flatmappable[T]) Flatmappable[T]
}

// Filterable defines types that support conditional filtering
type Filterable[T any] interface {
	Filter(predicate func(T) bool, errorMsg string) Filterable[T]
}

// Inspectable defines types that support side-effect operations
type Inspectable[T any] interface {
	Inspect(fn func(T)) Inspectable[T]
}

// ErrorInspectable defines types that support error inspection
type ErrorInspectable interface {
	InspectErr(fn func(error)) ErrorInspectable
}

// === Recovery and Error Handling Interfaces ===

// Recoverable defines types that support error recovery
type Recoverable[T any] interface {
	OrElse(fn func(error) Recoverable[T]) Recoverable[T]
}

// ErrorMappable defines types that support error transformation
type ErrorMappable interface {
	MapErr(fn func(error) error) ErrorMappable
}

// Panicable defines types that can panic on error
type Panicable[T any] interface {
	Must() T
	Expect(msg string) T
}

// === Combination Interfaces ===

// Monad combines the essential monadic operations
type Monad[T any] interface {
	Mappable[T]
	Flatmappable[T]
	Inspectable[T]
}

// ErrorHandlingMonad extends Monad with error handling capabilities
type ErrorHandlingMonad[T any] interface {
	Monad[T]
	ErrorAccessible
	ErrorMappable
	ErrorInspectable
	Recoverable[T]
	Must() T
	Expect(msg string) T
}

// Container represents a generic container with state checking
type Container[T any] interface {
	Checkable
	Unwrappable[T]
	Inspectable[T]
}

// === Conversion Interfaces ===

// ResultConvertible defines types that can be converted to Result
type ResultConvertible[T any] interface {
	ToResult(defaultValue T, err error) Result[T]
}

// Note: OptionConvertible removed as Option[T] is not idiomatic in Go

// ErrorConvertible defines types that can be converted to Error
type ErrorConvertible interface {
	ToError() Error
}

// === Utility Interfaces ===

// Combinable defines types that can be combined with logical operations
type Combinable[T any] interface {
	And(other Combinable[T]) Combinable[T]
	Or(other Combinable[T]) Combinable[T]
}

// Transposable defines types that support transpose operations
type Transposable[T, U any] interface {
	// Transpose operation for nested types
}

// === Specialized Interfaces ===

// ValueContainer represents containers that primarily hold values
type ValueContainer[T any] interface {
	Container[T]
	Mappable[T]
	Filterable[T]
}

// ErrorContainer represents containers that primarily handle errors
type ErrorContainer interface {
	Checkable
	ErrorAccessible
	ErrorMappable
	ErrorInspectable
	Must()
	Expect(msg string)
}

// Note: OptionalContainer removed as Option[T] is not idiomatic in Go

// === Functional Utilities ===

// Functor represents the functor pattern
type Functor[T any] interface {
	// Map: (T -> U) -> F[T] -> F[U]
	Map(fn func(T) T) Functor[T]
}

// Applicative extends Functor with application
type Applicative[T any] interface {
	Functor[T]
	// Apply would be here, but Go's type system makes this complex
}

// Alternative represents choice between alternatives
type Alternative[T any] interface {
	Or(other Alternative[T]) Alternative[T]
	And(other Alternative[T]) Alternative[T]
}

// === Collection Interfaces ===

// Collectible defines types that can be collected from slices
type Collectible[T, U any] interface {
	// Used by collection functions
}

// Traversable defines types that can be traversed
type Traversable[T any] interface {
	// Used for traverse operations
}

// === Meta Interfaces ===

// TypeInfo provides type information
type TypeInfo interface {
	TypeName() string
	IsEmpty() bool
}

// Serializable defines types that can be serialized
type Serializable interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// === Constraint Interfaces ===

// Ordered represents types that have ordering
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// Comparable represents types that can be compared
type Comparable interface {
	comparable
}
