package hook

import (
	"encoding/json"
	"io"
	"net/http"
)

func ApiCall(url string, typ any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return http.ErrHandlerTimeout
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, typ)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, typ)
}
