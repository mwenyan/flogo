package funit

import (
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/trigger"
)

func init() {
}

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

// Factory for trigger
type Factory struct {
}

// New implements trigger.Factory.New
func (t *Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	if err := metadata.MapToStruct(config.Settings, s, true); err != nil {
		return nil, err
	}
	trig := Trigger{}

	return &trig, nil
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Trigger is a stub for the Trigger implementation
type Trigger struct {
}

// Initialize implements trigger.Init.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	return nil
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() error {
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *Trigger) Stop() error {
	// stop the trigger
	return nil
}
