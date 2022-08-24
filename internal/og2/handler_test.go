package og2_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"hunter.io/og2/internal/og2"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	s1 := og2.Session{User: og2.User{Name: "John Doe"}}

	testCases := []struct {
		name     string
		input    og2.User
		expected og2.Session
	}{
		{
			name:     "test create user",
			input:    s1.User,
			expected: s1,
		},
	}

	for _, testCase := range testCases {
		require := require.New(t)

		b, err := json.Marshal(testCase.input)
		require.NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(b))
		w := httptest.NewRecorder()

		sessions := og2.NewSessions()
		h := og2.NewHandler(sessions).HandleUser()
		h.ServeHTTP(w, req)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		require.NoError(err)

		fmt.Println(string(body))

		var actual og2.Session
		err = json.Unmarshal(body, &actual)
		require.NoError(err)

		require.Equal(testCase.expected, actual)
	}
}
