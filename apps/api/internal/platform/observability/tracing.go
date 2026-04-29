package observability

func StartSpan(_ string) func() { return func() {} }
