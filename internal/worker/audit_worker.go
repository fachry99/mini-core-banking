// pakai goroutine untuk audit logging agar proses transfer tidak ter-block. Lifecycle-nya saya kontrol via channel dan context,
// sehingga bisa graceful shutdown.

package worker

import (
	"context"
	"log"
)

type AuditEvent struct {
	RequestID string
	Message   string
}

type AuditWorker struct {
	ch chan AuditEvent
}

func NewAuditWorker(buffer int) *AuditWorker {
	return &AuditWorker{
		ch: make(chan AuditEvent, buffer),
	}
}

func (w *AuditWorker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("audit worker stopped")
				return
			case event := <-w.ch:
				log.Println("AUDIT:", event.RequestID, event.Message)
			}
		}
	}()
}

func (w *AuditWorker) Publish(event AuditEvent) {
	w.ch <- event
}
