// internal/services/database_service.go
package services

import "fmt"

type DatabaseService struct {
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func (s *DatabaseService) ReadDatabase(folderPath string) ([]DataItem, error) {
	return nil, fmt.Errorf("not implemented")
}

type DataItem struct {
	Key   string
	Value string
}
