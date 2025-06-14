package runwork

import (
	"fmt"

	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

// SaveRecordWork saves a Record to the transaction log.
//
// It does nothing else.
type SaveRecordWork struct {
	Record *spb.Record
}

// Save implements Work.Save.
func (w *SaveRecordWork) Save(write func(*spb.Record)) {
	write(w.Record)
}

// DebugInfo implements Work.DebugInfo.
func (w *SaveRecordWork) DebugInfo() string {
	var recordType string
	switch x := w.Record.RecordType.(type) {
	case *spb.Record_Request:
		recordType = fmt.Sprintf("%T", x.Request.RequestType)
	default:
		recordType = fmt.Sprintf("%T", x)
	}

	return fmt.Sprintf(
		"SaveRecordWork(%s); Control(%v)",
		recordType, w.Record.GetControl())
}

// Accept implements Work.Accept.
func (w *SaveRecordWork) Accept(func(*spb.Record)) bool { return true }

// BypassOfflineMode implements Work.BypassOfflineMode.
func (w *SaveRecordWork) BypassOfflineMode() bool { return false }

// Process implements Work.Process.
func (w *SaveRecordWork) Process(func(*spb.Record), chan<- *spb.Result) {}

// Sentinel implements Work.Sentinel.
func (w *SaveRecordWork) Sentinel() any { return nil }
