package esios

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/domain/errors"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

const MOCK_TOKEN = "FAKE_TOKEN"

func Test_FetchPVPCPrices_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, pvpcPricesEndpoint, r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("accept"))
		require.Equal(t, MOCK_TOKEN, r.Header.Get("x-api-key"))
		require.Equal(t, "end_date=2023-09-08T23%3A59%3A59&geo_ids%5B%5D=1234&geo_ids%5B%5D=5678&start_date=2023-09-08T00%3A00%3A00", r.URL.RawQuery)

		res, err := os.ReadFile("./mocks/fetch_pvpc_response.json")
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}))
	defer server.Close()

	zone1, err := domain.NewZone(domain.ZoneDto{ID: "FOO", ExternalID: "1234", Name: "Foo Zone"})
	require.NoError(t, err)
	zone2, err := domain.NewZone(domain.ZoneDto{ID: "BAR", ExternalID: "5678", Name: "Bar Zone"})
	require.NoError(t, err)

	date, err := time.Parse("2006-01-02T15:04:05Z", "2023-09-08T17:54:36Z")
	require.NoError(t, err)

	adapter := NewEsiosAPI(server.URL, MOCK_TOKEN)
	prices, err := adapter.FetchPVPCPrices(context.Background(), []domain.Zone{zone1, zone2}, date)
	require.NoError(t, err)
	require.Len(t, prices, 2)
	snaps.MatchSnapshot(t, prices[0].Serialize(), prices[1].Serialize())
}

func Test_FetchPVPCPrices_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, pvpcPricesEndpoint, r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		require.Equal(t, "end_date=2023-09-08T23%3A59%3A59&geo_ids%5B%5D=1234&start_date=2023-09-08T00%3A00%3A00", r.URL.RawQuery)

		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("error"))
	}))
	defer server.Close()

	zone, err := domain.NewZone(domain.ZoneDto{ID: "ZON", ExternalID: "1234", Name: "Zone Name"})
	require.NoError(t, err)

	date, err := time.Parse("2006-01-02T15:04:05Z", "2023-09-08T17:54:36Z")
	require.NoError(t, err)

	adapter := NewEsiosAPI(server.URL, MOCK_TOKEN)
	prices, err := adapter.FetchPVPCPrices(context.Background(), []domain.Zone{zone}, date)
	require.Error(t, err)
	require.Equal(t, errors.ProviderError, errors.Code(err))
	require.Len(t, prices, 0)
}
