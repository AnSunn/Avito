package storage

import (
	"database/sql"
	_ "github.com/lib/pq" //it is required for ini() function
	"log"
)

type Storage struct {
	//How to connect to DB
	config *Config
	//DB FileDescriptor
	db *sql.DB
	//Subfield for repo interface (model user)
	segmentRepository *SegmentRepository
	actionRepository  *ActionRepository
}

// Constructor for store
func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

// Open store method
func (storage *Storage) Open() error {
	//it is not a connection, it just validates arguments in Open function
	db, err := sql.Open("postgres", storage.config.DatabaseURL)
	if err != nil {
		return err
	}
	//it establishes a connection
	if err := db.Ping(); err != nil {
		return err
	}
	storage.db = db
	log.Println("DB connection is successfully created")
	return nil
}

// Close store method
func (s *Storage) Close() {
	s.db.Close()
}

// Public for SegmentRepo
func (s *Storage) Segment() *SegmentRepository {
	if s.segmentRepository != nil {
		return s.segmentRepository
	}
	s.segmentRepository = &SegmentRepository{
		store: s,
	}
	return s.segmentRepository
}

// Public for ActionRepo
func (s *Storage) Action() *ActionRepository {
	if s.actionRepository != nil {
		return s.actionRepository
	}
	s.actionRepository = &ActionRepository{
		store: s,
		s: &SegmentRepository{
			store: s,
		},
	}
	return s.actionRepository
}
