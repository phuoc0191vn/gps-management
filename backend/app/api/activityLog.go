package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	utilities "ctigroupjsc.com/phuocnn/gps-management/uitilities"

	"gopkg.in/mgo.v2/bson"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"
	"github.com/julienschmidt/httprouter"
)

type ActivityLogHandler struct {
	ActivityLogRepository repository.ActivityLogRepository
	DeviceRepository      repository.DeviceRepository
	ReportRepository      repository.ReportRepository
	AccountRepository     repository.AccountRepository
}

func (h *ActivityLogHandler) GetInDay(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accountID := AccountID(r)
	if IsScopeAllowed(r) {
		accountID = p.ByName("accountID")
	}

	device, err := h.DeviceRepository.FindDeviceByStatus(accountID, model.StatusEnable)
	if err != nil || len(device) < 1 {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Code: http.StatusNotFound,
		})
		return
	}

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	activityLog, err := h.ActivityLogRepository.GetInDay(device[0].ID.Hex(), accountID, date)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: activityLog,
	})
}

func (h *ActivityLogHandler) CurrentLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accountID := AccountID(r)
	if IsScopeAllowed(r) {
		accountID = p.ByName("accountID")
	}

	device, err := h.DeviceRepository.FindDeviceByStatus(accountID, model.StatusEnable)
	if err != nil || len(device) < 1 {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Code: http.StatusNotFound,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code: http.StatusOK,
		Data: device[0],
	})
}

func (h *ActivityLogHandler) AllReport(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	accountIDs := make([]string, 0)
	condition := make(map[string]interface{})

	accountIDs = append(accountIDs, AccountID(r))
	if Scope(r) != model.ScopeUser {
		childAccounts, err := h.AccountRepository.GetChildAccounts(Email(r))
		if err != nil {
			childAccounts = make([]model.Account, 0)
		}

		for i := 0; i < len(childAccounts); i++ {
			accountIDs = append(accountIDs, childAccounts[i].ID.Hex())
		}
	}
	condition["accountID"] = bson.M{"$in": accountIDs}

	output, ok := GetQuery(r, DATATABLE_QUERY_OUTPUT)
	if ok && output == DATATABLE_QUERY_OUTPUT_DATATABLE {
		page := 1
		pageSize := 10
		var err error

		if r.URL.Query().Get("page") != "" {
			page, err = strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, ResponseBody{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				})
				return
			}
		}

		if r.URL.Query().Get("limit") != "" {
			pageSize, err = strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, ResponseBody{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
				})
				return
			}
		}

		total, data, err := h.ReportRepository.Pagination(page, pageSize, condition)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ResponseBody{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		WriteJSON(w, http.StatusOK, ResponseBody{
			Code: http.StatusOK,
			Data: struct {
				Total int
				Data  interface{}
			}{
				Total: total,
				Data:  data,
			},
		})
		return
	}

	// normal response
	data, err := h.ReportRepository.All()
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Data: data,
		Code: http.StatusOK,
	})
}

func (h *ActivityLogHandler) GenerateReport(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	start, end := GetRangeTimeFromQuery(r.URL.Query())
	start = start.UTC()
	end = end.UTC()

	id := p.ByName("id")
	accountID := AccountID(r)
	deviceID := id
	if Scope(r) != model.ScopeUser {
		device, err := h.DeviceRepository.FindDeviceByStatus(id, model.StatusEnable)
		if err != nil || len(device) < 1 {
			WriteJSON(w, http.StatusNotFound, ResponseBody{
				Code: http.StatusNotFound,
			})
			return
		}

		accountID = id
		deviceID = device[0].ID.Hex()
	}

	filter := model.Filter{
		AccountID: accountID,
		DeviceID:  deviceID,
		StartTime: start,
		EndTime:   end,
	}

	filename := fmt.Sprintf("Report_%d.csv", time.Now().Unix())
	account, err := h.AccountRepository.FindByID(accountID)
	if err == nil {
		filename = fmt.Sprintf("%s_Report_%d.csv", account.Email, time.Now().Unix())
	}

	report := model.Report{
		ID:          bson.NewObjectId(),
		AccountID:   accountID,
		Name:        filename,
		Filename:    filename,
		Status:      model.StatusUnprocessedReport,
		Type:        0,
		CreatedTime: utilities.TimeInUTC(time.Now()),
		Filter:      filter,
	}

	err = h.ReportRepository.Save(report)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: "unable to generate report",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code:    http.StatusOK,
		Message: "your report is generating",
		Data:    report.ID.Hex(),
	})
}

func (h *ActivityLogHandler) Download(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	report, err := h.ReportRepository.FindByID(p.ByName("id"))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Code: http.StatusNotFound,
		})
		return
	}

	bytesData, err := os.ReadFile(report.Filename)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Code:    http.StatusInternalServerError,
			Message: "unable to get report file",
		})
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", report.Name))
	w.Write(bytesData)
}

func (h *ActivityLogHandler) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	report, err := h.ReportRepository.FindByID(p.ByName("id"))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Code: http.StatusNotFound,
		})
		return
	}

	os.RemoveAll(report.Filename)
	err = h.ReportRepository.RemoveByID(report.ID.Hex())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Code:    http.StatusOK,
		Message: "delete report successfully",
	})
}

func GetRangeTimeFromQuery(query url.Values) (time.Time, time.Time) {
	now := time.Now().UTC()
	startStr := query.Get("startTime")
	startTime, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	endStr := query.Get("endTime")
	endTime, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		now = now.AddDate(0, 0, -7)
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	return startTime, endTime
}
