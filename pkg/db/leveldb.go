package lvldb

import "testgui/internal/db"

type lvldbClient struct {
	basepath string
}

// Add implements db.DBClient.
func (l *lvldbClient) Add(key string, value string) bool {
	panic("unimplemented")
}

// Close implements db.DBClient.
func (l *lvldbClient) Close() bool {
	panic("unimplemented")
}

// Exists implements db.DBClient.
func (l *lvldbClient) Exists(key string) bool {
	panic("unimplemented")
}

// Get implements db.DBClient.
func (l *lvldbClient) Get(key string) string {
	panic("unimplemented")
}

// Open implements db.DBClient.
func (l *lvldbClient) Open() bool {
	panic("unimplemented")
}

func NewlvldbClient(address string) db.DBClient {
	return &lvldbClient{
		basepath: address,
	}
}
