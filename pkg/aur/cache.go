package aur

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var NotInAurErr = errors.New("pkg not in aur")

func SendCachedInfoRequest(url, pkg string) (res InfoResult, err error) {
	resp, err := http.Get(fmt.Sprintf(url, pkg))
	if err != nil {
		return res, errors.Wrap(err, "received error from aur cache")
	}
	if resp.StatusCode == http.StatusNotFound {
		return InfoResult{}, NotInAurErr
	}
	if resp.StatusCode != http.StatusOK {
		return InfoResult{}, errors.New(fmt.Sprintf("received error from aur cache: %s", resp.Status))
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Error(err)
		return res, errors.Wrap(err, "couldn't decode result")
	}

	return
}
