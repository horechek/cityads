package cityads

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *ClientSuite) SetupSuite() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, testApiErrorData)
	}))
}

func (s *ClientSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *ClientSuite) TestErrRequest() {
	client := NewClient(
		s.server.URL+"/api/rest/webmaster/json/",
		"xxx",
	)

	params := url.Values{}
	params.Add("user_has_offer", "true")
	rep, err := client.request("offers/web", "GET", params)

	s.Assert().NotNil(err)
	s.Assert().Nil(rep)

	e := err.(ApiError)

	s.Assert().EqualValues("error description", e.ErrorName)
	s.Assert().EqualValues(400, e.Status)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

const (
	testApiErrorData = `{
	"status": 400,
	"error": "error description",
	"request_id": 200132,
	"data": []
}`
)
