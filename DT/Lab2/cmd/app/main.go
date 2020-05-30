package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/MishaNiki/Labs/DT/Lab2/internal/app/spamdetection"
	"github.com/MishaNiki/Labs/DT/Lab2/internal/app/storage"
)

var (
	loadStatis bool //Download statistics
	pathSpam   string
	pathOK     string
	message    string
)

func init() {
	flag.BoolVar(&loadStatis, "loadSt", false, "Download statistics")
	flag.StringVar(&pathSpam, "spam", "", "path to Spam")
	flag.StringVar(&pathOK, "ok", "", "path to messageOK")
	flag.StringVar(&message, "message", "", "spam check message")
}

func main() {
	flag.Parse()

	// Open DB
	stor := storage.New()
	if err := stor.Open(&storage.Config{
		User:     "postgres",
		Password: "postgres",
		DBName:   "Labs",
		SSLMode:  "disable",
	}); err != nil {
		log.Fatal(err)
		return
	}
	defer stor.Close()

	if loadStatis {
		if pathSpam == "" && pathOK == "" {
			log.Println("Empty paths to statistics files")
			return
		}

		dict := spamdetection.NewDictionary()

		if pathSpam != "" {
			spam := SplitMessage(loadFile(pathSpam))
			for _, v := range spam {
				dict.AddSpam(v)
			}
		}

		if pathOK != "" {
			msgok := SplitMessage(loadFile(pathOK))
			for _, v := range msgok {
				dict.AddOK(v)
			}
		}

		stor.LoadStatistics(dict)
		log.Println("[OK] Statistics uploaded to the database")

		return
	}

	if message == "" {
		log.Println("empty -message")
		return
	}

	spam, err := checkSpam(message, stor)

	if err != nil {
		log.Fatal(err)
		return
	}

	if spam {
		fmt.Println("[SPAM]\t", message)
	} else {
		fmt.Println("[NOT SPAM]\t", message)
	}

}

// SplitMessage ...
func SplitMessage(str string) []string {
	str = strings.Trim(str, `.,!?"'«»–\n: `)
	str = strings.ToLower(str)
	re := regexp.MustCompile(`[.,!?"'«»–\n:/ ]+`).Split(str, -1)
	return re
}

func loadFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func checkSpam(msg string, stor *storage.Storage) (bool, error) {

	arrmsg := SplitMessage(message)
	dict, err := stor.DownloadStatistics()
	if err != nil {
		return false, err
	}

	var sp, gp float64 = 1, 1

	for _, v := range arrmsg {
		sp *= dict.GetProbabilitySpam(v)
		gp *= dict.GetProbabilityOK(v)
	}

	z := sp / (sp + gp)

	fmt.Println("z = ", z)

	if z > 0.6 {
		return true, nil
	} else {
		return false, nil
	}
}

/*
	ТЗ

	Загрузить базу предложений спама и нормальных сообщений в одинаковыз пропорциях
	Пословно разделить каждое предложение добавить его в map из типоа слов

	Потом по формулк проверять полученный результат

	«Привет, дружок! Выслал тебе деньжат. Купи мне подарок»
*/
