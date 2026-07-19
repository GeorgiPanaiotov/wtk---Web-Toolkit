package spider

import (
	"net/http"
)

func Fetch(url string) (*http.Response, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return res, err
}
