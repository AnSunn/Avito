package api

//Full API initialization file
import (
	"encoding/json"
	"fmt"
	"github.com/AnSunn/ServerUserSegmentation/internal/app/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (api *API) PostSegment(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post Segment POST /segments")
	var segment models.Segment
	err := json.NewDecoder(req.Body).Decode(&segment)
	if err != nil {
		api.logger.Info("Invalid json received from the client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	a, err := api.storage.Segment().Create(&segment)
	if err != nil {
		api.logger.Info("Troubles while creating new segment:", err)
		msg := Message{
			StatusCode: 501,
			Message:    "Troubles while creating a new segment",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(a)

}

func (api *API) GetAllSegments(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	segments, err := api.storage.Segment().SelectAll()
	if err != nil {
		api.logger.Info(err)
		msg := Message{
			StatusCode: 501,
			Message:    "Problems extracting segments from the database.",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	api.logger.Info("Get All Segments GET /segments")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(segments)
}

func (api *API) GetSegmentByTitle(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Segment by name /api/v1/segments/{name}")
	title := mux.Vars(req)["title"]

	segment, ok, err := api.storage.Segment().FindSegmentByTitle(title)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (segments). err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "Troubles while accessing database table (segments)",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !ok {
		api.logger.Info("Can not find segment with this name in database")
		msg := Message{
			StatusCode: 404,
			Message:    "Segment with this name does not exists in database.",
			IsError:    true,
		}

		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(segment)

}

func (api *API) DeleteSegmentByTitle(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Delete Segment by Title DELETE /api/v1/segments/{title}")
	title := mux.Vars(req)["title"]

	_, ok, err := api.storage.Segment().FindSegmentByTitle(title)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (segments). err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "Troubles while accessing database table (segments)",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	if !ok {
		api.logger.Info("Can not find segment with this title in database")
		msg := Message{
			StatusCode: 404,
			Message:    "Segment with this title does not exist in database.",
			IsError:    true,
		}

		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	err = api.storage.Segment().DeleteByTitle(title)
	if err != nil {
		api.logger.Info("Troubles while deleting database element from table (segments). err:", err)
		msg := Message{
			StatusCode: 501,
			Message:    "Troubles while deleting database element from table (segments)",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(202)
	msg := Message{
		StatusCode: 202,
		Message:    fmt.Sprintf("Segment with title %s is successfully deleted (status field is changed to false).", title),
		IsError:    false,
	}
	json.NewEncoder(writer).Encode(msg)
}

//Actions

func (api *API) PostAction(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post Action POST /actions")
	var action models.Action
	err := json.NewDecoder(req.Body).Decode(&action)
	if err != nil {
		api.logger.Info("Invalid json received from the client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	err = api.storage.Action().Create(&action)
	if err != nil {
		api.logger.Info("Troubles while creating new action:", err)
		msg := Message{
			StatusCode: 501,
			Message:    "Troubles while creating new action",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	msg := Message{
		StatusCode: 201,
		Message:    "Actions were successfully created",
		IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
}

// Get Active Segments for User
func (api *API) GetActiveSegmentsForUser(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Segments by User_id GET /api/v1/segments/{user_id}")
	user_id, err := strconv.Atoi(mux.Vars(req)["user_id"])
	if err != nil {
		api.logger.Info("Troubles while parsing {user_id} param:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Unappropriated user_id value. Don't use user_id as uncasting to int value.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	segments, err := api.storage.Action().UserActiveSegmentsTitleByUserID(user_id)
	if err != nil {
		api.logger.Info(err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troubles extracting segments from actions database.",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	api.logger.Info("Get All Active User Segments GET /segments")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(segments)
}

func (api *API) GetActionDataForAPeriod(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Data for a period GET /api/v1/data/{year}{month}")
	year, err := strconv.Atoi(mux.Vars(req)["year"])
	if err != nil {
		api.logger.Info("Troubles while parsing {year} param:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Unappropriated year value. Don't use year as uncasting to int value.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	month, err := strconv.Atoi(mux.Vars(req)["month"])
	if err != nil {
		api.logger.Info("Troubles while parsing {month} param:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Unappropriated month value. Don't use month as uncasting to int value.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	data, err := api.storage.Action().ExtractDataForASpecialPeriod(year, month)
	if err != nil {
		api.logger.Info(err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troubles extracting data from actions database.",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	api.logger.Info("Get data from action table for a period GET /data")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(data)
}
