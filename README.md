# pertr

This is [http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper) transport
wrapper that handles network errors while download request body. When connection
problem occurs it repeats request to server with set
[Range](https://datatracker.ietf.org/doc/html/rfc7233#section-3.1) header and
continue download response body.

**Not production ready yet.**


## Usage

Just put new wrapper instance to client's transport:

```go
wrapper := pertr.New()

cli := http.Client{Transport: wrapper}
rsp, err := cli.Do(req)
```

You can use custom context to manage retry timeout. If context will expire
download will continue until next error occured.

```go
cxt := context.WithTimeout(context.Background(), 15 * time.Minute)
wrapper := pertr.New(pertr.WithContext(ctx))
```

If context specified wrapper cannot be neither reused or used in parallel operations.
To reuse wrapper just specify new context with:

```go
wrapper.SetContext(newctx)
```

Please do not call `SetContext()` until download operation finished.

Wrapper can use custom transport ([http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper)) for outgoing calls:

```go
mytr := http.Transport{}
wrapper := pertr.New(pertr.WithTransport(&mytr))
```


## Errors

`Body.Read()` will return first error from original first request.


## Notes

1. This wrapper will repeat request to server. Please use it with caution.

2. Wrapper will try to repeat request if server supports resume downloading.
It will detects if server responds with `http.StatusPartialContent` in first
retry request.

3. Go docs [says](https://pkg.go.dev/net/http#RoundTripper) that
[http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper) should not attempt
to interpret the response. However, this library does this.
