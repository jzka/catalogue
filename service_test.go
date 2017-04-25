package catalogue

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	p1 = Product{
		ID:          "1",
		Name:        "name1",
		Description: "description1",
		Price:       1.4, Stock: 2,
		ImageURL: "ImageUrl_1",
		Type:     []string{"type3"},
		TypeStr:  "type3",
	}

	p2 = Product{
		ID:          "2",
		Name:        "name2",
		Description: "description2",
		Price:       1.2, Stock: 2,
		ImageURL: "ImageUrl_2",
		Type:     []string{"type1"},
		TypeStr:  "type1",
	}
	p3 = Product{
		ID:          "3",
		Name:        "name3",
		Description: "description3",
		Price:       1.1,
		Stock:       3,
		ImageURL:    "ImageUrl_3",
		Type:        []string{"type1", "type3"},
		TypeStr:     "type1,type3",
	}

	p4 = Product{
		ID:          "4",
		Name:        "name4",
		Description: "description4",
		Price:       1.4,
		Stock:       4,
		ImageURL:    "ImageUrl_4",
		Type:        []string{"type3"},
		TypeStr:     "type3",
	}
	p5 = Product{
		ID:          "5",
		Name:        "name5",
		Description: "description5",
		Price:       1.5,
		Stock:       5,
		ImageURL:    "ImageUrl_5",
		Type:        []string{"type2", "type1"},
		TypeStr:     "type2,type1",
	}

	prods      = []Product{p1, p2, p3, p4, p5}
	categories = []string{"type1", "type2", "type3"}
)

var logger log.Logger

func TestCatalogueServiceGetList(t *testing.T) {
	logger = log.NewLogfmtLogger(os.Stderr)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub database connection error: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	cols := []string{"id", "name", "description", "price", "stock", "image_url", "type_name"}

	//Test Cases
	// pageSize 5, pageNumber 1, no order, no type
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(p1.ID, p1.Name, p1.Description, p1.Price, p1.Stock, p1.ImageURL, strings.Join(p1.Type, ",")).
		AddRow(p2.ID, p2.Name, p2.Description, p2.Price, p2.Stock, p2.ImageURL, strings.Join(p2.Type, ",")).
		AddRow(p3.ID, p3.Name, p3.Description, p3.Price, p3.Stock, p3.ImageURL, strings.Join(p3.Type, ",")).
		AddRow(p4.ID, p4.Name, p4.Description, p4.Price, p4.Stock, p4.ImageURL, strings.Join(p4.Type, ",")).
		AddRow(p5.ID, p5.Name, p5.Description, p5.Price, p5.Stock, p5.ImageURL, strings.Join(p5.Type, ",")))

	// pageSize 2, pageNumber 1, order Id, type type3
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(p1.ID, p1.Name, p1.Description, p1.Price, p1.Stock, p1.ImageURL, strings.Join(p1.Type, ",")).
		AddRow(p4.ID, p4.Name, p4.Description, p4.Price, p4.Stock, p4.ImageURL, strings.Join(p4.Type, ",")))

	// pageSize 3, pageNumber 1, order type, type no
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(p2.ID, p2.Name, p2.Description, p2.Price, p2.Stock, p2.ImageURL, strings.Join(p2.Type, ",")).
		AddRow(p1.ID, p1.Name, p1.Description, p1.Price, p1.Stock, p1.ImageURL, strings.Join(p1.Type, ",")).
		AddRow(p4.ID, p4.Name, p4.Description, p4.Price, p4.Stock, p4.ImageURL, strings.Join(p4.Type, ",")))

	s := NewCatalogueService(sqlxDB, logger)
	for _, testcase := range []struct {
		categories []string
		order      string
		pageNum    int
		pageSize   int
		want       []Product
	}{
		{
			categories: []string{},
			order:      "",
			pageNum:    1,
			pageSize:   5,
			want:       []Product{p1, p2, p3, p4, p5},
		},
		{
			categories: []string{"type3"},
			order:      "id",
			pageNum:    1,
			pageSize:   2,
			want:       []Product{p1, p4},
		},
		{
			categories: []string{},
			order:      "type",
			pageNum:    1,
			pageSize:   5,
			want:       []Product{p2, p1, p4},
		},
	} { //TODO: REFACTOR
		have, err := s.GetList(testcase.categories, testcase.order, testcase.pageNum, testcase.pageSize)
		if err != nil {
			t.Errorf(
				"List(%v, %s, %d, %d): returned error %s",
				testcase.categories, testcase.order, testcase.pageNum, testcase.pageSize,
				err.Error(),
			)
		}
		if want := testcase.want; !reflect.DeepEqual(want, have) {
			t.Errorf(
				"List(%v, %s, %d, %d): want %v, have %v",
				testcase.categories, testcase.order, testcase.pageNum, testcase.pageSize,
				want, have,
			)
		}
	}
}

func TestCatalogueServiceCount(t *testing.T) {
	logger = log.NewLogfmtLogger(os.Stderr)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub database connection error: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	cols := []string{"count"}

	//Test Cases
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).AddRow(3))
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).AddRow(1))
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).AddRow(1))

	s := NewCatalogueService(sqlxDB, logger)
	testcases := []struct { //annonymus struct
		types []string
		want  int
	}{
		{[]string{"type1"}, 3},
		{[]string{"type2"}, 1},
		{[]string{"type2", "type1"}, 1},
	}
	for _, testcase := range testcases {
		have, err := s.CountProducts(testcase.types)
		if err != nil {
			t.Errorf(
				"Count(%v): returned error %s",
				testcase.types, err.Error(),
				err.Error(),
			)
		}
		if want := testcase.want; want != have {
			t.Errorf("Count(%v): want %d, have %d", testcase.types, want, have)
		}
	}
}

func TestCatalogueServiceGetProduct(t *testing.T) {
	logger = log.NewLogfmtLogger(os.Stderr)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub database connection error: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	cols := []string{"id", "name", "description", "price", "stock", "image_url", "type_name"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(p1.ID, p1.Name, p1.Description, p1.Price, p1.Stock, p1.ImageURL, strings.Join(p1.Type, ",")).
		AddRow(p2.ID, p2.Name, p2.Description, p2.Price, p2.Stock, p2.ImageURL, strings.Join(p2.Type, ",")))

	s := NewCatalogueService(sqlxDB, logger)
	testcases := []struct { //annonymus struct
		id   string
		want interface{}
	}{
		{"1", "1"},
		{"0", ErrNotFound},
	}
	for i, testcase := range testcases {
		if i == 1 {
			_, have := s.GetProduct(testcase.id)
			if testcase.want != have {
				t.Errorf("GetProduct(%s):  want &s, have %s", testcase.id, testcase.want, have)
				continue
			}
		} else {
			have, err := s.GetProduct(testcase.id)
			if err != nil {
				t.Errorf("GetProduct(%s): %v", testcase.id, err)
				continue
			}
			if !reflect.DeepEqual(testcase.want, have.ID) {
				t.Errorf("GetProduct(%s): want %s, have %s", testcase.id, testcase.want, have.ID)
				continue
			}
		}
	}
}

func TestCatalogueServiceGetCategories(t *testing.T) {
	logger = log.NewLogfmtLogger(os.Stderr)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("stub database connection error: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	cols := []string{"name"}
	mock.ExpectQuery("SELECT name FROM type").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(categories[0]).
		AddRow(categories[1]).
		AddRow(categories[2]))
	s := NewCatalogueService(sqlxDB, logger)
	have, err := s.GetCategories()
	if err != nil {
		t.Errorf("Tags(): %v", err)
	}
	if !reflect.DeepEqual(categories, have) {
		t.Errorf("Tags(): want %v, have %v", categories, have)
	}
}
