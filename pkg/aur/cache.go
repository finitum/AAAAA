package aur

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// NotInAurErr will be returned when the package is not found in the AUR, meaning this dependency
// should not be considered for resolving dependencies.
var NotInAurErr = errors.New("pkg not in aur")

// SendCachedInfoRequest will do a HTTP GET request to the given url, after replacing the '%s' with the given package
// name. It expects to receive a InfoResult, or a 204 No Content (to signal that the package is not on the AUR).
func SendCachedInfoRequest(url, pkg string) (res InfoResult, err error) {
	resp, err := http.Get(fmt.Sprintf(url, pkg))
	if err != nil {
		return res, errors.Wrap(err, "received error from aur cache")
	}
	if resp.StatusCode == http.StatusNoContent {
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
