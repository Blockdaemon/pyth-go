package pyth

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetPriceAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		buf, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"jsonrpc": "2.0",
			"id": 0,
			"method": "getAccountInfo",
			"params": [
				"E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh",
				{"encoding": "base64"}
			]
		}`, string(buf))

		_, err = wr.Write([]byte(`{
			"jsonrpc": "2.0",
			"id": 0,
			"result": {
				"context": {
					"slot": 118773287
				},
				"value": {
					"data": [
						"` + base64.StdEncoding.EncodeToString(casePriceAccount) + `",
						"base64"
					],
					"executable": false,
					"lamports": 23942400,
					"owner": "gSbePebfvPy7tRqimPoVecS2UsBvYv46ynrzWocc92s",
					"rentEpoch": 274
				}
			}
		}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := NewClient(ProgramIDDevnet, server.URL, server.URL)
	acc, err := c.GetPriceAccount(context.Background(), solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh"))
	require.NoError(t, err)
	require.NotNil(t, acc)

	assert.Equal(t, &priceAccount_E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh, acc)
}
