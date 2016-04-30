events
======

Crazy Fast Event Dispatcher for golang.

1. NO REFLECTION
2. Async
3. Error Handler
4. Additional Event Store

> Events are somethings that happened in the past, 
> are executed async, 
> and you should call it with past name and do not rely on what the listener returns... 

Usage example:

```
helloListener := func(a interface{}) error {
    fmt.Print("Hello ", in)
    return nil
}
errorListener := func(e ListenerError){
    wasError = true
}
 
e := events.New()
e.On("hailed", helloListener)
e.OnError(errorListener) // If one listeners return error != nil

e.Raise("hailed", "World")

e.Wait() // optional
```

see `readme_test.go`


## Event repository with assertion on events (optional)

Is so useful in our tests to add an event store in order to make assertion on what
happened in your system

```
e := events.New()
e.AddInMemoryEventRepo() // or if you have one inject it see: AddEventRepo

e.On("hailed", helloListener)
e.Contains("hailed")     // true
``` 

## Test it

install deps `go get ./...`

and `go test ./... -v`
