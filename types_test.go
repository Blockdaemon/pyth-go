package pyth

import (
	_ "embed"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed tests/EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko.bin
	caseProductAccount []byte
	//go:embed tests/E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh.bin
	casePriceAccount []byte
)

func TestProductAccount(t *testing.T) {
	expected := Product{
		AccountHeader: AccountHeader{
			Magic:       Magic,
			Version:     V2,
			AccountType: AccountTypeProduct,
			Size:        161,
		},
		FirstPrice: solana.MustPublicKeyFromBase58("E36MyBbavhYKHVLWR79GiReNNnBDiHj6nWA7htbkNZbh"),
		Attrs: [464]byte{
			0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x0a,
			0x46, 0x58, 0x2e, 0x45, 0x55, 0x52, 0x2f, 0x55,
			0x53, 0x44, 0x0a, 0x61, 0x73, 0x73, 0x65, 0x74,
			0x5f, 0x74, 0x79, 0x70, 0x65, 0x02, 0x46, 0x58,
			0x0e, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x5f, 0x63,
			0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x03,
			0x55, 0x53, 0x44, 0x0b, 0x64, 0x65, 0x73, 0x63,
			0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x07,
			0x45, 0x55, 0x52, 0x2f, 0x55, 0x53, 0x44, 0x0e,
			0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x5f,
			0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x06, 0x45,
			0x55, 0x52, 0x55, 0x53, 0x44, 0x04, 0x62, 0x61,
			0x73, 0x65, 0x03, 0x45, 0x55, 0x52, 0x05, 0x74,
			0x65, 0x6e, 0x6f, 0x72, 0x04, 0x53, 0x70, 0x6f,
			0x74, 0x53, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	}

	var actual Product
	require.NoError(t, actual.UnmarshalBinary(caseProductAccount))

	assert.Equal(t, &expected, &actual)

	t.Run("GetAttrs", func(t *testing.T) {
		expected := map[string]string{
			"asset_type":     "FX",
			"base":           "EUR",
			"description":    "EUR/USD",
			"generic_symbol": "EURUSD",
			"quote_currency": "USD",
			"symbol":         "FX.EUR/USD",
			"tenor":          "Spot",
		}
		actual, err := actual.GetAttrs()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestPriceAccount(t *testing.T) {
	expected := PriceAccount{
		AccountHeader: AccountHeader{
			Magic:       Magic,
			Version:     V2,
			AccountType: AccountTypePrice,
			Size:        1200,
		},
		PriceType: 1,
		Exponent:  -5,
		Num:       10,
		NumQt:     0,
		LastSlot:  117136050,
		ValidSlot: 117491486,
		Twap: Ema{
			Val:   112674,
			Numer: 5644642336,
			Denom: 5009691136,
		},
		Twac: Ema{
			Val:   4,
			Numer: 2033641276,
			Denom: 5009691136,
		},
		Drv1:      1,
		Drv2:      0,
		Product:   solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko"),
		Next:      solana.PublicKey{},
		PrevSlot:  117491485,
		PrevPrice: 112717,
		PrevConf:  6,
		Drv3:      -2413575930482041166,
		Agg: PriceInfo{
			Price:   112717,
			Conf:    6,
			Status:  0,
			CorpAct: 0,
			PubSlot: 117491487,
		},
		Components: [32]PriceComp{
			{
				Publisher: solana.MustPublicKeyFromBase58("5U3bH5b6XtG99aVWLqwVzYPVpQiFHytBD68Rz2eFPZd7"),
				Agg: PriceInfo{
					PubSlot: 117491484,
				},
				Latest: PriceInfo{
					PubSlot: 117491485,
				},
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("4iVm6RJVU4R6kvc3KUDnE6cw4Ffb6769FzbXMu26sJrs"),
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("3djmXHmD9kuAydgFnSnWAjq4Kos5GnEx2KdFR2kvGiUw"),
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("86DsXwBCqFoCUiuB1t9oV2inHKQ5h2vFaNZ4GETvTHuz"),
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("rkTtobRtTCDLXbADsbVxHcfBr7Z8Z1JDSBM3kyk3LJe"),
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("2pfE7YYVhM9WaneVVF2kcwArMoconfjtq83oZfSurkkY"),
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("2vTC3XNpi7ED5T643KxVH5HqM7cSRKuUGnmMtKACY4Ju"),
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("45FYxKkPM1NhavyAHFTyXG2JCSsy5jD1UwwCz5UtHX5y"),
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("EevTjv14eGHqsxKvgpastHsuLr9FNPfzkP23wG61pT2U"),
				Agg: PriceInfo{
					Price:   113062,
					Conf:    1,
					Status:  1,
					CorpAct: 0,
					PubSlot: 116660829,
				},
				Latest: PriceInfo{
					Price:   113062,
					Conf:    1,
					Status:  1,
					CorpAct: 0,
					PubSlot: 116660829,
				},
			},
			{
				Publisher: solana.MustPublicKeyFromBase58("AKPWGLY5KpxbTx7DaVp4Pve8JweMjKbb1A19MyL2nrYT"),
				Agg: PriceInfo{
					Price:   111976,
					Conf:    16,
					Status:  1,
					CorpAct: 0,
					PubSlot: 116917242,
				},
				Latest: PriceInfo{
					Price:   111976,
					Conf:    16,
					Status:  1,
					CorpAct: 0,
					PubSlot: 116917242,
				},
			},
		},
	}

	var actual PriceAccount
	require.NoError(t, actual.UnmarshalBinary(casePriceAccount))

	assert.Equal(t, &expected, &actual)
}
