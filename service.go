package catalogue

import (
	"errors"
	"strings"

	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

// Service Catalogue - bussiness logic
type Service interface {
	GetList(categories []string, orderBy string, currentPageNum, currentPageSize int) ([]Product, error)
	CountProducts(categories []string) (int, error)
	GetProduct(id string) (Product, error)
	GetCategories() ([]string, error)
	//HealthCheck() []Health
	//List(tags []string, order string, pageNum, pageSize int) ([]Sock, error) // GET /catalogue

}

// Middleware is a decorator for our Service
type Middleware func(Service) Service

// catalogueService contains a connection to sql database,
// and a logger
type catalogueService struct {
	db     *sqlx.DB
	logger log.Logger
}

// Product describes the fields of a product in the catalogue.
type Product struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	ImageURL    string   `json:"imageUrl" db:"image_url"`
	Price       float32  `json:"price" db:"price"`
	Stock       int      `json:"stock" db:"stock"`
	Type        []string `json:"type" db:"-"`
	TypeStr     string   `json:"-" db:"type_name"`
}

// Health describes the health of a service.
type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

var (
	ErrDBConnection = errors.New("Database connection failed")
	ErrNotFound     = errors.New("Not found")
	baseQuery       = "SELECT product.product_id AS id, product.name, product.description, product.price, product.stock, product.image_url, GROUP_CONCAT(type.name) AS type_name FROM product JOIN product_type ON product.product_id=product_type.product_id JOIN type ON product_type.type_id=type.type_id"
)

// NewCatalogueService returns an implementation of the Service interface,
// with connection to database and a logger
func NewCatalogueService(db *sqlx.DB, logger log.Logger) Service {
	return &catalogueService{
		db:     db,
		logger: logger,
	}
}

func (s *catalogueService) GetList(categories []string, orderBy string, currentPageNum, currentPageSize int) ([]Product, error) {
	var prods []Product
	query := baseQuery

	for i, c := range categories {
		if i == 0 {
			query += " WHERE type.name= " + c
		} else {
			query += " OR type.name= " + c
		}
	}
	query += " GROUP BY id"
	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}
	query += ";"

	fmt.Println("before select")
	err := s.db.Select(&prods, query)
	if err != nil {
		fmt.Println("err not nil", err)
		s.logger.Log("db error while selecting: ", err)
		return []Product{}, ErrDBConnection
	}
	fmt.Println("after select")

	if len(prods) == 0 {
		fmt.Println("prods = 0")
	} else {
		fmt.Println("prods != 0")
	}

	for i, p := range prods {
		prods[i].Type = strings.Split(p.TypeStr, ",")
	}
	fmt.Println(prods)
	prods = cut(prods, currentPageNum, currentPageSize)
	return prods, nil
}

func cut(prods []Product, pageNum, pageSize int) []Product { //TODO: REFACTOR
	if pageNum == 0 || pageSize == 0 {
		return []Product{} // pageNum is 1-indexed
	}
	start := (pageNum * pageSize) - pageSize
	if start > len(prods) {
		return []Product{}
	}
	end := (pageNum * pageSize)
	if end > len(prods) {
		end = len(prods)
	}
	return prods[start:end]
}

func (s *catalogueService) CountProducts(categories []string) (int, error) {
	query := "SELECT COUNT(DISTINCT product.product_id) FROM product JOIN product_type ON product.product_id=product_type.product_id JOIN type ON product_type.type_id=type.type_id"

	for i, c := range categories {
		if i == 0 {
			query += " WHERE type.name= " + c
		} else {
			query += " OR type.name= " + c
		}
	}
	query += ";"
	var count int
	err := s.db.Get(&count, query)
	if err != nil {
		fmt.Println("err not nil", err)
		s.logger.Log("db error while selecting: ", err)
		return 0, ErrDBConnection
	}
	return count, nil
}

func (s *catalogueService) GetProduct(id string) (Product, error) {
	query := baseQuery + " WHERE product.product_id = '" + id + "'"
	fmt.Println("id = ", id)

	query += " LIMIT 1;"
	var prod Product
	err := s.db.Get(&prod, query)
	if err != nil {
		fmt.Println("err not nil", err)
		s.logger.Log("product not found ", id, err)
		return Product{}, ErrNotFound
	}
	prod.Type = strings.Split(prod.TypeStr, ",")
	return prod, nil
}

func (s *catalogueService) GetCategories() ([]string, error) {
	query := "SELECT name FROM type;"
	var types []string
	err := s.db.Select(&types, query)
	if err != nil {
		s.logger.Log("database err while getting types: ", err)
		return []string{}, err
	}
	return types, nil
}
