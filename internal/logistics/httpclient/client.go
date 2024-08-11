package httpclient

import (
	"encoding/json"
	"fmt"
	"github.com/ivanbulyk/logistics_engine_api/internal/logging"
	"io"
	"log/slog"
	"net/http"
)

func FetchMetrics(log *slog.Logger, endpoint string) error {
	const opLabel = "httpclient.FetchMetrics"

	log = log.With(
		slog.String("opLabel", opLabel),
	)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", opLabel, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", opLabel, err)
	}

	if resp.StatusCode != http.StatusOK {
		httpErr := map[string]any{}
		if err := json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
			return err
		}
		return fmt.Errorf("service responded with non OK status code: %s", httpErr["error"])
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("got error reading response body ", logging.Err(err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("got error closing response body ", logging.Err(err))
		}
	}(resp.Body)

	_ = body
	//log.Println(string(body))

	return nil
}
