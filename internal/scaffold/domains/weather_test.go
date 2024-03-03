package domains

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWeather_GetCelcius(t *testing.T) {
	t.Run("should return the temperature in celcius", func(t *testing.T) {
		weather := NewWeather("São Paulo", "SP", "Brasil", 25, "2021-01-01 12:00:00")
		assert.Equal(t, 25.0, weather.GetCelcius())
	})

	t.Run("should return the temperature in fahrenheit", func(t *testing.T) {
		weather := NewWeather("São Paulo", "SP", "Brasil", 25, "2021-01-01 12:00:00")
		assert.Equal(t, 77.0, weather.GetFahrenheit())
	})

	t.Run("should return the temperature in kelvin", func(t *testing.T) {
		weather := NewWeather("São Paulo", "SP", "Brasil", 25, "2021-01-01 12:00:00")
		assert.Equal(t, 298.15, weather.GetKelvin())

	})
}
