package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/product"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createTestServer() (*gin.Engine, error) {
	gin.SetMode(gin.TestMode)

	err := os.Setenv("TOKEN", "MITOKEN123")
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	storage := store.NewJSONstore("../../../test_products.json")
	repository := product.NewDefaultRepository(storage)
	service := product.NewDefaultService(repository)
	handler := NewProductHandler(service)

	group := router.Group("/products")
	{
		group.GET("/", handler.GetAll())
		group.POST("/", handler.Create())
		group.GET("/:id", handler.GetByID())
		group.GET("/search", handler.GetBySearch())
		group.PUT("/:id", handler.Update())
		group.PATCH("/:id", handler.UpdatePartial())
		group.DELETE("/:id", handler.Delete())
		group.GET("/consumer_price", handler.GetConsumerPrice())
	}

	return router, nil
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("TOKEN", "MITOKEN123")

	return req, httptest.NewRecorder()
}

func loadProducts(path string) ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func writeProducts(path string, list []domain.Product) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return err
}

func TestGetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}

		expectedBody := `{
			"data":[
				{"id":1,"name":"Oil - Margarine","quantity":439,"code_value":"S82254D","is_published":true,"expiration":"15/12/2021","price":71.42},
				{"id":2,"name":"Pineapple - Canned, Rings","quantity":345,"code_value":"M4637","is_published":true,"expiration":"09/08/2021","price":352.79},
				{"id":3,"name":"Wine - Red Oakridge Merlot","quantity":367,"code_value":"T65812","is_published":false,"expiration":"24/05/2021","price":179.23}
			]
		}`

		//Arrange
		//Crear server y definir rutas
		server, err := createTestServer()
		if err != nil {
			panic(err)
		}

		//Act
		request, respRecorder := createRequestTest(http.MethodGet, "/products/", "")
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, 200, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
	})
}

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}

		expectedBody := `{
			"data":{"id":2,"name":"Pineapple - Canned, Rings","quantity":345,"code_value":"M4637","is_published":true,"expiration":"09/08/2021","price":352.79}
		}`

		//Arrange
		//Crear server y definir rutas
		server, err := createTestServer()
		if err != nil {
			panic(err)
		}

		//Act
		request, respRecorder := createRequestTest(http.MethodGet, "/products/2", "")
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, 200, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
	})

	t.Run("Bad Request", func(t *testing.T) {
		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}

		expectedBody := `{
			"status":400,
			"code":"RequestError",
			"message":"The ID is not valid"
		}`

		//Arrange
		//Crear server y definir rutas
		server, err := createTestServer()
		if err != nil {
			panic(err)
		}

		//Act
		request, respRecorder := createRequestTest(http.MethodGet, "/products/a4", "")
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, 400, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
	})

	t.Run("Not Found", func(t *testing.T) {
		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}

		expectedBody := `{
			"status":404,
			"code":"NotFoundError",
			"message":"The requested resource was not found"
		}`

		//Arrange
		//Crear server y definir rutas
		server, err := createTestServer()
		if err != nil {
			panic(err)
		}

		//Act
		request, respRecorder := createRequestTest(http.MethodGet, "/products/1000", "")
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, 404, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
	})
}

func TestCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}

		expectedBody := `{
			"data":{"id":4,"name":"Pelota","quantity":12,"code_value":"PELO1","is_published":true,"expiration":"09/08/2023","price":35}
		}`

		//Arrange
		//Crear server y definir rutas
		server, err := createTestServer()
		if err != nil {
			panic(err)
		}

		//Cargar productos antes de la prueba
		defaultProducts, _ := loadProducts("../../../test_products.json")

		//Act
		request, respRecorder := createRequestTest(http.MethodPost, "/products/", `{
			"name":"Pelota","quantity":12,"code_value":"PELO1","is_published":true,"expiration":"09/08/2023","price":35
		}`)
		server.ServeHTTP(respRecorder, request)

		//Guardar y dejar los productos como estaban antes de la prueba
		writeProducts("../../../test_products.json", defaultProducts)

		//Assert
		assert.Equal(t, 201, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
	})

	t.Run("Unauthorized", func(t *testing.T) {
		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}

		expectedBody := `{
			"status":401,
			"code":"InvalidTokenError",
			"message":"The user token is not valid"
		}`

		//Arrange
		//Crear server y definir rutas
		server, err := createTestServer()
		if err != nil {
			panic(err)
		}

		//Act
		request, respRecorder := createRequestTest(http.MethodPost, "/products/", `{
			"name":"Pelota","quantity":12,"code_value":"PELO1","is_published":true,"expiration":"09/08/2023","price":35
		}`)
		request.Header.Set("TOKEN", "abc123")
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, 401, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
	})
}

func TestDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
			"Location":     []string{"/products/2"},
		}

		//Arrange
		//Crear server y definir rutas
		server, err := createTestServer()
		if err != nil {
			panic(err)
		}

		//Cargar productos antes de la prueba
		defaultProducts, _ := loadProducts("../../../test_products.json")

		//Act
		request, respRecorder := createRequestTest(http.MethodDelete, "/products/2", "")
		server.ServeHTTP(respRecorder, request)

		//Guardar y dejar los productos como estaban antes de la prueba
		writeProducts("../../../test_products.json", defaultProducts)

		//Assert
		assert.Equal(t, 204, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.Equal(t, "", respRecorder.Body.String())
	})
}
