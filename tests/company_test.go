package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mogw/micro-company/internal/auth"
	"github.com/mogw/micro-company/internal/company"
	"github.com/mogw/micro-company/internal/config"
	"github.com/mogw/micro-company/internal/db"
	"github.com/mogw/micro-company/internal/kafka"
	"github.com/stretchr/testify/assert"
)

var (
	router *gin.Engine
)

func setupRouter() {
	cfg := config.LoadConfig()

	mongoClient, err := db.ConnectMongo(cfg.MongoURI)
	if err != nil {
		panic("Failed to connect to MongoDB")
	}

	kafkaProducer := kafka.NewProducer(cfg.KafkaBroker)

	companyRepo := company.NewRepository(mongoClient, "companydb", "companies")
	companyService := company.NewService(companyRepo, kafkaProducer)
	companyHandler := company.NewHandler(companyService)

	router = gin.Default()
	router.Use(auth.JWTMiddleware(cfg.JWTSecret))

	companyHandler.RegisterRoutes(router)
}

func TestMain(m *testing.M) {
	setupRouter()
	m.Run()
}

func getToken() string {
	// Implement JWT token generation or use a mock token for testing purposes
	return "your_test_jwt_token"
}

func TestCreateCompany(t *testing.T) {
	token := getToken()

	comp := company.Company{
		Name:              "Test Company",
		Description:       "This is a test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              "Corporations",
	}

	jsonValue, _ := json.Marshal(comp)
	req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdCompany company.Company
	json.Unmarshal(w.Body.Bytes(), &createdCompany)
	assert.Equal(t, comp.Name, createdCompany.Name)
	assert.Equal(t, comp.Description, createdCompany.Description)
	assert.Equal(t, comp.AmountOfEmployees, createdCompany.AmountOfEmployees)
	assert.Equal(t, comp.Registered, createdCompany.Registered)
	assert.Equal(t, comp.Type, createdCompany.Type)
}

func TestGetCompany(t *testing.T) {
	token := getToken()

	comp := company.Company{
		Name:              "Test Company Get",
		Description:       "This is a test company for get",
		AmountOfEmployees: 20,
		Registered:        true,
		Type:              "NonProfit",
	}

	jsonValue, _ := json.Marshal(comp)
	req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdCompany company.Company
	json.Unmarshal(w.Body.Bytes(), &createdCompany)

	req, _ = http.NewRequest("GET", "/companies/"+createdCompany.UUID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedCompany company.Company
	json.Unmarshal(w.Body.Bytes(), &fetchedCompany)
	assert.Equal(t, createdCompany.UUID, fetchedCompany.UUID)
	assert.Equal(t, createdCompany.Name, fetchedCompany.Name)
	assert.Equal(t, createdCompany.Description, fetchedCompany.Description)
	assert.Equal(t, createdCompany.AmountOfEmployees, fetchedCompany.AmountOfEmployees)
	assert.Equal(t, createdCompany.Registered, fetchedCompany.Registered)
	assert.Equal(t, createdCompany.Type, fetchedCompany.Type)
}

func TestUpdateCompany(t *testing.T) {
	token := getToken()

	comp := company.Company{
		Name:              "Test Company Update",
		Description:       "This is a test company for update",
		AmountOfEmployees: 30,
		Registered:        true,
		Type:              "Cooperative",
	}

	jsonValue, _ := json.Marshal(comp)
	req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdCompany company.Company
	json.Unmarshal(w.Body.Bytes(), &createdCompany)

	update := map[string]interface{}{
		"name":                "Updated Company Name",
		"amount_of_employees": 50,
	}
	jsonValue, _ = json.Marshal(update)
	req, _ = http.NewRequest("PATCH", "/companies/"+createdCompany.UUID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/companies/"+createdCompany.UUID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedCompany company.Company
	json.Unmarshal(w.Body.Bytes(), &updatedCompany)
	assert.Equal(t, "Updated Company Name", updatedCompany.Name)
	assert.Equal(t, 50, updatedCompany.AmountOfEmployees)
}

func TestDeleteCompany(t *testing.T) {
	token := getToken()

	comp := company.Company{
		Name:              "Test Company Delete",
		Description:       "This is a test company for delete",
		AmountOfEmployees: 40,
		Registered:        true,
		Type:              "Sole Proprietorship",
	}

	jsonValue, _ := json.Marshal(comp)
	req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdCompany company.Company
	json.Unmarshal(w.Body.Bytes(), &createdCompany)

	req, _ = http.NewRequest("DELETE", "/companies/"+createdCompany.UUID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	req, _ = http.NewRequest("GET", "/companies/"+createdCompany.UUID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}