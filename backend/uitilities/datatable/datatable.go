package datatable

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DATATABLE_QUERY_COLUMN_DATA         = "columns[%d][data]"
	DATATABLE_QUERY_COLUMN_NAME         = "columns[%d][name]"
	DATATABLE_QUERY_COLUMN_SEARCHABLE   = "columns[%d][searchable]"
	DATATABLE_QUERY_COLUMN_ORDERABLE    = "columns[%d][orderable]"
	DATATABLE_QUERY_COLUMN_SEARCH_VALUE = "columns[%d][search][value]"
	DATATABLE_QUERY_COLUMN_SEARCH_REGEX = "columns[%d][search][regex]"
	DATATABLE_QUERY_ORDER_COLUMN        = "order[0][column]"
	DATATABLE_QUERY_ORDER_DIR           = "order[0][dir]"
	DATATABLE_QUERY_START               = "start"
	DATATABLE_QUERY_LENGTH              = "length"
	DATATABLE_QUERY_SEARCH_VALUE        = "search[value]"
	DATATABLE_QUERY_SEARCH_REGEX        = "search[regex]"
	QUERY_ORDER_ASC                     = "asc"
	QUERY_ORDER_DESC                    = "desc"
	DATATABLE_DEFAULT_START             = 0
	DATATABLE_DEFAULT_LENGTH            = 10
	MESSAGE_VALIDATION_ERROR            = "Invalid Data"
)

type DataTableColumnQuery struct {
	Data        string
	Name        string
	SearchValue string
	SearchRegex string
	Searchable  bool
	Orderable   bool
}

type DataTableQuery struct {
	Columns     []DataTableColumnQuery
	OrderColumn int
	OrderDir    string
	Start       int
	Length      int
	SearchValue string
	SearchRegex bool
}

type Frequency struct {
	Label string `bson:"_id" json:"label"`
	Data  int    `bson:"count" json:"data"`
}

type Search struct {
	Value string
	Regex bool
}

type DataTableResult struct {
	Data            interface{} `json:"data"`
	RecordsTotal    int         `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
}

type ResponseBody struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func MongoDataTable(queryParams url.Values, condition map[string]interface{}, selection map[string]interface{}, col *mgo.Collection) (*DataTableResult, error) {
	if condition == nil {
		condition = make(map[string]interface{})
	}

	if selection == nil {
		selection = make(map[string]interface{})
	}

	dataTableQuery, ok := ParseDataTableQuery(queryParams)
	if !ok {
		return nil, fmt.Errorf(MESSAGE_VALIDATION_ERROR)
	}

	if !ValidData(dataTableQuery) {
		return nil, fmt.Errorf(MESSAGE_VALIDATION_ERROR)
	}

	total := 0
	data := make([]interface{}, 0)

	if dataTableQuery.SearchValue != "" {
		searchCondition := make([]interface{}, 0)
		for _, cols := range dataTableQuery.Columns {
			if cols.Data == "" {
				continue
			}
			searchCondition = append(searchCondition, bson.M{cols.Data: bson.M{"$regex": dataTableQuery.SearchValue, "$options": "i"}})
		}

		if len(searchCondition) > 0 {
			tmpOr, ok := condition["$or"].([]interface{})
			if !ok {
				tmpOr = make([]interface{}, 0)
			}

			tmpOr = append(tmpOr, searchCondition...)
			condition["$or"] = tmpOr
		}
	}

	query := col.Find(condition).Select(selection)
	total, e := query.Count()
	if e != nil {
		return nil, e
	}

	if total > 0 {
		e = query.Skip(dataTableQuery.Start).Limit(dataTableQuery.Length).Sort(fmt.Sprintf("%s%s", dataTableQuery.OrderDir, FindSortField(dataTableQuery))).All(&data)

		if e != nil {
			return nil, e
		}
	}

	result := new(DataTableResult)
	result.Data = data
	result.RecordsFiltered = total
	result.RecordsTotal = total
	return result, nil
}

func FindSortField(data DataTableQuery) string {
	index := data.OrderColumn
	val := data.Columns[index].Data
	if val == "" {
		val = "_id"
	}

	return val
}

func ValidData(data DataTableQuery) bool {
	if data.OrderColumn >= len(data.Columns) {
		return false
	}

	return true
}

func ParseDataTableQuery(queryParams url.Values) (data DataTableQuery, existed bool) {
	columnNumber := 0
	for {
		columnData, ok := queryParams[fmt.Sprintf(DATATABLE_QUERY_COLUMN_DATA, columnNumber)]
		if !ok || len(columnData) < 1 {
			break
		}

		column := DataTableColumnQuery{}
		column.Data = columnData[0]
		column.Name = queryParams.Get(fmt.Sprintf(DATATABLE_QUERY_COLUMN_NAME, columnNumber))
		column.SearchValue = queryParams.Get(fmt.Sprintf(DATATABLE_QUERY_COLUMN_SEARCH_VALUE, columnNumber))
		column.SearchRegex = queryParams.Get(fmt.Sprintf(DATATABLE_QUERY_COLUMN_SEARCH_REGEX, columnNumber))
		column.Searchable = ParseBoolQuery(queryParams, fmt.Sprintf(DATATABLE_QUERY_COLUMN_SEARCHABLE, columnNumber))
		column.Orderable = ParseBoolQuery(queryParams, fmt.Sprintf(DATATABLE_QUERY_COLUMN_ORDERABLE, columnNumber))
		columnNumber++

		data.Columns = append(data.Columns, column)
	}

	if columnNumber == 0 {
		return
	}

	existed = true
	data.OrderColumn = ParseIntQueryWithDefault(queryParams, DATATABLE_QUERY_ORDER_COLUMN, 0)
	data.OrderDir = "+"
	if queryParams.Get(DATATABLE_QUERY_ORDER_DIR) == QUERY_ORDER_DESC {
		data.OrderDir = "-"
	}
	data.Start = ParseIntQueryWithDefault(queryParams, DATATABLE_QUERY_START, DATATABLE_DEFAULT_START)
	data.Length = ParseIntQueryWithDefault(queryParams, DATATABLE_QUERY_LENGTH, DATATABLE_DEFAULT_LENGTH)
	data.SearchValue = queryParams.Get(DATATABLE_QUERY_SEARCH_VALUE)
	data.SearchRegex = ParseBoolQuery(queryParams, DATATABLE_QUERY_SEARCH_REGEX)
	return
}

func MakeDateCondition(fromValue, toValue string) map[string]interface{} {
	dateCondition := make(map[string]interface{})

	if fromValue != "" {
		fromDate, e := ParseISODate(fromValue)
		if e == nil && !fromDate.IsZero() {
			dateCondition["$gte"] = fromDate
		}
	}

	if toValue != "" {
		toDate, e := ParseISODate(toValue)
		if e == nil && !toDate.IsZero() {
			dateCondition["$lte"] = toDate
		}
	}

	return dateCondition
}

func ParseDateCondition(from, to, format string) (dateCondition map[string]interface{}, mappingTime map[string]Frequency, e error) {
	dateCondition = make(map[string]interface{})
	mappingTime = make(map[string]Frequency)
	fromDate := time.Time{}
	toDate := time.Time{}

	if from != "" {
		fromDate, e = ParseISODate(from)

		if e != nil {
			return
		}

		from = fromDate.Format(format)
		dateCondition["$gte"] = from
	}

	if to != "" {
		toDate, e = ParseISODate(to)

		if e != nil {
			return
		}

		to = toDate.Format(format)
		dateCondition["$lte"] = to
	}

	if from != "" && to != "" {
		tmp := fromDate
		for from <= to {
			from = tmp.Format(format)
			mappingTime[from] = Frequency{
				Label: from,
			}
			tmp = tmp.AddDate(0, 0, 1)
		}
	}

	return
}

func ParseBoolQuery(queryParams url.Values, key string) bool {
	trueStr := "true"
	val := strings.TrimSpace(strings.ToLower(queryParams.Get(key)))
	if val == trueStr {
		return true
	}

	return false
}

func ParseIntQueryWithDefault(queryParams url.Values, key string, defaultVal int) int {
	val, e := strconv.ParseInt(strings.TrimSpace(strings.ToLower(queryParams.Get(key))), 10, 64)

	if e != nil {
		return defaultVal
	}

	return int(val)
}

func ParseISODate(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}
