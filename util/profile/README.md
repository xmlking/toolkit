# Profile

## Usage

```go
func BigIntFactorial(x big.Int) *big.Int {
    // Arguments to a defer statement is immediately evaluated and stored.
    // The deferred function receives the pre-evaluated values when its invoked.
    defer profile.Duration(time.Now(), "IntFactorial")

    y := big.NewInt(1)
    for one := big.NewInt(1); x.Sign() > 0; x.Sub(x, one) {
        y.Mul(y, x)
    }

    return x.Set(y)
}
```
