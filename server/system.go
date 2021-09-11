package server

import (
	"adp_backend/config"
	"adp_backend/model"

	utime "adp_backend/infra/time"
)

func AddAuditLog(e *config.Env, userName, sourceIP, event, eventArgs, eventResult string) error {
	var al model.AuditLog
	al.LoginUser = userName
	al.SourceIp = sourceIP
	al.EventArgs = eventArgs
	al.Event = event
	al.EventResult = eventResult
	al.CreatedAt = utime.CurTime()

	return nil
}
