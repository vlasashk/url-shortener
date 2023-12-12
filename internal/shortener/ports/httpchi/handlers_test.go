package httpchi_test

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/vlasashk/url-shortener/internal/shortener/mocks"
	"github.com/vlasashk/url-shortener/internal/shortener/models/logger"
	"github.com/vlasashk/url-shortener/internal/shortener/models/service"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type UnitTestSuite struct {
	suite.Suite
	service   service.Service
	testAlias string
	testUrl   string
}
type TestCase struct {
	testName      string
	storageOutput func()
	expectedCode  int
	expectedResp  string
	urlParamID    string
	reqBody       string
	reqMethod     string
	reqTarget     string
}

func (suite *UnitTestSuite) SetupTest() {
	suite.service = mocks.NewService(suite.T())
	suite.testAlias = "testAlias"
	suite.testUrl = "https://test.com"
}

func (suite *UnitTestSuite) TestCrateAlias() {
	testCases := []TestCase{
		{
			storageOutput: func() {
				suite.service.(*mocks.Service).On("CrateAlias", suite.testUrl).Return(suite.testAlias, nil).Once()
			},
			expectedCode: http.StatusCreated,
			expectedResp: `{"alias":"testAlias"}`,
			reqBody:      `{"original":"https://test.com"}`,
			reqMethod:    "POST",
			reqTarget:    "/alias",
		},
		{
			storageOutput: func() {
				suite.service.(*mocks.Service).On("CrateAlias", suite.testUrl).Return("", errors.New("any err")).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"error":"alias creation fail"}`,
			reqBody:      `{"original":"https://test.com"}`,
			reqMethod:    "POST",
			reqTarget:    "/alias",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedResp:  `{"error":"invalid JSON"}`,
			reqBody:       `{"fail":"https://test.com"}`,
			reqMethod:     "POST",
			reqTarget:     "/alias",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusBadRequest,
			expectedResp:  `{"error":"bad JSON"}`,
			reqBody:       `{"fail"}`,
			reqMethod:     "POST",
			reqTarget:     "/alias",
		},
	}
	log := logger.New(zerolog.InfoLevel)
	for _, tc := range testCases {
		tc.storageOutput()
		handler := httpchi.NewHandler(suite.service, log)
		req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, strings.NewReader(tc.reqBody))
		w := httptest.NewRecorder()

		handler.CrateAlias(w, req)

		body, err := io.ReadAll(w.Body)
		bodyStr := strings.TrimSpace(string(body))
		suite.NoError(err)
		suite.Equal(tc.expectedCode, w.Code)
		suite.Equal(tc.expectedResp, bodyStr)
	}
}

func (suite *UnitTestSuite) TestGetOrigURL() {
	testCases := []TestCase{
		{
			storageOutput: func() {
				suite.service.(*mocks.Service).On("GetOrigURL", suite.testAlias).Return(suite.testUrl, nil).Once()
			},
			expectedCode: http.StatusFound,
			expectedResp: `<a href="https://test.com">Found</a>.`,
			urlParamID:   suite.testAlias,
			reqMethod:    "GET",
			reqTarget:    "/",
		},
		{
			storageOutput: func() {
				suite.service.(*mocks.Service).On("GetOrigURL", suite.testAlias).Return("", errors.New("any err")).Once()
			},
			expectedCode: http.StatusNotFound,
			expectedResp: `{"error":"alias search fail"}`,
			urlParamID:   suite.testAlias,
			reqMethod:    "GET",
			reqTarget:    "/",
		},
	}
	log := logger.New(zerolog.InfoLevel)
	for _, tc := range testCases {
		tc.storageOutput()
		handler := httpchi.NewHandler(suite.service, log)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("alias", tc.urlParamID)
		req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, strings.NewReader(tc.reqBody))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()

		handler.GetOrigURL(w, req)

		body, err := io.ReadAll(w.Body)
		bodyStr := strings.TrimSpace(string(body))
		suite.NoError(err)
		suite.Equal(tc.expectedCode, w.Code)
		suite.Equal(tc.expectedResp, bodyStr)
	}
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
