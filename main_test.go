package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//"golang.org/x/text/message"
)

var city = "moscow"

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET",
			fmt.Sprintf("/cafe?city=%s&count=%d", city, v.count),
			nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		body := strings.TrimSpace(response.Body.String())

		var act int

		if body == "" {
			act = 0
		} else {
			act = len(strings.Split(body, ","))
		}

		assert.Equal(t, v.want, act)
	}
}

func TestCafeSearch(t *testing.T) {
	requests := []struct {
		search string
		want   int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		handler := http.HandlerFunc(mainHandle)
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET",
			fmt.Sprintf("/cafe?city=%s&search=%s", city, v.search),
			nil)
		handler.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		body := strings.TrimSpace(resp.Body.String())

		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}
		assert.Equal(t, v.want, len(cafes))

		search := strings.ToLower(v.search)

		for _, cafe := range cafes {
			cafeLower := strings.ToLower(cafe)
			assert.True(t, strings.Contains(cafeLower, search))
		}
	}
}

func TestCafeWhenNot(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?city=moscow",
		"/cafe?city=tula&count=2",
		"/cafe?city=moscow&search=кофе",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)
		handler.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}

}
