package common

import (
	"net/url"
	"strings"

	"github.com/QOSGroup/cassini/log"
)

// ParseUrls parse URLs
func ParseUrls(firstURL, secondURL string) (fus, sus []url.URL, err error) {
	fus, err = parseURLs(firstURL)
	if err != nil {
		log.Error("Parse url error: ", err)
		return nil, nil, err
	}
	if strings.EqualFold(secondURL, "") {
		sus = fus
	} else {
		sus, err = parseURLs(secondURL)
		if err != nil {
			log.Error("Parse url error: ", err)
			return nil, nil, err
		}

	}
	return
}

// parseURLs parse URL
func parseURLs(urls string) ([]url.URL, error) {
	var err error
	us := strings.Split(urls, ",")
	ret := make([]url.URL, len(us))
	var ur *url.URL
	for i, u := range us {
		ur, err = url.Parse(u)
		if err != nil {
			return nil, err
		}
		ret[i] = *ur
	}
	return ret, nil
}
