package validators

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateGetAllPaginationParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("default_values", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/", nil)

		limit, skip, isValid := ValidateGetAllPaginationParams(ctx)

		assert.True(t, isValid)
		assert.Equal(t, 10, limit) // Default limit
		assert.Equal(t, 0, skip)   // Default skip
	})

	t.Run("valid_parameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=20&skip=5", nil)

		limit, skip, isValid := ValidateGetAllPaginationParams(ctx)

		assert.True(t, isValid)
		assert.Equal(t, 20, limit)
		assert.Equal(t, 5, skip)
	})

	t.Run("zero_values", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=0&skip=0", nil)

		limit, skip, isValid := ValidateGetAllPaginationParams(ctx)

		assert.True(t, isValid)
		assert.Equal(t, 0, limit)
		assert.Equal(t, 0, skip)
	})

	t.Run("negative_values", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=-10&skip=-5", nil)

		limit, skip, isValid := ValidateGetAllPaginationParams(ctx)

		assert.True(t, isValid)
		assert.Equal(t, -10, limit) // Negative values are allowed by the validator
		assert.Equal(t, -5, skip)   // Service layer will handle normalization
	})

	t.Run("invalid_limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=abc&skip=5", nil)

		_, _, isValid := ValidateGetAllPaginationParams(ctx)

		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid limit parameter", response["error"])
		assert.Equal(t, "Limit must be a valid integer", response["details"])
	})

	t.Run("invalid_skip", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=10&skip=xyz", nil)

		_, _, isValid := ValidateGetAllPaginationParams(ctx)

		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid skip parameter", response["error"])
		assert.Equal(t, "Skip must be a valid integer", response["details"])
	})

	t.Run("both_parameters_invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=abc&skip=xyz", nil)

		_, _, isValid := ValidateGetAllPaginationParams(ctx)

		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid limit parameter", response["error"])
	})

	t.Run("large_integer_values", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=1000000&skip=500000", nil)

		limit, skip, isValid := ValidateGetAllPaginationParams(ctx)

		assert.True(t, isValid)
		assert.Equal(t, 1000000, limit) // The validator allows any integer
		assert.Equal(t, 500000, skip)   // Service layer will handle normalization if needed
	})

	t.Run("float_values", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=10.5&skip=5.2", nil)

		_, _, isValid := ValidateGetAllPaginationParams(ctx)

		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid limit parameter", response["error"])
	})

	t.Run("empty_parameter_values", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("GET", "/?limit=&skip=", nil)

		_, _, isValid := ValidateGetAllPaginationParams(ctx)

		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid limit parameter", response["error"])
	})
}
