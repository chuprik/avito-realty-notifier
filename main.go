package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ItemType struct {
	ID        string
	Title     string
	Link      string
	Address   string
	Price     int
	Comission string
}

var (
	fileName = "avito-notifier.json"
	link     = "https://www.avito.ru"

	savedData  []ItemType
	parsedData []ItemType

	searchURL  *string
	webhookURL *string
)

func init() {
	searchURL = flag.String("s", "", "URL with results")
	webhookURL = flag.String("w", "", "Slack webhookURL URL")
}

func main() {
	flag.Parse()

	var err error

	savedData, err = loadData(fileName)
	errorCheck(err, "Error: Trying to load data")

	parsedData, err = getParsedItems(*searchURL)
	errorCheck(err, "Error: Trying to get parsed items")

	var newItems []*ItemType
	newItemsCount := 0

	for i := range parsedData {
		if isNewItem(parsedData[i], savedData) {
			newItems = append(newItems, &(parsedData[i]))
			newItemsCount++
		}
	}

	if newItemsCount > 0 {
		notifier(newItems, newItemsCount)
	}

	err = saveData(parsedData, fileName)
	errorCheck(err, "Error")
}

// hell function :D
func notifier(items []*ItemType, count int) {
	var buffer bytes.Buffer
	for _, item := range items {
		buffer.WriteString("<")
		buffer.WriteString(item.Link)
		buffer.WriteString("|")
		buffer.WriteString(item.Title)
		buffer.WriteString(">")
		if item.Price > 0 {
			buffer.WriteString(" за *")
			buffer.WriteString(strconv.Itoa(item.Price))
			buffer.WriteString("* руб. в месяц (_")
			buffer.WriteString(item.Comission)
			buffer.WriteString("_)")
		}
		buffer.WriteString("\n")
		buffer.WriteString(item.Address)
		buffer.WriteString("\n\n")
	}

	m := map[string]interface{}{
		"text": buffer.String(),
	}
	mJSON, _ := json.Marshal(m)
	contentReader := bytes.NewReader(mJSON)
	http.Post(*webhookURL, "application/json", contentReader)
}

func isNewItem(item ItemType, items []ItemType) bool {
	for _, oldItem := range items {
		if item.ID == oldItem.ID && item.Price >= oldItem.Price {
			return false
		}
	}
	return true
}

func loadData(file string) ([]ItemType, error) {
	filePath := getFilePath(file)
	exists, err := exists(filePath)
	errorCheck(err, "Error")
	if !exists {
		fmt.Printf("File %s does not exists\n", filePath)
		return []ItemType{}, nil
	}
	jsonData, err := ioutil.ReadFile(filePath)
	errorCheck(err, "Error")
	var data []ItemType
	json.Unmarshal(jsonData, &data)
	return data, nil
}

func saveData(items []ItemType, file string) error {
	data, err := json.Marshal(items)
	errorCheck(err, "Error")
	filePath := getFilePath(file)
	return ioutil.WriteFile(filePath, data, 0644)
}

func getFilePath(file string) string {
	filePath, _ := os.Getwd()
	if exists, _ := exists(filePath + "/.config"); !exists {
		err := os.Mkdir(filePath+"/.config", 0755)
		errorCheck(err, "Error")
	}
	filePath += "/.config/" + file

	return filePath
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getParsedItems(url string) ([]ItemType, error) {
	doc, err := goquery.NewDocument(url)
	errorCheck(err, "Error")

	var items []ItemType

	doc.Find(".catalog-list .item").Each(func(i int, s *goquery.Selection) {
		id, errGoquery := s.Attr("id")
		if !errGoquery {
			err = errors.New("Can not get attribute 'id'")
		}

		a := s.Find("h3.title a")
		title := strings.TrimSpace(a.Text())
		href, errHref := a.Attr("href")
		if !errHref {
			err = errors.New("Can not get attribute 'href'")
		}
		link := link + href
		commission := strings.TrimSpace(s.Find(".about__commission").Text())

		// yet another hell
		html, _ := s.Find(".about").Html()
		match := regexp.MustCompile(`.*<(div|a)`).FindString(html)
		price, _ := strconv.Atoi(regexp.MustCompile(`[^\d]`).ReplaceAllString(match, ""))

		address := strings.TrimSpace(s.Find(".address").Text())

		items = append(items, ItemType{id, title, link, address, price, commission})
	})

	return items, err
}

func errorCheck(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
