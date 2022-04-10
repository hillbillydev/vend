package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vend/vend"
	"github.com/vend/vend/inmemory"
)

func TestGetProducts(t *testing.T) {
	t.Run("No items should return 404.", func(t *testing.T) {
		server, _ := setupTest()

		req, err := http.NewRequest("GET", "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusNotFound)
		}
	})

	t.Run("Should find two items.", func(t *testing.T) {
		server, storage := setupTest()
		storage.CreateProduct("Club", 100)
		storage.CreateProduct("Bike", 1000)

		req, err := http.NewRequest("GET", "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}

		var res []vend.Product
		json.NewDecoder(rr.Body).Decode(&res)

		if len(res) != 2 {
			t.Errorf("handler did not return 2 items, but instead returned %d items.", len(res))
		}
	})

}

func setupTest() (*Server, *inmemory.InMemory) {
	storage := inmemory.NewInMemoryInMemory()
	s := NewServer("8080", storage)

	return s, storage
}
