## Viki API Signer

### Installation

```sh
go get -u github.com/viki-org/gomods/signer
```

### Example

```go

signer := signer.New("viki_app_id", "secret")

// Sign a URL string
url := "https://api.viki.io/v4/clips.json"
signedURL := signer.Sign(url)
// "https://api.viki.io/v4/clips.json?app=viki_app_id&t=1436429448&sig=sha1hmacsignature"

// You can also pass in a path
path := "/v4/clips.json"
signedPath := signer.Sign(path)
// "/v4/clips.json?app=viki_app_id&t=1436429448&sig=sha1hmacsignature"

// Or obtain the signature directly
sig := signer.GetSignature(path)
// "sha1hmacsignature"

```

### Testing

```sh
go test
```
