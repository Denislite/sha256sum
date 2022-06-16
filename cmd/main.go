package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sha256sum/internal/repository"
	"sha256sum/internal/service"
	"sha256sum/internal/utils"
	"time"
)

var (
	path     string
	hashType string
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

	signals = make(chan os.Signal, 1)
}

func main() {
	utils.CheckSignal(signals)

	log.Println("|| 🛠 Connecting to DB      ||")

	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("|| 🎯 Connection successful ||")

	r := repository.NewRepository(db)

	s := service.NewService(r, hashType)

	log.Println("|| 🗄 Check DB data         ||")

	_, err = s.DirectoryHash(path)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("|| 🎱 Database was checked  ||")
	log.Printf("||==========================||")

	ticker := time.NewTicker(15 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				result, err := s.CompareHash(path)
				if err != nil {
					log.Fatalln(err)
				}
				if result != nil {
					log.Println("|| Changed files:")
					for _, hash := range result {
						log.Printf("|| %s %s %s \n",
							hash.FileName, hash.OldHash, hash.NewHash)
					}
					log.Fatalln("|| ❌  Files was changed, restarting pod...")
				}
				log.Println("|| ✅  Directory was checked, all right")
			}
		}
	}()
	<-signals

	ticker.Stop()
}
