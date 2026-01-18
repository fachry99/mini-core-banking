package audit

import (
	"time"

	"github.com/fachry/mini-core-banking/internal/logger"
)

func LogTransfer(
	requestID string,
	from string,
	to string,
	amount int64,
	status string,
) {
	logger.Log.Printf(
		"[AUDIT] request_id=%s from=%s to=%s amount=%d status=%s time=%s",
		requestID,
		from,
		to,
		amount,
		status,
		time.Now().Format(time.RFC3339),
	)
}
