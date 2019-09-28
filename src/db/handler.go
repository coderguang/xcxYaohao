package db

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func UpdateDownloadToDb(data *define.DownloadHistoryUrl) {
	//err := globalDb.Where(define.DownloadHistoryUrl{Title: data.Title, URL: data.URL}).Assign(define.DownloadHistoryUrl{Status: data.Status, Tips: data.Tips}).FirstOrCreate(data).Error
	err := globalDb.Where(define.DownloadHistoryUrl{Title: data.Title, URL: data.URL}).Assign(map[string]interface{}{
		"title":       data.Title,
		"url":         data.URL,
		"status":      data.Status,
		"download_dt": data.DownloadDt,
		"tips":        data.Tips}).FirstOrCreate(data).Error
	if err != nil {
		sglog.Error("update download url db error,", err)
	}
}

func UpdateCardData(data *define.CardData) error {
	err := globalDb.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateNoticeData(data *define.NoticeData) error {
	go func(d *define.NoticeData) {
		// err := globalDb.Where(define.NoticeData{Token: data.Token}).Assign(define.NoticeData{
		// 	Token:        data.Token,
		// 	Status:       data.Status,
		// 	Name:         data.Name,
		// 	Title:        data.Title,
		// 	CardType:     data.CardType,
		// 	Code:         data.Code,
		// 	Phone:        data.Phone,
		// 	EndDt:        data.EndDt,
		// 	Desc:         data.Desc,
		// 	RenewTimes:   data.RenewTimes,
		// 	NoticeTimes:  data.NoticeTimes,
		// 	RequireTimes: data.RequireTimes,
		// 	FinalLogin:   data.FinalLogin,
		// 	CreateDt:     data.CreateDt,
		// 	ShareTimes:   data.ShareTimes}).FirstOrCreate(data).Error

		err := globalDb.Where(define.NoticeData{Token: data.Token}).Assign(map[string]interface{}{
			"token":           data.Token,
			"status":          data.Status,
			"name":            data.Name,
			"title":           data.Title,
			"card_type":       data.CardType,
			"code":            data.Code,
			"phone":           data.Phone,
			"end_dt":          data.EndDt,
			"desc":            data.Desc,
			"renew_times":     data.RenewTimes,
			"notice_times":    data.NoticeTimes,
			"require_times":   data.RequireTimes,
			"final_notice_dt": data.FinalLogin,
			"create_dt":       data.CreateDt,
			"share_times":     data.ShareTimes}).FirstOrCreate(data).Error

		if err != nil {
			sglog.Error("update notice data error,token:", data.Token, err)
		}

	}(data)
	return nil
}
