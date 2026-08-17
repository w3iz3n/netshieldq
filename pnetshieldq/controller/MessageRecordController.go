package controller

import (
	"encoding/json"
	"net/http"
	"pnetshieldq/dao"
)

type MessageRecordController struct {
	dao *dao.MessageRecordDao
}

func NewMessageRecordController(dao *dao.MessageRecordDao) *MessageRecordController {
	return &MessageRecordController{dao: dao}
}

func (ctrl *MessageRecordController) FindAllMessages(w http.ResponseWriter, r *http.Request) {
	records, err := ctrl.dao.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将 records 序列化为 JSON 字符串
	jsonData, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 输出 JSON 字符串到控制台
	//fmt.Println(string(jsonData))

	// 设置响应头并返回 JSON 数据
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
