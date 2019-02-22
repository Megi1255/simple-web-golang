package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// 입력으로 사용할 JSON body
	testdata := []struct {
		input    string
		expected map[string]interface{}
	}{
		{
			`{
    "api_name": "CreateUser",
    "user_name": "gin"
}`,
			map[string]interface{}{
				"result":     "OK",
				"api_name":   "CreateUser",
				"user_name":  "gin",
				"user_score": float64(100),
			},
		},
		{
			`{
    "api_name": "GetUserScore",
    "user_name": "gin"
}`,
			map[string]interface{}{
				"result":     "OK",
				"api_name":   "GetUserScore",
				"user_name":  "gin",
				"user_score": float64(100),
			},
		},
	}
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()

	for _, d := range testdata {
		// 테스트 http 요청 생성
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/gateway", strings.NewReader(d.input))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		// 검사
		var got map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Errorf("Unable to parse JSON body '%s', err: %v", w.Body, err)
		}
		if !reflect.DeepEqual(got, d.expected) {
			t.Errorf("api_name got: %v wanted: %v", got, d.expected)
		}
	}
}

func BenchmarkGateway(b *testing.B) {
	d := `{
    "api_name": "GetUserScore",
    "user_name": "gin"
}`
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/gateway", strings.NewReader(d))
	req.Header.Set("Content-Type", "application/json")

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
}
