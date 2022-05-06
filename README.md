# sha256sum

Implementation on Golang.<br> 
The sha256sum command computes and checks a SHA256 encrypted message digest.

---
## :hammer: Installation:

1. Like Docker container: <br>
configure .env file with your dir path
```
    DIR_PATH=YOUR_PATH_HERE
```
2. after build docker-compose
```
    docker-compose build
    docker-compose up -d postgres
```
3. Also you can use application without docker
```
    go mod download
    go build cmd/main.go
        or
    ./build.sh {your_os} // win32, win64, linux32, linux64, osx32, osx64
```
---
## :floppy_disk: Usage:

You can use it with option like:
1. **`-f`** (path to file):
```
    docker-compose run hasher -f=/local/your_path/example.txt
    go run cmd/main.go -f=/your_path/example.txt
```
2. **`-d`** (path to dir):
```
    docker-compose run hasher -d=/local/your_path/
    go run cmd/main.go -d /your_path/
```
3. **`-a`**(hash algorithm): <br>
   Now you can use sha256, sha512, md5 (default sha256)
```
    docker-compose run hasher -d=/local/your_path/ -a=md5
    go run cmd/main.go -d=/your_path/ -a=md5
```

4. **`-h`** (options docs):
```
    docker-compose run hasher -h
    go run cmd/main.go -h
```

5. **`-check`** (compare new and old hashes): <br>
    Also, you can use **`-a`** flag with that
```
    docker-compose run hasher -check=/local/your_path/ -a=md5
    go run cmd/main.go -check=/your_path/ -a=md5
```

6. **`-deleted`** (check deleted files in your dir):
```
    docker-compose run hasher -deleted=/local/your_path/
    go run cmd/main.go -deleted=/your_path/
```
---
## :mag_right: hashsum pkg:
View documentation with
```
    go doc hashsum
        or
    godoc -http :6060 // you can use other ports
```