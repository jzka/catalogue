package catalogue

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type getListRequest struct {
	Categories []string `json:"categories"`
	OrderBy    string   `json:"order"`
	PageSize   int      `json:"page_size"`
	PageNum    int      `json:"page"`
}

type getListResponse struct {
	Products []Product `json:"products"`
	Err      error     `json:"err,omitempty"`
}

type countProductRequest struct {
	Categories []string `json:"categories"`
}
type countProductResponse struct {
	Count int   `json:"count"`
	Err   error `json:"err,omitempty"`
}

type getProductRequest struct {
	ID string `json:"id"`
}
type getProductResponse struct {
	Product Product `json:"product"`
	Err     error   `json:"err,omitempty"`
}

type getCategoriesRequest struct {
}

type getCategoriesResponse struct {
	Categories []string `json:"categories"`
	Err        error    `json:"err,omitempty"`
}

//endp wrapper
type Endpoints struct {
	GetListEndpoint       endpoint.Endpoint
	CountProductsEndpoint endpoint.Endpoint
	GetProductEndpoint    endpoint.Endpoint
	GetCategoriesEndpoint endpoint.Endpoint
}

func MakeGetListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getListRequest)
		prods, err := svc.GetList(req.Categories, req.OrderBy, req.PageNum, req.PageSize)
		return getListResponse{Products: prods, Err: err}, err
	}
}

func MakeCountProductsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countProductRequest)
		count, err := svc.CountProducts(req.Categories)
		return countProductResponse{Count: count, Err: err}, err
	}
}

func MakeGetProductEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getProductRequest)
		prod, err := svc.GetProduct(req.ID)
		return getProductResponse{Product: prod, Err: err}, err
	}
}

func MakeGetCategoriesEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cat, err := svc.GetCategories()
		return getCategoriesResponse{Categories: cat, Err: err}, err
	}
}
