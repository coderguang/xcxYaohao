package db

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func UpdateDownloadToDb(data *define.DownloadHistoryUrl) {
	//err := globalDb.Where(define.DownloadHistoryUrl{Title: data.Title, URL: data.URL}).Assign(define.DownloadHistoryUrl{Status: data.Status, Tips: data.Tips}).FirstOrCreate(data).Error
	err := globalDb.Where(define.DownloadHistoryUrl{Title: data.Title, URL: data.URL}).Assign(define.DownloadHistoryUrl{Status: data.Status, Tips: data.Tips}).FirstOrCreate(data).Error
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
		err := globalDb.Where(define.NoticeData{Token: data.Token}).Assign(define.NoticeData{
			Token:         data.Token,
			Status:        data.Status,
			Name:          data.Name,
			Title:         data.Title,
			CardType:      data.CardType,
			Code:          data.Code,
			Phone:         data.Phone,
			EndDt:         data.EndDt,
			Desc:          data.Desc,
			RenewTimes:    data.RenewTimes,
			NoticeTimes:   data.NoticeTimes,
			RequireTimes:  data.RequireTimes,
			FinalLoginDt:  data.FinalLoginDt,
			CreateDt:      data.CreateDt,
			FinalNoticeDt: data.FinalNoticeDt,
			ShareTimes:    data.ShareTimes}).FirstOrCreate(data).Error
		if err != nil {
			sglog.Error("update notice data error,token:", data.Token, err)
		}

	}(data)
	return nil
}
