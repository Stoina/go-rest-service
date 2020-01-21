package restserver

import (
	"encoding/json"
	"net/url"

	db "github.com/Stoina/go-database"
	repo "github.com/Stoina/go-rest-server/repo"
)

// SQLRepository exported
// SQLRepository ...
type SQLRepository struct {
	name         	string
	url          	string
	dbConnection 	*db.Connection
	settings 		*SQLRepositorySettings
}

// NewSQLRepository exported
// NewSQLRepository ...
func NewSQLRepository(name string, url string, dbConn *db.Connection, settings *SQLRepositorySettings) *SQLRepository {
	return &SQLRepository{
		name:         name,
		url:          url,
		dbConnection: dbConn,
		settings: 	  settings}
}

// Post exported
// With a post request to the container resource you can create a new resource.
func (sqlRepo SQLRepository) Post(contentType string, content string) *repo.RepositoryResult {
	if contentType == "application/json" {
		return insertJSON(sqlRepo, content)
	}

	return nil
}

// Put exported
// With a put request to the container resource you can overwrite the resource with the representation in the request.
func (sqlRepo SQLRepository) Put(par string) *repo.RepositoryResult {
	return nil
}

// Patch exported
// With the http patch method, individual properties of a resource can be manipulated in a targeted manner.
func (sqlRepo SQLRepository) Patch(par string) *repo.RepositoryResult {
	return nil
}

// Delete exported
// With this method an existing resource can be deleted
func (sqlRepo SQLRepository) Delete(par string) *repo.RepositoryResult {
	return nil
}

// Get exported
// Get ...
func (sqlRepo SQLRepository) Get(calledURL *url.URL) *repo.RepositoryResult {
	return sqlRepo.Query("select * from \"" + sqlRepo.settings.TableName + "\"")
}

// Name exported
// Name ....
func (sqlRepo SQLRepository) Name() string {
	return sqlRepo.name
}

// URL exported
// URL ...
func (sqlRepo SQLRepository) URL() string {
	return sqlRepo.url
}

// Query exported
// Query ...
func (sqlRepo *SQLRepository) Query(query string) *repo.RepositoryResult {

	var resultError error
	responseMessage := ""
	responseData := ""

	queryResult, resultError := sqlRepo.dbConnection.Query(query)

	if resultError == nil {
		data, resultError := queryResult.ConvertToJSON()
		
		if resultError == nil {
			responseMessage = "Data loaded successfully"
			responseData = data
		}	
	}
	
	if resultError == nil {
		return repo.NewRepositoryResult(responseData, false, "", responseMessage, true)
	} 
	
	return repo.NewRepositoryResult(responseData, true, resultError.Error(), responseMessage, false)
}

func insertJSON(sqlRepository SQLRepository, jsonContent string) *repo.RepositoryResult {

	var jsonValues map[string]interface{}
	json.Unmarshal([]byte(jsonContent), &jsonValues)
	
	columndAndValueCount := len(jsonValues)

	columns := make([]string, columndAndValueCount) 
	values := make([]interface{}, columndAndValueCount)

	index := 0
	for key, value := range jsonValues {
		columns[index] = key
		values[index] = value

		index++
	}
	
	var resultError error
	responseMessage := ""
	responseData := ""

	insertStatement := db.NewInsertStatement(sqlRepository.settings.TableName, columns, values)
	dbResult, resultError := sqlRepository.dbConnection.Insert(insertStatement)

	if resultError == nil {
		responseData, resultError = dbResult.ConvertToJSON()
	}
	
	if resultError == nil {
		return repo.NewRepositoryResult(responseData, false, "", responseMessage, true)
	} 
	
	return repo.NewRepositoryResult(responseData, true, resultError.Error(), responseMessage, false)
}