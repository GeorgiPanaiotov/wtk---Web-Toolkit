package spider

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func Fetch(url string) (*http.Response, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return res, err
}

func ExtractLinks(base *url.URL, doc *html.Node) ([]string, error) {
	var links []string
	var walk func(*html.Node)

	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := base.Parse(attr.Val)
					if err != nil {
						continue
					}

					u.Fragment = ""
					u.RawQuery = ""

					links = append(links, u.String())
					break
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(doc)
	return links, nil
}
