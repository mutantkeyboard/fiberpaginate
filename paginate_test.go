package fiberpaginate

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Response struct {
	Page            int         `json:"page"`
	Limit           int         `json:"limit"`
	Offset          int         `json:"offset"`
	Start           int         `json:"start"`
	Sort            []SortField `json:"sort"`
	NextPageURL     string      `json:"next_PageURL"`
	PreviousPageURL string      `json:"prev_PageURL"`
}

// go test -run Test_PaginateWithQueries
func Test_PaginateWithQueries(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New(Config{
		DefaultSort: "id",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:            pageInfo.Page,
			Limit:           pageInfo.Limit,
			Offset:          pageInfo.Offset,
			Start:           pageInfo.Start(),
			Sort:            pageInfo.Sort,
			NextPageURL:     pageInfo.NextPageURL(c.BaseURL()),
			PreviousPageURL: pageInfo.PreviousPageURL(c.BaseURL()),
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?page=2&limit=20", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 2, respBody.Page)
	utils.AssertEqual(t, 20, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 20, respBody.Start)
	utils.AssertEqual(t, "http://example.com?page=3&limit=20", respBody.NextPageURL)
	utils.AssertEqual(t, "http://example.com?page=1&limit=20", respBody.PreviousPageURL)
	utils.AssertEqual(t, []SortField{{Field: "id", Order: ASC}}, respBody.Sort)
}

// go test -run TestPreviousAndNextPageURLMethods
func TestPreviousAndNextPageURLMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		pageInfo     PageInfo
		baseURL      string
		expectedNext string
		expectedPrev string
	}{
		{
			name:         "Middle page",
			pageInfo:     PageInfo{Page: 2, Limit: 10},
			baseURL:      "https://example.com/users",
			expectedNext: "https://example.com/users?page=3&limit=10",
			expectedPrev: "https://example.com/users?page=1&limit=10",
		},
		{
			name:         "First page",
			pageInfo:     PageInfo{Page: 1, Limit: 20},
			baseURL:      "https://example.com/users",
			expectedNext: "https://example.com/users?page=2&limit=20",
			expectedPrev: "",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextURL := tt.pageInfo.NextPageURL(tt.baseURL)
			prevURL := tt.pageInfo.PreviousPageURL(tt.baseURL)

			utils.AssertEqual(t, tt.expectedNext, nextURL)
			utils.AssertEqual(t, tt.expectedPrev, prevURL)
		})
	}
}

func Test_PaginateWithOffset(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?offset=20&limit=20", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 1, respBody.Page)
	utils.AssertEqual(t, 20, respBody.Limit)
	utils.AssertEqual(t, 20, respBody.Offset)
	utils.AssertEqual(t, 20, respBody.Start) // This should be 20, matching the current Start() implementation
	utils.AssertEqual(t, []SortField{{Field: "id", Order: ASC}}, respBody.Sort)
}

// go test -run Test_PaginateCheckDefaultsWhenNoQueries
func Test_PaginateCheckDefaultsWhenNoQueries(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 1, respBody.Page)
	utils.AssertEqual(t, 10, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 0, respBody.Start)
	utils.AssertEqual(t, []SortField{{Field: "id", Order: ASC}}, respBody.Sort)
}

// go test -run Test_PaginateCheckDefaultsWhenNoPage
func Test_PaginateCheckDefaultsWhenNoPage(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?limit=20", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 1, respBody.Page)
	utils.AssertEqual(t, 20, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 0, respBody.Start)
	utils.AssertEqual(t, []SortField{{Field: "id", Order: ASC}}, respBody.Sort)
}

// go test -run Test_PaginateCheckDefaultsWhenNoLimit
func Test_PaginateCheckDefaultsWhenNoLimit(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New())

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?page=2", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 2, respBody.Page)
	utils.AssertEqual(t, 10, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 10, respBody.Start)
	utils.AssertEqual(t, []SortField{{Field: "id", Order: ASC}}, respBody.Sort)
}

// go test -run Test_PaginateConfigDefaultPageDefaultLimit
func Test_PaginateConfigDefaultPageDefaultLimit(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		DefaultPage:  100,
		DefaultLimit: MaxLimit,
		DefaultSort:  "name",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 100, respBody.Page)
	utils.AssertEqual(t, MaxLimit, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 9900, respBody.Start)
	utils.AssertEqual(t, []SortField{{Field: "name", Order: ASC}}, respBody.Sort)
}

// go test -run Test_PaginateConfigPageKeyLimitKey
// go test -run Test_PaginateConfigPageKeyLimitKey
func Test_PaginateConfigPageKeyLimitKey(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		PageKey:     "site",
		LimitKey:    "size",
		DefaultSort: "id",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?site=2&size=5", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 2, respBody.Page)
	utils.AssertEqual(t, 5, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 5, respBody.Start)
	utils.AssertEqual(t, []SortField{{Field: "id", Order: ASC}}, respBody.Sort)
}

// go test -run Test_PaginateNegativeDefaultPageDefaultLimitValues
func Test_PaginateNegativeDefaultPageDefaultLimitValues(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		DefaultPage:  -1,
		DefaultLimit: -1,
		DefaultSort:  "id",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	utils.AssertEqual(t, nil, json.Unmarshal(bodyBytes, &respBody))

	utils.AssertEqual(t, 1, respBody.Page)
	utils.AssertEqual(t, 10, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 0, respBody.Start)
	utils.AssertEqual(t, []SortField{{Field: "id", Order: ASC}}, respBody.Sort)
}

// go test -run Test_PaginateFromContextWithoutNew
func Test_PaginateFromContextWithoutNew(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

	body := resp.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	utils.AssertEqual(t, nil, err)

	var respBody Response
	err = json.Unmarshal(bodyBytes, &respBody)
	utils.AssertEqual(t, true, err != nil) // Expecting an error because the response should be empty

	// Assert that all fields are zero values
	utils.AssertEqual(t, 0, respBody.Page)
	utils.AssertEqual(t, 0, respBody.Limit)
	utils.AssertEqual(t, 0, respBody.Offset)
	utils.AssertEqual(t, 0, respBody.Start)
	utils.AssertEqual(t, []SortField(nil), respBody.Sort)
}

// go test -run Test_PaginateWithMultipleSorting
func Test_PaginateWithMultipleSorting(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(New(Config{
		SortKey:      "sort",
		DefaultSort:  "id",
		AllowedSorts: []string{"id", "name", "date"},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(Response{
			Page:   pageInfo.Page,
			Limit:  pageInfo.Limit,
			Offset: pageInfo.Offset,
			Start:  pageInfo.Start(),
			Sort:   pageInfo.Sort,
		})
	})

	testCases := []struct {
		name           string
		url            string
		expectedSort   []SortField
		expectedStatus int
	}{
		{"Default Sort", "/", []SortField{{Field: "id", Order: ASC}}, 200},
		{"Single Sort", "/?sort=name", []SortField{{Field: "name", Order: ASC}}, 200},
		{"Multiple Sort", "/?sort=name,-date", []SortField{{Field: "name", Order: ASC}, {Field: "date", Order: DESC}}, 200},
		{"Invalid Sort", "/?sort=invalid", []SortField{{Field: "id", Order: ASC}}, 200},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, tc.url, nil))
			utils.AssertEqual(t, nil, err)
			utils.AssertEqual(t, tc.expectedStatus, resp.StatusCode)

			var result Response
			err = json.NewDecoder(resp.Body).Decode(&result)
			utils.AssertEqual(t, nil, err)

			utils.AssertEqual(t, tc.expectedSort, result.Sort)
		})
	}
}

// go test -run TestParseSortQuery
func TestParseSortQuery(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		allowedSorts []string
		defaultSort  string
		expected     []SortField
	}{
		{
			name:         "Empty query",
			query:        "",
			allowedSorts: []string{"id", "name", "date"},
			defaultSort:  "id",
			expected:     []SortField{{Field: "id", Order: ASC}},
		},
		{
			name:         "Single allowed field",
			query:        "name",
			allowedSorts: []string{"id", "name", "date"},
			defaultSort:  "id",
			expected:     []SortField{{Field: "name", Order: ASC}},
		},
		{
			name:         "Multiple fields with mixed order",
			query:        "name,-date,id",
			allowedSorts: []string{"id", "name", "date"},
			defaultSort:  "id",
			expected: []SortField{
				{Field: "name", Order: ASC},
				{Field: "date", Order: DESC},
				{Field: "id", Order: ASC},
			},
		},
		{
			name:         "Disallowed field",
			query:        "email,name",
			allowedSorts: []string{"id", "name", "date"},
			defaultSort:  "id",
			expected:     []SortField{{Field: "name", Order: ASC}},
		},
		{
			name:         "All disallowed fields",
			query:        "email,phone",
			allowedSorts: []string{"id", "name", "date"},
			defaultSort:  "id",
			expected:     []SortField{{Field: "id", Order: ASC}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSortQuery(tt.query, tt.allowedSorts, tt.defaultSort)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseSortQuery() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_PaginateEdgeCases tests edge cases for the pagination middleware.
//
// Specifically, it tests that the middleware correctly handles:
//   - Negative page numbers
//   - Page numbers of 0
//   - Negative limit numbers
//   - Limits of 0
//
// In each of these cases, the middleware should return a PageInfo with a page
// number of 1 and a limit of 10 (the default limit).
func Test_PaginateEdgeCases(t *testing.T) {
	t.Parallel()
	app := fiber.New()

	app.Use(New(Config{
		DefaultSort:  "id",
		DefaultLimit: 10, // Explicitly set the default limit to 10
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}
		return c.JSON(pageInfo)
	})

	testCases := []struct {
		name          string
		url           string
		expectedPage  int
		expectedLimit int
	}{
		{"Negative page", "/?page=-1", 1, 10},
		{"Page zero", "/?page=0", 1, 10},
		{"Negative limit", "/?limit=-10", 1, 10},
		{"Limit zero", "/?limit=0", 1, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := app.Test(httptest.NewRequest("GET", tc.url, nil))
			utils.AssertEqual(t, nil, err)
			utils.AssertEqual(t, 200, resp.StatusCode)

			var result PageInfo
			utils.AssertEqual(t, nil, json.NewDecoder(resp.Body).Decode(&result))
			utils.AssertEqual(t, tc.expectedPage, result.Page)
			utils.AssertEqual(t, tc.expectedLimit, result.Limit)
		})
	}
}

/* BENCHMARK TESTING */

func BenchmarkPaginateMiddleware(b *testing.B) {
	app := fiber.New()
	app.Use(New())

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, _ := FromContext(c)
		return c.JSON(pageInfo)
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/?page=2&limit=20&sort=name,-date", nil)
		_, err := app.Test(req, -1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPaginateMiddlewareWithCustomConfig(b *testing.B) {
	app := fiber.New()
	app.Use(New(Config{
		PageKey:      "p",
		LimitKey:     "l",
		SortKey:      "s",
		DefaultPage:  1,
		DefaultLimit: 30,
		DefaultSort:  "id",
		AllowedSorts: []string{"id", "name", "date"},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, _ := FromContext(c)
		return c.JSON(pageInfo)
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/?p=3&l=25&s=name,-id", nil)
		_, err := app.Test(req, -1)
		if err != nil {
			b.Fatal(err)
		}
	}
}
