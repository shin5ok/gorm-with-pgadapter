/*
This is just for local test with Spanner Emulator
Note: Before running this test, run spanner emulator and create an instance as "test-instance"
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	fakeServing Serving

	itemTestID = "d169f397-ba3f-413b-bc3c-a465576ef06e"
	userTestID string

/*
testDbName = genId()
*/
)

func init() {
	/*
		// TODO: setup some schemas just for test that will be destroyed at the end of the test
		schemaFiles, _ := filepath.Glob("schemas/*_ddl.sql")
		if err := testutil.InitData(ctx, fakeDbString, schemaFiles); err != nil {
			log.Fatal(err)
		}
	*/

	//prepareTestDB(spannerPgString, testDbName)

	//testSpannerPgString := rewriteDBString(spannerPgString, testDbName)
	testSpannerPgString := spannerPgString
	db, err := gorm.Open(postgres.Open(testSpannerPgString), &gorm.Config{
		DisableNestedTransaction: true,
	})

	if err != nil {
		log.Fatal(err)
	}
	fakeServing = Serving{
		Client: dbClient{sc: db},
	}
}

/*
	func prepareTestDB(originString, testDbName string) {
		prepareSpannerPgString := rewriteDBString(spannerPgString, "game")
		preDb, err := gorm.Open(postgres.Open(prepareSpannerPgString), &gorm.Config{
			DisableNestedTransaction: true,
		})
		if err != nil {
			panic(err)
		}
		defer func() {
			db, _ := preDb.DB()
			db.Close()
		}()
		createTestDb(preDb, testDbName)
	}

	func rewriteDBString(originSpannerPgString string, name string) string {
		r := regexp.MustCompile(`(.*)\s+database=.*`)
		results := r.FindAllStringSubmatch(originSpannerPgString, -1)
		newDbString := fmt.Sprintf("new DB string is '%s database=%s'", results[0][1], name)
		log.Println("Generating DB just for test", newDbString)
		return newDbString
	}

	func createTestDb(db *gorm.DB, name string) {
		sql := "create database %s " + name
		db.Exec(sql)
	}
*/

func Test_run(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeServing.pingPong)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected: %d. Got: %d, Message: %s", http.StatusOK, rr.Code, rr.Body)
	}

}

func Test_createUser(t *testing.T) {

	path := "test-user"
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("user_name", path)

	r := &http.Request{}
	req, err := http.NewRequestWithContext(r.Context(), "POST", "/api/user/"+path, nil)
	assert.Nil(t, err)
	newReq := req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeServing.createUser)
	handler.ServeHTTP(rr, newReq)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected: %d. Got: %d, Message: %s", http.StatusOK, rr.Code, rr.Body)
	}
	var u User
	json.Unmarshal(rr.Body.Bytes(), &u)
	userTestID = u.Id

}

// This test depends on Test_createUser
func Test_addItemUser(t *testing.T) {

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("user_id", userTestID)
	ctx.URLParams.Add("item_id", itemTestID)

	r := &http.Request{}
	uriPath := fmt.Sprintf("/api/user_id/%s/%s", userTestID, itemTestID)
	req, err := http.NewRequestWithContext(r.Context(), "PUT", uriPath, nil)
	assert.Nil(t, err)
	newReq := req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeServing.addItemToUser)
	handler.ServeHTTP(rr, newReq)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected: %d. Got: %d, Message: %s", http.StatusOK, rr.Code, rr.Body)
	}

}

func Test_getUserItems(t *testing.T) {

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("user_id", userTestID)

	r := &http.Request{}
	uriPath := fmt.Sprintf("/api/user_id/%s", userTestID)
	req, err := http.NewRequestWithContext(r.Context(), "GET", uriPath, nil)
	assert.Nil(t, err)
	newReq := req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeServing.addItemToUser)
	handler.ServeHTTP(rr, newReq)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected: %d. Got: %d, Message: %s", http.StatusOK, rr.Code, rr.Body)
	}
}

func Test_cleaning(t *testing.T) {
	t.Cleanup(func() {})
}
