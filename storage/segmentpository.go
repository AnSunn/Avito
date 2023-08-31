package storage

import (
	"database/sql"
	"fmt"
	"github.com/AnSunn/ServerUserSegmentation/internal/app/models"
	"log"
	"time"
)

type SegmentRepository struct {
	store *Storage
}

var (
	tableSegment string = "segments"
)

// For Post request
func (s *SegmentRepository) Create(a *models.Segment) (*models.Segment, error) {
	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id, status", tableSegment)
	if err := s.store.db.QueryRow(query, a.Title).Scan(&a.ID, &a.Status); err != nil {
		return nil, err
	}
	return a, nil
}

// Get all request and helper for FindByID
func (s *SegmentRepository) SelectAll() ([]*models.Segment, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableSegment)
	rows, err := s.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	segments := make([]*models.Segment, 0)
	for rows.Next() {
		a := models.Segment{}
		err := rows.Scan(&a.ID, &a.Title, &a.Status)
		if err != nil {
			log.Println(err)
			continue
		}
		segments = append(segments, &a)
	}
	return segments, nil
}

// Helper for Delete by id and GET by id request
func (s *SegmentRepository) FindSegmentByTitle(title string) (*models.Segment, bool, error) {
	seg := &models.Segment{}
	founded := false
	query := fmt.Sprintf("SELECT * FROM %s WHERE title='%s'", tableSegment, title)
	row := s.store.db.QueryRow(query)
	if err := row.Scan(&seg.ID, &seg.Title, &seg.Status); err != nil {
		if err == sql.ErrNoRows {
			return seg, founded, fmt.Errorf("Segment by title %s: no such segment", title)
		}
		return seg, founded, fmt.Errorf("Segment by title %s: %v", title, err)
	}
	founded = true
	return seg, founded, nil
}

// Helper for ExtractDataForASpecialPeriod and GetActiveSegmentsByUserId
func (s *SegmentRepository) FindSegmentById(id int) (*models.Segment, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", tableSegment)
	row := s.store.db.QueryRow(query, id)
	seg := &models.Segment{}
	if err := row.Scan(&seg.ID, &seg.Title, &seg.Status); err != nil {
		if err == sql.ErrNoRows {
			return seg, fmt.Errorf("Segment by id %d: no such segment", id)
		}
		return seg, fmt.Errorf("Segment by id %d: %v", id, err)
	}
	return seg, nil
}

// For DELETE request
func (s *SegmentRepository) DeleteByTitle(title string) error {
	segment, ok, err := s.FindSegmentByTitle(title)
	if err != nil {
		return err
	}
	if ok {
		query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2", tableSegment)
		_, err = s.store.db.Exec(query, false, segment.ID)
		if err != nil {
			return err
		}
		actions, err1 := s.store.Action().UserActiveActions()
		if err1 != nil {
			return err1
		}
		for _, val := range actions {
			if val.Segment_id == segment.ID {
				end_date := time.Now().Format(time.RFC3339)
				query = fmt.Sprintf("UPDATE %s SET end_date = $1 WHERE segment_id = $2 and id =$3", tableAction)
				_, err = s.store.db.Exec(query, end_date, segment.ID, val.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
