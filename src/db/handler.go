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
	//go func(d *define.CardData) {
	err := globalDb.Create(data).Error
	if err != nil {
		sglog.Error("caarete card data error,", err)
	}
	//}(data)
	return err
}

func UpdateNoticeData(data *define.NoticeData) error {
	//go func(d *define.NoticeData) {
	err := globalDb.Where(define.NoticeData{Token: data.Token}).Assign(map[string]interface{}{
		"status":          data.Status,
		"name":            data.Name,
		"platform":        data.Platform,
		"title":           data.Title,
		"card_type":       data.CardType,
		"code":            data.Code,
		"phone":           data.Phone,
		"end_dt":          data.EndDt,
		"desc":            data.Desc,
		"renew_times":     data.RenewTimes,
		"notice_times":    data.NoticeTimes,
		"require_times":   data.RequireTimes,
		"final_login_dt":  data.FinalLoginDt,
		"create_dt":       data.CreateDt,
		"final_notice_dt": data.FinalNoticeDt,
		"share_times":     data.ShareTimes,
		"scene_id":        data.SceneId,
		"share_to_num":    data.ShareToNum,
		"shared_by":       data.SharedBy,
		"ad_complete_dt":  data.AdCompleteDt,
		"ad_times":        data.AdTimes}).FirstOrCreate(data).Error
	if err != nil {
		sglog.Error("update notice data error,token:", data.Token, err)
	}

	//}(data)
	return err
}

func UpdateStatisData(data *define.StatisticsData) error {
	err := globalDb.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateNoticeFinalData(data *define.NoticeFinalTime) error {
	err := globalDb.Where(define.NoticeFinalTime{Title: data.Title}).Assign(define.NoticeFinalTime{Time: data.Time}).FirstOrCreate(data).Error
	if err != nil {
		return err
	}
	return nil
}
