package crypto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/config"
)

var symbolMap = map[string]string{
	"bitcoin":  "BTC",
	"ethereum": "ETH",
	"solana":   "SOL",
	"dogecoin": "DOGE",
	"cardano":  "ADA",
	"polkadot": "DOT",
	"ripple":   "XRP",
}

func FetchCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		ids := strings.Join(cfg.CryptoCoins, ",")
		url := fmt.Sprintf(
			"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd&include_24hr_change=true&include_market_cap=true&include_24hr_vol=true",
			ids,
		)

		resp, err := http.Get(url)
		if err != nil {
			return ResponseMsg{Error: fmt.Errorf("request failed: %w", err)}
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return ResponseMsg{Error: fmt.Errorf("API returned %d", resp.StatusCode)}
		}

		var data apiResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return ResponseMsg{Error: fmt.Errorf("decode failed: %w", err)}
		}

		var coins []CoinData
		for _, id := range cfg.CryptoCoins {
			if coin, ok := data[id]; ok {
				sym := symbolMap[id]
				if sym == "" {
					sym = strings.ToUpper(id)
					if len(sym) > 4 {
						sym = sym[:4]
					}
				}
				coins = append(coins, CoinData{
					ID:        id,
					Symbol:    sym,
					Price:     coin.USD,
					Change24h: coin.Change24h,
					MarketCap: coin.MarketCap,
					Volume24h: coin.Volume24h,
				})
			}
		}

		return ResponseMsg{Coins: coins}
	}
}
