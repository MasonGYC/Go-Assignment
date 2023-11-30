package main

type Page struct {
	id      int
	content any
}

func NewPage(id int, content any) Page {
	return Page{
		id:      id,
		content: content,
	}
}
