package createurl_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/handlers/createurl"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/handlers/createurl/mocks"
)

type mocksToUse struct {
	creator *mocks.AliasCreator
}

func initMocks(t *testing.T) *mocksToUse {
	return &mocksToUse{
		creator: mocks.NewAliasCreator(t),
	}
}

func TestCrateAlias(t *testing.T) {
	defaultURL := "https://test.com"
	defaultAlias := "testAlias"

	testCases := []struct {
		name          string
		storageOutput func(m *mocksToUse)
		expectedCode  int
		expectedResp  string
		reqBody       string
		reqMethod     string
		reqTarget     string
	}{
		{
			name: "CrateAliasSuccess",
			storageOutput: func(m *mocksToUse) {
				m.creator.On("CreateAlias", mock.Anything, defaultURL).Return(defaultAlias, nil).Once()
			},
			expectedCode: http.StatusCreated,
			expectedResp: `{"alias":"testAlias"}`,
			reqBody:      `{"original":"https://test.com"}`,
			reqMethod:    "POST",
			reqTarget:    "/alias",
		},
		{
			name: "CrateAliasFail",
			storageOutput: func(m *mocksToUse) {
				m.creator.On("CreateAlias", mock.Anything, defaultURL).Return("", errors.New("any err")).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"error":"alias creation fail"}`,
			reqBody:      `{"original":"https://test.com"}`,
			reqMethod:    "POST",
			reqTarget:    "/alias",
		},
		{
			name:          "CrateAliasInvalidJSON",
			storageOutput: func(_ *mocksToUse) {},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedResp:  `{"error":"invalid JSON"}`,
			reqBody:       `{"fail":"https://test.com"}`,
			reqMethod:     "POST",
			reqTarget:     "/alias",
		},
		{
			name:          "CrateAliasBadJsonStructureFail",
			storageOutput: func(_ *mocksToUse) {},
			expectedCode:  http.StatusBadRequest,
			expectedResp:  `{"error":"bad JSON"}`,
			reqBody:       `{"fail"}`,
			reqMethod:     "POST",
			reqTarget:     "/alias",
		},
	}
	log := zerolog.New(os.Stderr)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resourceMock := initMocks(t)
			tc.storageOutput(resourceMock)

			handler := createurl.New(log, resourceMock.creator)
			req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, strings.NewReader(tc.reqBody))
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			body, err := io.ReadAll(w.Body)
			bodyStr := strings.TrimSpace(string(body))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedResp, bodyStr)
		})
	}
}
