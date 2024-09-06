// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

type Site struct {
	Title       string
	Description string
	IconUrl     string
	Image       string
}

func hasKeyWithValue(attributes map[string]string, key, value string) bool {
	if val, ok := attributes[key]; ok {
		if val == value {
			return true
		}
	}
	return false
}

func attrToMap(attr []html.Attribute) map[string]string {
	r := make(map[string]string, len(attr))
	for _, a := range attr {
		r[a.Key] = a.Val
	}
	return r
}

func ExtractData(tok html.Token, site Site) Site {
	var attrMap map[string]string
	if tok.Data == "meta" {
		if attrMap == nil {
			attrMap = attrToMap(tok.Attr)
		}

		// fmt.Println("tok.Attr =>", tok.Attr)
		// fmt.Println("attrMap =>", attrMap)

		if hasKeyWithValue(attrMap, "property", "og:description") || hasKeyWithValue(attrMap, "name", "description") {
			// we have the check above, thus we skip error handling here
			site.Description, _ = attrMap["content"]
		}
		if hasKeyWithValue(attrMap, "property", "og:image") {
			site.Image, _ = attrMap["content"]
		}
	} else if tok.Data == "link" {
		if attrMap == nil {
			attrMap = attrToMap(tok.Attr)
		}
		if hasKeyWithValue(attrMap, "type", "image/x-icon") {
			site.IconUrl, _ = attrMap["href"]
		}
	}

	return site
}

func Extract(r io.Reader) Site {
	site := Site{}
	lexer := html.NewTokenizer(r)

	var inHead bool
	var inTitle bool
	for {
		tokenType := lexer.Next()
		tok := lexer.Token()

		// fmt.Println("tok.Attr =>", tok.Attr)
		// fmt.Println("tok.Data =>", tok.Data)

		// fmt.Println("raw =>", raw)
		// fmt.Println("tok =>", tok)
		// fmt.Println("tokenType =>", tokenType)

		site = ExtractData(tok, site)

		switch tokenType {
		case html.StartTagToken: // <head>
			if tok.Data == "head" {
				inHead = true
			}

			// keep lexing if not in head of document
			if !inHead {
				continue
			}

			if tok.Data == "title" {
				inTitle = true
			}
		case html.TextToken:
			if inTitle {
				site.Title = tok.Data
				inTitle = false
			}
		case html.EndTagToken: // </head>
			if tok.Data == "head" {
				return site
			}
		}
	}
}

func main() {
	client := http.Client{
		Timeout: time.Millisecond * 250,
	}
	resp, err := client.Get("https://github.com/")
	if err != nil {
		fmt.Println("err =>", err)
	} else if resp.StatusCode > http.StatusPermanentRedirect {
		fmt.Println("resp.StatusCode =>", err)
	}

	site := Extract(resp.Body)
	err = resp.Body.Close()

	fmt.Println("site.Title =>", site.Title)
	fmt.Println("site.Description =>", site.Description)
	fmt.Println("site.IconUrl =>", site.IconUrl)
	fmt.Println("site.Image =>", site.Image)
}
