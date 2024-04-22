package engine

import (
	"context"

	"github.com/rss3-network/protocol-go/schema"
	activityx "github.com/rss3-network/protocol-go/schema/activity"
	"github.com/rss3-network/protocol-go/schema/tag"
)

type Worker interface {
	// Name is the name of the worker.
	Name() string
	// Platform returns the display name of the worker as the `platform` in the final Activity response.
	Platform() string
	// Tags the possible `tag` of the worker, displayed in the final Activity response.
	Tags() []tag.Tag
	// Types the possible `type` of the worker, displayed in the final Activity response.
	Types() []schema.Type
	// Filter the SourceFilter of the worker(network, state, start logics, etc.).
	Filter() SourceFilter
	// Match the matching logic of worker and returns if the task is matched
	Match(ctx context.Context, task Task) (bool, error)
	// Transform the core logic of the worker and returns the Activity.
	Transform(ctx context.Context, task Task) (*activityx.Activity, error)
}
