package helper

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ .env fayl topilmadi yoki yuklanmadi")
	}
}

func ENV(key string) string {
	return os.Getenv(key)
}

var cyrillicAlphabet = []string{
	"а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н",
	"о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я",
	"А", "Б", "В", "Г", "Д", "Е", "Ё", "Ж", "З", "И", "Й", "К", "Л", "М", "Н",
	"О", "П", "Р", "С", "Т", "У", "Ф", "Х", "Ц", "Ч", "Ш", "Щ", "Ъ", "Ы", "Ь", "Э", "Ю", "Я",
}
var latinAlphabet = []string{
	"a", "b", "v", "g", "d", "e", "yo", "j", "z", "i", "y", "k", "l", "m", "n",
	"o", "p", "r", "s", "t", "u", "f", "h", "ts", "ch", "sh", "sch", "", "y", "", "e", "yu", "ya",
	"a", "b", "v", "g", "d", "e", "yo", "j", "z", "i", "y", "k", "l", "m", "n",
	"o", "p", "r", "s", "t", "u", "f", "h", "ts", "ch", "sh", "sch", "", "y", "", "e", "yu", "ya",
}

func Slug(data string) string {
	for i, cyr := range cyrillicAlphabet {
		data = strings.ReplaceAll(data, cyr, latinAlphabet[i])
	}
	data = strings.ToLower(data)
	reg := regexp.MustCompile(`[^\w\d\- ]`)
	data = reg.ReplaceAllString(data, "")
	data = strings.ReplaceAll(data, " ", "-")
	reg = regexp.MustCompile(`\-{2,}`)
	data = reg.ReplaceAllString(data, "-")
	return data
}
