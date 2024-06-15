package origurl_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/handlers/origurl"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/handlers/origurl/mocks"
)

type mocksToUse struct {
	provider *mocks.OriginalProvider
}

func initMocks(t *testing.T) *mocksToUse {
	return &mocksToUse{
		provider: mocks.NewOriginalProvider(t),
	}
}

func TestGetAlias(t *testing.T) {
	defaultURL := "http://test.com"
	defaultAlias := "testAlias"

	testCases := []struct {
		name          string
		storageOutput func(m *mocksToUse)
		expectedCode  int
		expectedResp  string
		urlParamID    string
		reqMethod     string
		reqTarget     string
	}{
		{
			name: "GetAliasSuccess",
			storageOutput: func(m *mocksToUse) {
				m.provider.On("GetOrigURL", mock.Anything, defaultAlias).Return(defaultURL, nil).Once()
			},
			expectedCode: http.StatusFound,
			expectedResp: `<a href="https://test.com">Found</a>.`,
			urlParamID:   defaultAlias,
			reqMethod:    "GET",
			reqTarget:    "/",
		},
		{
			name: "GetAliasFail",
			storageOutput: func(m *mocksToUse) {
				m.provider.On("GetOrigURL", mock.Anything, defaultAlias).Return("", errors.New("any err")).Once()
			},
			expectedCode: http.StatusNotFound,
			expectedResp: `{"error":"alias search fail"}`,
			urlParamID:   defaultAlias,
			reqMethod:    "GET",
			reqTarget:    "/",
		},
	}
	log := zerolog.New(os.Stderr)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resourceMock := initMocks(t)
			tc.storageOutput(resourceMock)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("alias", tc.urlParamID)

			handler := origurl.New(log, resourceMock.provider)

			req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

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
