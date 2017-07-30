package context

import (
	"fmt"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	FAKE_RESPONSE_TEXT = "Hi there"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ContextSuite struct{}

var _ = Suite(&ContextSuite{})

func (s *ContextSuite) TestExecuteHttp(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, FAKE_RESPONSE_TEXT)
	}))
	defer ts.Close()

	mockReq := newMockRequest(c, ts.URL)
	makeTestCall(c, mockReq)
}

func newMockRequest(c *C, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Fatal(err)
	}
	return req
}

func makeTestCall(c *C, req *http.Request) {
	client := NewContextHttpClient(nil, 0)
	ctx := NewCustomContext()

	if err := client.ExecuteHttp(ctx, req, func(response *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer response.Body.Close()

		c.Assert(response.StatusCode, Equals, http.StatusOK)

		buffer, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		text := strings.TrimSpace(string(buffer))
		c.Assert(text, Equals, FAKE_RESPONSE_TEXT)

		return nil
	}); err != nil {
		c.Fatal(err)
	}
}
