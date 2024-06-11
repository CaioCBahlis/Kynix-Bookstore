package Scrappers

import (
	b64 "encoding/base64"
	"strings"
	"time"
	"github.com/gocolly/colly/v2"
)

type Book struct{
	Title string `json:"title"`
	IMGURL string `json:"imgurl"`
	Link1 string  `json:"link1"`
	Link2 string  `json:"link2"`
	Link3 string  `json:"link3"`
	Link4 string  `json:"link4"`
	Lupdate time.Time `json:"lupdate"`
}


func ZLibrary_Scrapper(c *colly.Collector, url string)  (ScrappedBook Book) {
	
	MyBook := Book{}

	c.OnHTML("div.card.mt-2", func(h *colly.HTMLElement) {
		
		

		IMGUrlParent := h.DOM.Find("div.card-body.text-center")
		Title := h.DOM.Find("h1").Text()
		IMGURL, _ := IMGUrlParent.Find("img").Attr("src")
		Down1_Parent := h.DOM.Find("div.mt-2")
		Mirror_Parent := h.DOM.Find("div.mt-3")

		DOWN_Encoded, _ := Down1_Parent.Find("a.btn.btn-success.download_now").Attr("onclick")
		DOWN_URL := B64_Formatting(DOWN_Encoded)

		Mirror_Url, _ :=  Mirror_Parent.Find("a[id=mirror1]").Attr("href")
		//no, its not a bug, its actually right, only the first mirror has an uncoded href
		
		Mirror_Encoded2, _ := Mirror_Parent.Find("a[id=mirror2]").Attr("onclick")
		Mirror_Url2 := B64_Formatting(Mirror_Encoded2)


		Mirror_Encoded3, _ := Mirror_Parent.Find("a[id=mirror3]").Attr("onclick")
		Mirror_Url3 := B64_Formatting(Mirror_Encoded3)


		MyBook.Title = Title
		MyBook.IMGURL = IMGURL
		MyBook.Link1 = DOWN_URL
		MyBook.Link2 = Mirror_Url
		MyBook.Link3 = Mirror_Url2
		MyBook.Link4 = Mirror_Url3
		MyBook.Lupdate = time.Now()

		
		
})


	c.Visit(url)
	return MyBook
	
}


func B64_Formatting(encodedb64 string) string{
	encodedb64 = strings.TrimPrefix(encodedb64, "openLinkNewTab")
	encodedb64 = strings.Trim(encodedb64, "'()")
	link, err := b64.StdEncoding.DecodeString(encodedb64)
	if err == nil{
		return string(link[:])
	}
	return encodedb64
	
}
