package weather

import "github.com/stretchr/testify/mock"

type MockWeather struct {
	mock.Mock
}

func (m *MockWeather) GetWeather(city string, state string) (Response, error) {
	args := m.Called(city, state)
	return args.Get(0).(Response), args.Error(1)
}
