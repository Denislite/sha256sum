# sha256sum

Implementation on Golang.<br> 
The sha256sum command computes and checks a SHA256 encrypted message digest.

---
## Installation:

coming soon

---
## Usage:

You can use it with option like:
1. -f (path to file):
```
    go run cmd/feature_{last_num}/main.go -f /base/.../dir/example.txt
```
2. -d (path to dir):
```
    go run cmd/feature_{last_num}/main.go -d /base/.../dir/
```
3. -a (hash algorithm):
```
    go run cmd/feature_{last_num}/main.go -d /base/.../dir/ -a md5
```
Now you can use sha256, sha512, md5 (default sha256)

4. -h (options docs):
```
    go run cmd/feature_{last_num}/main.go -h
```
