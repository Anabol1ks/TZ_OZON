package e2e

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExchangeRateE2E(t *testing.T) {
	for i := 0; i < 20; i++ {
		i := i
		t.Run(fmt.Sprintf("date-%d", i), func(t *testing.T) {
			t.Parallel()
			date := time.Now().AddDate(0, 0, i).Format("02/01/2006")
			url := fmt.Sprintf("http://localhost:8080/scripts/XML_daily.asp?date_req=%s", date)
			resp, err := http.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				require.Contains(t, string(body), "<ValCurs")
			} else if resp.StatusCode == 500 {
				require.Contains(t, string(body), "Internal Server Error")
			} else {
				t.Fatalf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
			}
		})
	}
}
