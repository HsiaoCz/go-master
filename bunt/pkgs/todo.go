package pkgs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/HsiaoCz/go-master/bunt/types"
)

// get location by ip

func GetLocationByIP(ip string) (string, error) {
	// get location by ip
	// ...
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var location types.IPLocation
	err = json.Unmarshal(body, &location)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	return location.City, nil
}
