package history_dto

import (
	"my-project/helper"
	history_model "my-project/modul/history/model"
)

type Response struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	TableName string `json:"table_name"`
	RowID     int64  `json:"row_id"`
	Action    string `json:"action"`
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
	IP        string `json:"ip"`
	API       string `json:"api"`
	Method    string `json:"method"`
	CreatedAt string `json:"created_at"`
}

func ToResponse(history history_model.History) Response {
	return Response{
		ID:        history.ID,
		UserID:    history.UserID,
		TableName: history.TableName,
		RowID:     history.RowID,
		Action:    history.Action,
		OldValue:  history.OldValue,
		NewValue:  history.NewValue,
		IP:        history.IP,
		API:       history.API,
		Method:    history.Method,
		CreatedAt: helper.FormatDate(history.CreatedAt),
	}
}
