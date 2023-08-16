# options

Realization of pattern **functional parameter**.

Example:
```go
package a

type A struct{
    foo string
    bar string
}

func New(opts ...Option[A]) (*A, error){
    a := A{
        foo: "defaultFoo",
        bar: "defaultBar",
    }
    if err := ApplyOptions(&a, opts...); err != nil{
        return nil, fmt.Errorf("applying opts: %w", err)
    }
    return &a, nil
}

func WithFoo(v string) Option[A]{
    return func(target *A) error{
        target.foo = v
    }
}

func WithBar(v string) Option[A]{
    return func(target *A) error{
        if v == ""{
            return errors.New("bar can bot be empty")
        }
        target.bar = v
    }
}
```

```go
package main

func main(){
    a1, err := a.New()
    fatalOnErr(err)
    a2, err := a.New(a.WithFoo("customFoo"))
    fatalOnErr(err)
    a3, err := a.New(a.WithFoo(""), a.WithBar("")) // returns non-nil err

    for idx, el := range []*A{a1, a2, a3}{
        fmt.Printf("%d: %+v\n", idx, el)
    }
}

func fatalOnErr(err error){
    if err != nil{
        log.Fatal(err)
    }
}

func FatalOnNil(err error){
    if err == nil{
        log.Fatal(err)
    }
}
```