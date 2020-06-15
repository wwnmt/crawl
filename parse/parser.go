package parse

import (
	"crawl/engine"
	"crawl/model"
	"regexp"
	"strings"
)

var tagRe = regexp.MustCompile(`<a href="([^"]+)">([^"]+)</a>`)
var bookListRe = regexp.MustCompile(`<a href="([^"]+)" title="([^"]+)"`)
var authorRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?">([^<]+)</a>`)
var publishRe = regexp.MustCompile(`<span class="pl">出版社:</span> ([^<"]+)[\d\D]*?<br/>`)
var dateRe = regexp.MustCompile(`<span class="pl">出版年:</span> ([^<"]+)[\d\D]*?<br/>`)
var pagesRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<"]+)[\d\D]*?<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span> ([^<"]+)[\d\D]*?<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var isbnRe = regexp.MustCompile(`<span class="pl">ISBN:</span> ([^<"]+)[\d\D]*?<br/>`)
var infoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<"]+)</p></div>`)

func NilParser([]byte) engine.ParseResult {
	return engine.ParseResult{}
}

func TagParser(contents []byte) engine.ParseResult {
	match := tagRe.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}

	for _, m := range match {
		temp := string(m[1])
		if strings.HasPrefix(temp, "/tag/") {
			result.Items = append(result.Items, m[2])
			result.Requests = append(result.Requests, engine.Request{
				Url:       "https://book.douban.com" + string(m[1]),
				ParseFunc: BookListParser,
			})
		}
	}
	return result
}

func BookListParser(contents []byte) engine.ParseResult {
	match := bookListRe.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}

	for _, m := range match {
		bookName := string(m[2])
		result.Items = append(result.Items, bookName)
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			ParseFunc: func(bytes []byte) engine.ParseResult {
				return BookParser(bytes, bookName)
			},
		})
	}
	return result
}

func extraString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

func BookParser(contents []byte, bookName string) engine.ParseResult {
	var book = model.Book{
		Name:        bookName,
		Author:      extraString(contents, authorRe),
		Publisher:   extraString(contents, publishRe),
		PublishDate: extraString(contents, dateRe),
		Score:       strings.Replace(extraString(contents, scoreRe), " ", "", -1),
		Pages:       extraString(contents, pagesRe),
		Price:       extraString(contents, priceRe),
		ISBN:        extraString(contents, isbnRe),
		Info:        strings.Replace(extraString(contents, infoRe), " ", "", -1),
	}

	result := engine.ParseResult{
		Items: []interface{}{book},
	}

	return result
}
