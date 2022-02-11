# pertr

This is [http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper) transport
wrapper that handles network errors while download request body. When connection
problem occurs it repeats request to server with set `Range` header and continue
download response body.

**Not production ready yet.**


## Usage

```go
cxt := context.WithTimeout(context.Background(), 1 * time.Second)
req, err := http.NewRequestWithContext(ctx, "GET", "https://sbercloud.ru/", http.NoBody)
if err != nil {
    panic(err)
}

cli := http.Client{Transport: pertr.Wrapper{}}
rsp, err := cli.Do(req)
```

## Warning!

1. This wrapper will repeat request to server, even resend POST request body. Please
use it with caution.

2. Go docs [says](https://pkg.go.dev/net/http#RoundTripper) that
[http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper) should not attempt
to interpret the response. However, this library does this.
