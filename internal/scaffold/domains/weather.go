package domains

type Weather struct {
	City        string
	State       string
	Country     string
	Temperature float64
	LastUpdated string
}

func NewWeather(city, state, country string, temperature float64, lastUpdated string) Weather {
	return Weather{
		City:        city,
		State:       state,
		Country:     country,
		Temperature: temperature,
		LastUpdated: lastUpdated,
	}
}

func (w *Weather) GetCelcius() float64 {
	return w.Temperature
}

func (w *Weather) GetFahrenheit() float64 {
	return w.Temperature*1.8 + 32
}

func (w *Weather) GetKelvin() float64 {
	return w.Temperature + 273.15
}
