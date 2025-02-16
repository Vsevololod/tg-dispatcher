package lib

import "strings"

func GetVideoIdFromUrl(url string) string {
	temp := url[strings.LastIndex(url, "/")+1:]
	return temp[:strings.Index(temp, "?")]

}
