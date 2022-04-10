package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/vend/vend"
)

type Server struct {
	r       chi.Router
	port    string
	storage vend.Storage
}

func NewServer(port string, storage vend.Storage) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	s := &Server{
		r:       r,
		storage: storage,
		port:    port,
	}

	r.Get("/products", s.getProducts)
	r.Post("/products", s.postProducts)
	r.Post("/sales", s.postSales)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	port := fmt.Sprintf(":%s", s.port)
	return http.ListenAndServe(port, s.r)
}

func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.storage.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(products) == 0 {
		http.Error(w, "No products", http.StatusNotFound)
		return
	}

	bs, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bs)
}

type postProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (req postProductRequest) IsValid() error {
	if len(req.Name) == 0 {
		return errors.New("name length can't be 0")
	}

	if req.Price < 0 {
		return fmt.Errorf("expected price to be >= 0 but it was %d", req.Price)
	}

	return nil
}

func (s *Server) postProducts(w http.ResponseWriter, r *http.Request) {
	var req postProductRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = req.IsValid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := s.storage.CreateProduct(req.Name, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bs)
}

type (
	postSalesRequest struct {
		DollarDiscount int `json:"dollarDiscount"`
		Sales          []struct {
			ProductID string `json:"productId"`
			Quantity  int    `json:"quantity"`
		} `json:"sales"`
	}

	postSalesResponse struct {
		Products   []salesProductResponse `json:"products"`
		TotalPrice int                    `json:"totalPrice"`
	}

	salesProductResponse struct {
		vend.Product
		Quantity   int `json:"quantity"`
		TotalPrice int `json:"totalPrice"`
	}
)

func (req postSalesRequest) IsValid() error {
	for i, s := range req.Sales {
		_, err := uuid.Parse(s.ProductID)
		if err != nil {
			return fmt.Errorf("not a valid id in item with index of [%d]. %w", i, err)
		}

		if s.Quantity < 1 {
			return fmt.Errorf("not a valid quantity in item with index of [%d], expected it to be more then 0 but got %d.", i, s.Quantity)
		}
	}

	return nil
}

func productsToSalesProdcutResponse(products []vend.Product) []salesProductResponse {
	var result []salesProductResponse
	counter := make(map[string]salesProductResponse)

	for _, p := range products {
		res, ok := counter[p.ID.String()]
		if !ok {
			counter[p.ID.String()] = salesProductResponse{
				Product:    p,
				Quantity:   1,
				TotalPrice: p.Price,
			}
			continue
		}

		res.Quantity++
		res.TotalPrice += p.Price
		counter[p.ID.String()] = res
	}

	for _, sale := range counter {
		result = append(result, sale)
	}

	return result
}

func (s *Server) postSales(w http.ResponseWriter, r *http.Request) {
	var req postSalesRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = req.IsValid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sales []vend.LineItem
	for _, s := range req.Sales {
		sales = append(sales, vend.LineItem{
			ProductID: uuid.MustParse(s.ProductID),
			Quantity:  s.Quantity,
		})

	}

	products, err := s.storage.Sale(sales)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	products = vend.ApplyDiscountOnProducts(products, req.DollarDiscount)

	res := postSalesResponse{
		Products:   productsToSalesProdcutResponse(products),
		TotalPrice: vend.TotalPrice(products),
	}

	bs, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bs)
}
