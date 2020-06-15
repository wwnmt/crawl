package model

type Book struct {
	Name        string
	Author      string
	Publisher   string
	PublishDate string
	Pages       string
	Price       string
	Score       string
	Info        string
	ISBN        string
}

func (b Book) String() string {
	return "书名:" + b.Name + " 作者:" + b.Author + " 出版社:" + b.Publisher +
		" 出版日期:" + b.PublishDate + " 页数:" + b.Pages +
		" 定价:" + b.Price + " 评分:" + b.Score + " ISBN:" + b.ISBN +
		"\n简介:" + b.Info
}
