package delivery

import (
	"fmt"
	"net/http"
	"time"
)

var (
	QUERY_URL = "http://www.kuai" + "di100.com/query?type=%s&postid=%s&id=1&valicode=&temp=%d"
)

func QueryRemote(com, nu string) (*http.Response, error) {
	//request like a browser
	url := fmt.Sprintf(QUERY_URL, com, nu, time.Now().Unix())
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "http://www.kuai"+"di100.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:47.0) Gecko/20100101 Firefox/47.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	return http.DefaultClient.Do(req)
}
