package requester

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockRequester struct {
	mock.Mock
}

func (m *MockRequester) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}
