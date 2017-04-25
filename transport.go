package catalogue

//TODO: ADD PROMETHEUS + LOGING MIDLE
import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

var (
	ErrBadRouting = "bad routing error"
)

func MakeHttpHandler(ctx context.Context, endpoint Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/catalogue").Handler(httptransport.NewServer(
		endpoint.GetListEndpoint,
		decodeGetListRequest,
		encodeGetListResponse,
		options...,
	))
	r.Methods("GET").Path("/catalogue/{id}").Handler(httptransport.NewServer(
		endpoint.GetProductEndpoint,
		decodeGetProductRequest,
		encodeGetProductResponse,
		options...,
	))
	r.Methods("GET").Path("/tags").Handler(httptransport.NewServer(
		endpoint.CountProductsEndpoint,
		decodeGetCategoriesRequest,
		encodeResponse,
		options...,
	))
	return r
}

// Encode Functions

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil err")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func encodeGetListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(getListResponse)
	if resp.Err != nil {
		encodeError(ctx, resp.Err, w)
		return nil
	}
	return encodeResponse(ctx, w, resp.Products)
}

func encodeGetProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(getProductResponse)
	if resp.Err != nil {
		encodeError(ctx, resp.Err, w)
		return nil
	}
	return encodeResponse(ctx, w, resp.Product)
}

//Decode Functions

func decodeGetListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	pageNum := 1
	pageSize := 5
	order := "id"
	categories := []string{}

	if page := r.FormValue("page"); page != "" {
		pageNum, _ = strconv.Atoi(page)
	}
	if size := r.FormValue("page_size"); size != "" {
		pageSize, _ = strconv.Atoi(size)
	}
	if orderby := r.FormValue("order"); orderby != "" {
		order = strings.ToLower(orderby)
	}
	if cat := r.FormValue("categories"); cat != "" {
		categories = strings.Split(cat, ",")
	}
	return getListRequest{
		Categories: categories,
		OrderBy:    order,
		PageSize:   pageSize,
		PageNum:    pageNum,
	}, nil
}

func decodeGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	return getProductRequest{
		ID: id,
	}, nil
}

func decodeGetCategoriesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func decodeCountProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	cat := []string{}
	if catstr := r.FormValue("categories"); catstr != "" {
		cat = strings.Split(catstr, ",")
	}
	return countProductRequest{
		Categories: cat,
	}, nil
}
