/*
* Editor:KaiWen
* PATH:models/upload_history.go
* Description:
 */

package models

import "time"

type UploadHistory struct {
	ID         int       `json:"id"`
	FileName   string    `json:"name"`
	UploadTime time.Time `json:"uploadTime"`
	Status     string    `json:"status"`
}

var UploadHistoryList []UploadHistory
