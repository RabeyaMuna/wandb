package stream

import (
	"github.com/Khan/genqlient/graphql"
	"github.com/wandb/wandb/core/internal/featurechecker"
	"github.com/wandb/wandb/core/internal/observability"
	"github.com/wandb/wandb/core/internal/runupserter"
	"github.com/wandb/wandb/core/internal/runwork"
	"github.com/wandb/wandb/core/internal/settings"
	"github.com/wandb/wandb/core/internal/wboperation"
	spb "github.com/wandb/wandb/core/pkg/service_go_proto"
)

// RecordIngester turns Records into Work.
//
// Records coming from the client via interprocess communication, or those
// read from a transaction log, pass through here first.
type RecordIngester struct {
	ExtraWork          runwork.ExtraWork
	FeatureProvider    *featurechecker.ServerFeaturesCache
	GraphqlClientOrNil graphql.Client
	Logger             *observability.CoreLogger
	Operations         *wboperation.WandbOperations
	Run                *StreamRun

	Settings *settings.Settings
}

// Ingest pushes a record into the run work pipeline.
func (p *RecordIngester) Ingest(record *spb.Record) {
	var work runwork.Work

	if record.GetRun() != nil {
		work = &runupserter.RunUpdateWork{
			Record: record,

			StreamRunUpserter: p.Run,

			Settings:           p.Settings,
			BeforeRunEndCtx:    p.ExtraWork.BeforeEndCtx(),
			Operations:         p.Operations,
			FeatureProvider:    p.FeatureProvider,
			GraphqlClientOrNil: p.GraphqlClientOrNil,
			Logger:             p.Logger,
		}
	} else {
		// Legacy style for handling records where the code to process them
		// lives in handler.go and sender.go directly.
		work = runwork.WorkFromRecord(record)
	}

	p.ExtraWork.AddWork(work)
}
