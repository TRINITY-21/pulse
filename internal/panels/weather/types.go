package weather

import "time"

type apiResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Visibility int `json:"visibility"`
	Sys        struct {
		Sunrise int64 `json:"sunrise"`
		Sunset  int64 `json:"sunset"`
	} `json:"sys"`
}

type Data struct {
	City        string
	Temp        float64
	FeelsLike   float64
	TempMin     float64
	TempMax     float64
	Condition   string
	Description string
	Humidity    int
	Pressure    int
	WindSpeed   float64
	Clouds      int
	Visibility  float64 // km
	Sunrise     int64
	Sunset      int64
}

type ResponseMsg struct {
	Data  Data
	Error error
}

type TickMsg time.Time
