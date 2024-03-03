package viacep

import "github.com/stretchr/testify/mock"

type MockViaCEP struct {
	mock.Mock
}

func (m *MockViaCEP) GetAddress(cep string) (Response, error) {
	args := m.Called(cep)
	return args.Get(0).(Response), args.Error(1)
}
