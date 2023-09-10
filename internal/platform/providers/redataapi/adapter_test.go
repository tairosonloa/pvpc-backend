package redataapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"pvpc-backend/internal/domain"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func Test_FetchPVPCPrices_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, pvpcPricesEndpoint, r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		require.Equal(t, "end_date=2023-09-08T23%3A59&geo_ids=1234&start_date=2023-09-08T00%3A00&time_trunc=hour", r.URL.RawQuery)

		res, err := os.ReadFile("./mocks/fetch_pvpc_response.json")
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}))
	defer server.Close()

	zone, err := domain.NewZone(domain.ZoneDto{ID: "ZON", ExternalID: "1234", Name: "Zone Name"})
	require.NoError(t, err)

	date, err := time.Parse("2006-01-02T15:04:05Z", "2023-09-08T17:54:36Z")
	require.NoError(t, err)

	adapter := NewREDataAPI(server.URL)
	prices, err := adapter.FetchPVPCPrices(context.Background(), []domain.Zone{zone}, date)
	require.NoError(t, err)
	require.Len(t, prices, 1)
	snaps.MatchSnapshot(t, prices[0].Serialize())
}

func Test_FetchPVPCPrices_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, pvpcPricesEndpoint, r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		require.Equal(t, "end_date=2023-09-08T23%3A59&geo_ids=1234&start_date=2023-09-08T00%3A00&time_trunc=hour", r.URL.RawQuery)

		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("error"))
	}))
	defer server.Close()

	zone, err := domain.NewZone(domain.ZoneDto{ID: "ZON", ExternalID: "1234", Name: "Zone Name"})
	require.NoError(t, err)

	date, err := time.Parse("2006-01-02T15:04:05Z", "2023-09-08T17:54:36Z")
	require.NoError(t, err)

	adapter := NewREDataAPI(server.URL)
	prices, err := adapter.FetchPVPCPrices(context.Background(), []domain.Zone{zone}, date)
	require.NoError(t, err)
	require.Len(t, prices, 0)
}
