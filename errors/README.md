# errors
This `errors` package is designed based on [Failure is your Domain](https://middlemost.com/failure-is-your-domain/) blog:

    The tricky part about errors is that they need to be different things to different consumers of them. 
    In any given system, we have at least 3 consumer roles — the application, the end user, & the operator.

## Features 
- Allows specifying `Logical Operation` that caused the failure. helps `the operator`
- Allows specifying `Machine-readable` error code. helps `the application`
- Allows specifying `Human-readable` message.  helps `the end user`
- Allows specifying nested error.
- Support gRPC styles errors with [Details](https://jbrandhorst.com/post/grpc-errors/)

**ErrorCode** helps to categorize errors into:
1. System Errors - Only recoverable after fixing system failures. e.g., disk fill, database down., certs expaired.  
2. Temporary Errors - Recoverable immediately after retry with exponential backoff
3. Data Errors - Input validation errors which need to be reprocessed after fixing the data issues  

Each _feature_ can be added to the previous `wrapped` or `leaf` error, using [available wrapper constructors](https://github.com/cockroachdb/errors#Available-wrapper-constructors) <br/>
While using wrapped errors, you can reveal each **character/trait** by unwrapping like an **Onion** 🧅. <br/>
using `Unwrap()` / `As()` / `Is()` or helpers functions `GetAllDetails`, `FlattenDetails`, `HasAssertionFailure`, `GetCategory`, `GetCode`, `GetOperation` etc.
 
## Interface 

```go
type ErrorOperation interface {
	error
	Operation() string
}

type ErrorCoder interface {
	error
	Code() codes.Code
	Category() categories.Category
}
```

## TODO
- make this package generic (use `int` instead fo `codes.Code`?) and move it to `toolkit`
- `codes.Code` should be in consuming application
 
## Reference
- [Failure is your Domain](https://middlemost.com/failure-is-your-domain/)
- [Creating domain specific error helpers in Go with errors.As](https://blog.carlmjohnson.net/post/2020/working-with-errors-as/)

- [advanced gRPC Error Usage](https://jbrandhorst.com/post/grpc-errors/)
- [cockroachdb errors](https://github.com/cockroachdb/errors)
- [atlas-app-toolkit errors](https://github.com/infobloxopen/atlas-app-toolkit/tree/master/errors)
- [go-multierror](https://github.com/hashicorp/go-multierror)
