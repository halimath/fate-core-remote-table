package usecase

import "context"

// Func is a generic function type that is used to define use case functions that
// provide a return value.
type Func[I, O any] func(context.Context, I) (O, error)

// Proc is a generic function type that is used to define use case functions
// that produce no result (apart from error) - so they are considered used case
// "procedures"
type Proc[I any] func(context.Context, I) error
