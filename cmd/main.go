package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sha256sum/internal/repository"
	"sha256sum/internal/service"
	"sha256sum/internal/utils"
	"strconv"
	"time"
)

var (
	path     string
	hashType string
	duration int
	signals  chan os.Signal
)

func init() {
	flag.StringVar(&path, "d", "", "/example/.../dir/ || you can check hashsum sum by dir path")
	flag.StringVar(&hashType, "a", "sha256", "available: md5, sha512 || default: sha256")
	flag.Parse()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	duration, err = strconv.Atoi(os.Getenv("DURATION_TIME"))
	if err != nil {
		duration = 15
	}

	signals = make(chan os.Signal, 1)
}

func main() {
	utils.CheckSignal(signals)

	log.Println("### 🛠 Connecting to DB")

	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("### 🎯 Connection successful")

	r := repository.NewRepository(db)

	k, container, err := utils.NewK8SConnection()
	if err != nil {
		log.Fatalln(err)
	}

	s := service.NewService(r, hashType, k, container)

	log.Println("### 🗄 Check DB data")

	_, err = s.DirectoryHash(path)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("### 🎱 Database was checked")

	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	go s.DirectoryCheck(ticker, path)
	<-signals

	ticker.Stop()
}
