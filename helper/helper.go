package helper

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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

func FormatDateOnly(date time.Time) string {
	return date.Format("02-01-2006")
}

func FormatTime(date time.Time) string {
	return date.Format("15:04:05")
}

func FormatDate(date time.Time) string {
	return date.Format("02-01-2006 15:04:05")
}

var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"seq": func(start, end int, current int) []int {
		var pages []int
		window := 2

		if start < 1 {
			start = 1
		}

		if current-window > start {
			pages = append(pages, 1)
			if current-window > 2 {
				pages = append(pages, -1)
			}
		}

		for i := current - window; i <= current+window; i++ {
			if i >= 1 && i <= end {
				pages = append(pages, i)
			}
		}

		if current+window < end {
			if current+window < end-1 {
				pages = append(pages, -1)
			}
			pages = append(pages, end)
		}

		return pages
	},
}

func View(ctx echo.Context, layoutName, viewName string, data interface{}) error {
	tmpl, err := template.New("layout").
		Funcs(templateFuncs).
		ParseFiles("views/"+layoutName, "views/"+viewName)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Template parsing error: "+err.Error())
	}

	return tmpl.ExecuteTemplate(ctx.Response().Writer, "layout", data)
}
