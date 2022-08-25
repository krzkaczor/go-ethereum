package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type AddressVerifier struct {
	blacklist  map[common.Address]bool
	ofacApiUrl string
}

func NewAddressVerifier(blackListPath string, ofacApiUrl string) (*AddressVerifier, error) {
	blacklist, err := readBlacklistFile(blackListPath)
	if err != nil {
		return nil, err
	}

	return &AddressVerifier{blacklist, ofacApiUrl}, nil
}

func (a *AddressVerifier) IsAddressAllowed(addr common.Address) (bool, error) {
	if isBlacklisted(a.blacklist, addr) {
		return false, nil
	}

	if a.ofacApiUrl != "" {
		return checkAPI(a.ofacApiUrl, addr)
	}
	return true, nil
}

type BlacklistedAddressesJson []common.Address

func readBlacklistFile(path string) (map[common.Address]bool, error) {
	if path == "" {
		return make(map[common.Address]bool), nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var contents BlacklistedAddressesJson
	if err := json.Unmarshal(bytes, &contents); err != nil {
		return nil, err
	}

	blacklist := make(map[common.Address]bool, len(contents))
	for _, address := range contents {
		blacklist[address] = true
	}

	log.Info("Read blacklist from file: ", "length", len(contents))

	return blacklist, nil
}

func isBlacklisted(blacklist map[common.Address]bool, addr common.Address) bool {
	return blacklist[addr]
}

func checkAPI(apiUrl string, addr common.Address) (bool, error) {
	// @todo add support to blacklist config
	url := apiUrl + addr.String()
	log.Info("Validating address OFAC status", "addr", addr, "url", url)
	resp, err := getOrRetry(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result response
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result.AddressAllowed, nil
}

type response struct {
	AddressAllowed bool `json:"addressAllowed"`
}

func getOrRetry(url string) (*http.Response, error) {
	var lastErr error

	for i := 5; i > 0; i-- {
		resp, err := http.Get(url)
		lastErr = err

		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		return resp, nil
	}
	return nil, lastErr
}
