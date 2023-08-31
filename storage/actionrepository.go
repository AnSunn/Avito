package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AnSunn/ServerUserSegmentation/internal/app/models"
	"github.com/AnSunn/ServerUserSegmentation/internal/app/models/responses"
	"log"
	"strings"
	"time"
)

type ActionRepository struct {
	store *Storage
	s     *SegmentRepository
}

var (
	tableAction string = "actions"
)

// For Post request
func (ac *ActionRepository) Create(a *models.Action) error {
	//Segments that have to be deleted. Start
	if a.Segment_title_to_del != nil {
		for _, seg := range a.Segment_title_to_del {
			segment, ok, err := ac.s.FindSegmentByTitle(seg)
			if err != nil {
				log.Println("There is an error in searching title process in segment table (delete action). The error: ", err)
			}
			if ok {
				end_date := time.Now().Format(time.RFC3339)
				query := fmt.Sprintf("UPDATE %s SET end_date = '%s' WHERE user_id = %d AND segment_id = %d", tableAction, end_date, a.User_id, segment.ID)
				_, err := ac.store.db.Exec(query)
				if err != nil {
					log.Println("There is an error in updating end_date for this user (delete action). The error: ", err)
				}
			} else {
				log.Println("There is no data with this title (delete action)")
			}
		}
	}
	//End

	//Segments that have to be added. Start
	if a.Segment_title_to_add != nil {
		act, err1 := ac.store.Action().UserActiveSegmentsTitleByUserID(a.User_id)
		for _, seg := range a.Segment_title_to_add {
			segment, ok, err := ac.s.FindSegmentByTitle(seg.Title)
			if err != nil {
				return err
			}
			if ok {
				IsActive := false
				if err1 != nil && !strings.Contains(err1.Error(), "no such active actions") {
					return err1
				}
				if err1 == nil {
					for _, val := range act {
						if seg.Title == val.Segment_title {
							IsActive = true
							break
						}
					}
				}
				start_date := time.Now().Format(time.RFC3339)
				if !IsActive {
					if seg.Days == 0 {
						query := fmt.Sprintf("INSERT INTO %s (user_id, segment_id, start_date) VALUES ($1, $2, $3) RETURNING id", tableAction)
						if err := ac.store.db.QueryRow(query, a.User_id, segment.ID, start_date).Scan(&a.ID); err != nil {
							return err
						}
					} else {
						end_date := time.Now().AddDate(0, 0, seg.Days).Format(time.RFC3339)
						query := fmt.Sprintf("INSERT INTO %s (user_id, segment_id, start_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id", tableAction)
						if err := ac.store.db.QueryRow(query, a.User_id, segment.ID, start_date, end_date).Scan(&a.ID); err != nil {
							return err
						}
					}
				}
			} else {
				return errors.New("There is a problem adding data to action table")
			}
		}
	}
	//End
	return nil
}

func (a *ActionRepository) FindActionBySegmentId(id int) (*models.Action, error) {
	//!!!!
	action := &models.Action{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE segment_id =$1", tableAction)
	row := a.store.db.QueryRow(query, id)
	var start_date_raw, end_date_raw sql.NullTime
	if err := row.Scan(&action.ID, &action.User_id, &action.Segment_id, &start_date_raw, &end_date_raw); err != nil {
		if err == sql.ErrNoRows {
			return action, fmt.Errorf("Action with id = %d: no such action", id)
		}
		return action, fmt.Errorf("Action with id = %d: %v", id, err)
	}
	action.Start_date = start_date_raw.Time.Format(time.RFC3339)
	action.End_date = end_date_raw.Time.Format(time.RFC3339)
	return action, nil
}

func (a *ActionRepository) UserActiveActionsByUserID(id int) ([]*models.Action, error) {

	today := time.Now().Format(time.RFC3339)
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND '%s' between start_date and "+
		"COALESCE (end_date, $2)", tableAction, today)
	rows, err := a.store.db.Query(query, id, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := make([]*models.Action, 0)
	var start_date_raw, end_date_raw sql.NullTime
	for rows.Next() {
		action := models.Action{}
		if err = rows.Scan(&action.ID, &action.User_id, &action.Segment_id, &start_date_raw, &end_date_raw); err != nil {
			if err == sql.ErrNoRows {
				return data, fmt.Errorf("Active actions for user_id = %d: no such active actions", id)
			}
			return data, fmt.Errorf("Active actions for user_id = %d: %v", id, err)
		}
		action.Start_date = start_date_raw.Time.Format(time.RFC3339)
		action.End_date = end_date_raw.Time.Format(time.RFC3339)
		data = append(data, &action)
	}
	return data, nil
}

func (a *ActionRepository) UserActiveSegmentsTitleByUserID(id int) ([]*responses.ActiveUserSegments, error) {
	data, err := a.UserActiveActionsByUserID(id)
	if err != nil {
		return nil, err
	}
	user_segments := make([]*responses.ActiveUserSegments, 0)
	for _, val := range data {
		user_segment := responses.ActiveUserSegments{}
		segments, _ := a.store.Segment().FindSegmentById(val.Segment_id)
		user_segment.User_id = val.User_id
		user_segment.Segment_title = segments.Title
		user_segments = append(user_segments, &user_segment)
	}
	return user_segments, nil
}

// For delete segment function
func (a *ActionRepository) UserActiveActions() ([]*models.Action, error) {
	today := time.Now().Format(time.RFC3339)
	query := fmt.Sprintf("SELECT * FROM %s WHERE '%s' between start_date and "+
		"COALESCE (end_date, $1)", tableAction, today)
	rows, err := a.store.db.Query(query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := make([]*models.Action, 0)
	var start_date_raw, end_date_raw sql.NullTime
	for rows.Next() {
		action := models.Action{}
		if err = rows.Scan(&action.ID, &action.User_id, &action.Segment_id, &start_date_raw, &end_date_raw); err != nil {
			if err == sql.ErrNoRows {
				return data, fmt.Errorf("Active actions for all users = %d: nobody has active actions")
			}
			return data, fmt.Errorf("Active actions for all users = %d: %v", err)
		}
		action.Start_date = start_date_raw.Time.Format(time.RFC3339)
		action.End_date = end_date_raw.Time.Format(time.RFC3339)
		data = append(data, &action)
	}
	return data, nil
}

// For additional task #1
func (a *ActionRepository) ExtractDataForASpecialPeriod(year, month int) ([]*responses.DataByPeriod, error) {
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	firstMonthString := firstOfMonth.Format(time.RFC3339)
	LastMonthString := lastOfMonth.Format(time.RFC3339)

	query := fmt.Sprintf("SELECT * FROM %s WHERE  to_date(start_date::TEXT, 'YYYY MM DD') between '%s' and '%s'", tableAction, firstMonthString, LastMonthString)
	rows, err := a.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := make([]*responses.DataByPeriod, 0)
	for rows.Next() {
		action := models.Action{}
		var startDateRaw, endDateRaw sql.NullTime
		err = rows.Scan(&action.ID, &action.User_id, &action.Segment_id, &startDateRaw, &endDateRaw)
		if err != nil {
			return nil, err
		}
		action.Start_date = startDateRaw.Time.Format(time.RFC3339)
		segmentTitle, _ := a.store.Segment().FindSegmentById(action.Segment_id)
		selectedData := responses.DataByPeriod{
			User_id:       action.ID,
			Segment_title: segmentTitle.Title,
			Operation:     responses.Add,
			Date:          action.Start_date,
		}
		data = append(data, &selectedData)
	}
	query = fmt.Sprintf("SELECT * FROM %s WHERE to_date(end_date::TEXT, 'YYYY MM DD') between '%s' and '%s'", tableAction, firstMonthString, LastMonthString)
	rows, err = a.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		action := models.Action{}
		var startDateRaw, endDateRaw sql.NullTime
		err = rows.Scan(&action.ID, &action.User_id, &action.Segment_id, &startDateRaw, &endDateRaw)
		if err != nil {
			return nil, err
		}
		action.End_date = endDateRaw.Time.Format(time.RFC3339)
		segmentTitle, _ := a.store.Segment().FindSegmentById(action.Segment_id)
		selectedData := responses.DataByPeriod{
			User_id:       action.ID,
			Segment_title: segmentTitle.Title,
			Operation:     responses.Delete,
			Date:          action.End_date,
		}
		data = append(data, &selectedData)
	}
	return data, err
}
