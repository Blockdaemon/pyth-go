//  Copyright 2022 Blockdaemon Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	//go:embed tests/BmA9Z6FjioHJPpjT39QazZyhDRUdZy2ezwx4GiDdE2u2.bin
	caseMappingAccount []byte
)

func TestProductAccount(t *testing.T) {
	expected := ProductAccount{
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

	var actual ProductAccount
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

func TestMappingAccount(t *testing.T) {
	expected := MappingAccount{
		AccountHeader: AccountHeader{
			Magic:       Magic,
			Version:     V2,
			AccountType: AccountTypeMapping,
			Size:        2168,
		},
		Num: 66,
		Products: [640]solana.PublicKey{
			solana.MustPublicKeyFromBase58("89GseEmvNkzAMMEXcW9oTYzqRPXTsJ3BmNerXmgA1osV"),
			solana.MustPublicKeyFromBase58("JCnD5WiurZfoeVPEi2AXVgacg73Wd2iRDDjZDbSwdr9D"),
			solana.MustPublicKeyFromBase58("G89jkM5wFLpmnbvRbeePUumxsJyzoXaRfgBVjyx2CPzQ"),
			solana.MustPublicKeyFromBase58("GaBJpKtnyUbyKe34XuyegR7W98a9PT5cg985G974NY8R"),
			solana.MustPublicKeyFromBase58("Fwosgw2ikRvdzgKcQJwMacyczk3nXgoW3AtVtyVvXSAb"),
			solana.MustPublicKeyFromBase58("Bjn6jsUJKEe8JDR4aTHhRRDZzeKHMavvnF9UWVFPaFqE"),
			solana.MustPublicKeyFromBase58("CpPmHbFqkfejPcF8cvxyDogm32Sqo3YGMFBgv3kR1UtG"),
			solana.MustPublicKeyFromBase58("6U4PrvMwfMcBkG7Zrc4oxYqJwrfTMWmgA9hS6fjDJkmo"),
			solana.MustPublicKeyFromBase58("7MudLeJnT2GCPZ66oeAqd6jenF9fGARrB1pLo5nBT3KM"),
			solana.MustPublicKeyFromBase58("J6zuHzycf8XLd85QHDUwMtVPxJGPJptPSC9dyioKXCnb"),
			solana.MustPublicKeyFromBase58("4nyATHv6KnZY5fVTqQLq9DkcstfqGYd834Jmbch2bf3i"),
			solana.MustPublicKeyFromBase58("3m1y5h2uv7EQL3KaJZehvAJa4yDNvgc5yAdL9KPMKwvk"),
			solana.MustPublicKeyFromBase58("2ciUuGZiee5macAMeQ7bHGTJtwcYTgnt6jdmQnnKZrfu"),
			solana.MustPublicKeyFromBase58("3Mnn2fX6rQyUsyELYms1sBJyChWofzSNRoqYzvgMVz5E"),
			solana.MustPublicKeyFromBase58("6MEwdxe4g1NeAF9u6KDG14anJpFsVEa2cvr5H6iriFZ8"),
			solana.MustPublicKeyFromBase58("6NpdXrQEpmDZ3jZKmM2rhdmkd3H6QAk23j2x8bkXcHKA"),
			solana.MustPublicKeyFromBase58("2weC6fjXrfaCLQpqEzdgBHpz6yVNvmSN133m7LDuZaDb"),
			solana.MustPublicKeyFromBase58("4zvUzWGBxZA9nTgBZWAf1oGYw6nCEYRscdt14umTNWhM"),
			solana.MustPublicKeyFromBase58("C5wDxND9E61RZ1wZhaSTWkoA8udumaHnoQY6BBsiaVpn"),
			solana.MustPublicKeyFromBase58("25tCF4ChvZyNP67xwLuYoAKuoAcSV13xrmP9YTwSPnZY"),
			solana.MustPublicKeyFromBase58("EWxGfxoPQSNA2744AYdAKmsQZ8F9o9M7oKkvL3VM1dko"),
			solana.MustPublicKeyFromBase58("CiTV5gD8G53M1EQdo32jy5riYRU8fMFSVWC5wJj3vjcr"),
			solana.MustPublicKeyFromBase58("3K2hkXeoxNeRgjGTU6unJ4WSRaZ3FZxhABTgk8wcbPpX"),
			solana.MustPublicKeyFromBase58("6NF21VSjK5qt5t6JZXZtMZ1kXEwTuDBqSYL2ev7dMgx3"),
			solana.MustPublicKeyFromBase58("8zDnpUALoDEZoufcVPjTZSpLXpWVZJbRQGEXVMBy14SW"),
			solana.MustPublicKeyFromBase58("BzJGxCqttFdmu3J1MC7N8qAhPHgf6ZHHdoDcX91idXLK"),
			solana.MustPublicKeyFromBase58("63VWd2FVbukVozZ1okHt8wVMq7enAYFXYnmp2DUQtBBJ"),
			solana.MustPublicKeyFromBase58("APY6HpicQTabn46sZPt2pAqunFpkws9ySdujj5BHcZAK"),
			solana.MustPublicKeyFromBase58("6C4PJ4bMuLFmvHRqSkmGeyoSGAKMfPG1um1k1suryfs"),
			solana.MustPublicKeyFromBase58("3BtxtRxitVDcsd7pPUWUnFm9KvmNDy9usS4gE6pUFhpH"),
			solana.MustPublicKeyFromBase58("8g9qN2XBoTr53dcescpHjQUkhKL6pHrcHzHQ9MWpkLJa"),
			solana.MustPublicKeyFromBase58("Ch1iziUNx4japBFbJd7BzapHuHMBJ6rtoNzdjDixKqaW"),
			solana.MustPublicKeyFromBase58("Pc6A6JVVSXZosyFTtykoQBf5TZ1hhfRTcuo6Zf88AUg"),
			solana.MustPublicKeyFromBase58("EssaQC37YW2LVXTsEVjijNte3FTUw21MJBkZYHDdyakc"),
			solana.MustPublicKeyFromBase58("2UE6gC5FuVPWuKqZamRfcEc5MjtvpRoW6L1anCGW4skS"),
			solana.MustPublicKeyFromBase58("5oXPa1o1fN7rpCUE2cgKFD9jYh4m4kFPjxjsjiHTF3Nr"),
			solana.MustPublicKeyFromBase58("C9Ua6p5Db4MU6E2umizxz7ehRuRRPZNyYbu4LK6G5oUn"),
			solana.MustPublicKeyFromBase58("FCXufwMoZhytNjKWVrJ2cEGbexaFD1nMDk7GqJRin1Rg"),
			solana.MustPublicKeyFromBase58("HBH8e7ZGWaPaNW17LSHFQPdAXEaU8ejAmCYkr1566eAt"),
			solana.MustPublicKeyFromBase58("DH9LThgjEXXvk6avCGJNHVd46bdgtuj9xVk2d9kqQpnW"),
			solana.MustPublicKeyFromBase58("BvcxLiPtnz6yfaC1TfPQCQnxjqRfb4CRXxZj4i5Wh8FM"),
			solana.MustPublicKeyFromBase58("9SKi6wLvdo8A45RqjwF1ZMJixteS7tsAxYBkvJgvxW6Z"),
			solana.MustPublicKeyFromBase58("B68oPzvPMNjrdmMEAX9zT1ucqm4R3QkrquMjPQjrPmz2"),
			solana.MustPublicKeyFromBase58("F6mXkFitT5T9MGP1QS5J3aDufN8xFVeNykPDTohoAitw"),
			solana.MustPublicKeyFromBase58("76B8fdtbYnpba2io43rt7MpAQHCe2objsc637f8auC3G"),
			solana.MustPublicKeyFromBase58("31HTfSgBs7PJmY6YgRKaA3ionPmBHrzPnbwZSqRGs2Zx"),
			solana.MustPublicKeyFromBase58("4Yprdh5xpNgpsuDTPmfxn1ky7YjXmioZ4h8vGdCaBDsE"),
			solana.MustPublicKeyFromBase58("os3is9HtWPHW4EXpGAkdr2prdWVs2pS8qKtf2ZYJdBw"),
			solana.MustPublicKeyFromBase58("A8Q4MoqpiEp4zU36XoVB3c9q6XkeqM5uaM9LQv4p9vQy"),
			solana.MustPublicKeyFromBase58("Fssu4winxvACKTyT28hDAsonR8DWsBtrButRH2KL77x2"),
			solana.MustPublicKeyFromBase58("Hzs3LGujZGkqyLVkad6w6CkjhzwDo2srQ3Tk4QPCoAD2"),
			solana.MustPublicKeyFromBase58("71k9hopyryKPUWug1iKiJCkbEsz1C7EptMN2t1dtgNmA"),
			solana.MustPublicKeyFromBase58("AM3Sf65EAMRe4MofThv2C8jhfQ4gJ1Q9HCz38Mybojt"),
			solana.MustPublicKeyFromBase58("DDdPuysfkxPq5Y1ZtTSk1H5n7iBKc9wtEKUwd1TNu3Gc"),
			solana.MustPublicKeyFromBase58("BTj5x8YZL5F8Z16zbXfDrN7zhZKccQN4F1H7M8W3XprK"),
			solana.MustPublicKeyFromBase58("GyMxxTi4EHhtropCenKb5wdPkHHqgczgMJPhXkdDtiK9"),
			solana.MustPublicKeyFromBase58("CLsD9kiEs7LVyZpMH2mTUmtwLZ1C1dqpf5HhMHDQcRCJ"),
			solana.MustPublicKeyFromBase58("5kWV4bhHeZANzg5MWaYCQYEEKHjur5uz1mu5vuLHwiLB"),
			solana.MustPublicKeyFromBase58("EmQ2KSgm6uPp7PgiuYxibdMKxoPqpJC38Ti87ujNgteq"),
			solana.MustPublicKeyFromBase58("Xg61xoGzEhiJUxXvLRaodYEZRSVFRem9KUVDc89wzvH"),
			solana.MustPublicKeyFromBase58("HgFFJ9vGcjMqgAbeN4xuQLtoyPuecJbdvo6oMDCnEUnk"),
			solana.MustPublicKeyFromBase58("DfkWGYqPRdVrJeFtYHSPGUueW5YnUZbW8C4KTCJ67Dsp"),
			solana.MustPublicKeyFromBase58("6KrH5AMu8frBgJ7hbRD4L9w68LVurbdJVXLachBdi3q4"),
			solana.MustPublicKeyFromBase58("4dQPbzdUQFa3YquQQbndRv1oTtbC2BUEFM3KwbTixYAn"),
			solana.MustPublicKeyFromBase58("7UHB783Nh4avW3Yw9yoktf2KjxipU56KPahA51RnCCYE"),
			solana.MustPublicKeyFromBase58("C6wHqCjtiHkTY3wKdehup1sScb6R6e8YV7dkK3P17w98"),
		},
	}

	var actual MappingAccount
	require.NoError(t, actual.UnmarshalBinary(caseMappingAccount))

	assert.Equal(t, &expected, &actual)
}
