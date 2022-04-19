# sha256sum

Implementation on Golang.<br> 
The sha256sum command computes and checks a SHA256 encrypted message digest.

---
## Installation:
via Golang
```
    go mod download
    go build cmd/main.go
```
via Docker
```
    docker-compose build
    docker-compose up
```

---
## Usage:

You can use it with option like:
1. -f (path to file):
```
    go run cmd/main.go -f /base/.../dir/example.txt
```
2. -d (path to dir):
```
    go run cmd/main.go -d /base/.../dir/
```
3. -a (hash algorithm):
```
    go run cmd/main.go -d /base/.../dir/ -a md5
```
Now you can use sha256, sha512, md5 (default sha256)

4. -h (options docs):
```
    go run cmd/main.go -h
```

5. -check (compare new and old hashes):
```
    go run cmd/main.go -check /base/.../dir/
```
Also, you can use -a flag with that
```
    go run cmd/main.go -check /base/.../dir/ -a md5
```

---
## hashsum pkg:
view documentation with
```
    go doc hashsum
```
or
```
    godoc -http :6060 // you can use other ports
```
like HTML page