package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"pulse/internal/config"
)

func FetchCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf(
			"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
			cfg.WeatherCity, cfg.WeatherAPIKey,
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

		condition := ""
		description := ""
		if len(data.Weather) > 0 {
			condition = data.Weather[0].Main
			description = data.Weather[0].Description
		}

		return ResponseMsg{
			Data: Data{
				City:        data.Name,
				Temp:        data.Main.Temp,
				FeelsLike:   data.Main.FeelsLike,
				TempMin:     data.Main.TempMin,
				TempMax:     data.Main.TempMax,
				Condition:   condition,
				Description: description,
				Humidity:    data.Main.Humidity,
				Pressure:    data.Main.Pressure,
				WindSpeed:   data.Wind.Speed,
				Clouds:      data.Clouds.All,
				Visibility:  float64(data.Visibility) / 1000.0,
				Sunrise:     data.Sys.Sunrise,
				Sunset:      data.Sys.Sunset,
			},
		}
	}
}
