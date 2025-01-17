// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package storage

import (
	"fmt"
	"time"

	"github.com/linuxboot/contest/pkg/event"
	"github.com/linuxboot/contest/pkg/event/frameworkevent"
	"github.com/linuxboot/contest/pkg/event/testevent"
	"github.com/linuxboot/contest/pkg/xcontext"
)

type EventStorage interface {
	// Test events storage interface
	StoreTestEvent(ctx xcontext.Context, event testevent.Event) error
	GetTestEvents(ctx xcontext.Context, eventQuery *testevent.Query) ([]testevent.Event, error)

	// Framework events storage interface
	StoreFrameworkEvent(ctx xcontext.Context, event frameworkevent.Event) error
	GetFrameworkEvent(ctx xcontext.Context, eventQuery *frameworkevent.Query) ([]frameworkevent.Event, error)
}

// TestEventEmitter implements Emitter interface from the testevent package
type TestEventEmitter struct {
	header testevent.Header
	// allowedEvents restricts the events this emitter will accept, if set
	allowedEvents *map[event.Name]bool
}

// TestEventFetcher implements the Fetcher interface from the testevent package
type TestEventFetcher struct {
}

// TestEventEmitterFetcher implements Emitter and Fetcher interface of the testevent package
type TestEventEmitterFetcher struct {
	TestEventEmitter
	TestEventFetcher
}

// Emit emits an event using the selected storage layer
func (e TestEventEmitter) Emit(ctx xcontext.Context, data testevent.Data) error {
	if e.allowedEvents != nil {
		if _, ok := (*e.allowedEvents)[data.EventName]; !ok {
			return fmt.Errorf("teststep %s is not allowed to emit unregistered event %s", e.header.TestName, data.EventName)
		}
	}
	event := testevent.Event{Header: &e.header, Data: &data, EmitTime: time.Now()}
	if err := storage.StoreTestEvent(ctx, event); err != nil {
		return fmt.Errorf("could not persist event data %v: %v", data, err)
	}
	return nil
}

// Fetch retrieves events based on QueryFields that are used to build a Query object for TestEvents
func (ev TestEventFetcher) Fetch(ctx xcontext.Context, queryFields ...testevent.QueryField) ([]testevent.Event, error) {
	eventQuery, err := testevent.QueryFields(queryFields).BuildQuery()
	if err != nil {
		return nil, fmt.Errorf("unable to build a query: %w", err)
	}

	if isStronglyConsistent(ctx) {
		return storage.GetTestEvents(ctx, eventQuery)
	}
	return storageAsync.GetTestEvents(ctx, eventQuery)
}

// NewTestEventEmitter creates a new Emitter object associated with a Header
func NewTestEventEmitter(header testevent.Header) testevent.Emitter {
	return TestEventEmitter{header: header}
}

// NewTestEventEmitterWithAllowedEvents creates a new Emitter object associated with a Header
func NewTestEventEmitterWithAllowedEvents(header testevent.Header, allowedEvents *map[event.Name]bool) testevent.Emitter {
	return TestEventEmitter{header: header, allowedEvents: allowedEvents}
}

// NewTestEventFetcher creates a new Fetcher object associated with a Header
func NewTestEventFetcher() testevent.Fetcher {
	return TestEventFetcher{}
}

// NewTestEventEmitterFetcher creates a new EmitterFetcher object associated with a Header
func NewTestEventEmitterFetcher(header testevent.Header) testevent.EmitterFetcher {
	return TestEventEmitterFetcher{
		TestEventEmitter{header: header},
		TestEventFetcher{},
	}
}

// NewTestEventEmitterFetcherWithAllowedEvents creates a new EmitterFetcher object associated with a Header
func NewTestEventEmitterFetcherWithAllowedEvents(header testevent.Header, allowedEvents *map[event.Name]bool) testevent.EmitterFetcher {
	return TestEventEmitterFetcher{
		TestEventEmitter{header: header},
		TestEventFetcher{},
	}
}

// FrameworkEventEmitter implements Emitter interface from the frameworkevent package
type FrameworkEventEmitter struct {
}

// FrameworkEventFetcher implements the Fetcher interface from the frameworkevent package
type FrameworkEventFetcher struct {
}

// FrameworkEventEmitterFetcher implements Emitter and Fetcher interface from the frameworkevent package
type FrameworkEventEmitterFetcher struct {
	FrameworkEventEmitter
	FrameworkEventFetcher
}

// Emit emits an event using the selected storage engine
func (ev FrameworkEventEmitter) Emit(ctx xcontext.Context, event frameworkevent.Event) error {
	if err := storage.StoreFrameworkEvent(ctx, event); err != nil {
		return fmt.Errorf("could not persist event %v: %v", event, err)
	}
	return nil
}

// Fetch retrieves events based on QueryFields that are used to build a Query object for FrameworkEvents
func (ev FrameworkEventFetcher) Fetch(ctx xcontext.Context, queryFields ...frameworkevent.QueryField) ([]frameworkevent.Event, error) {
	eventQuery, err := frameworkevent.QueryFields(queryFields).BuildQuery()
	if err != nil {
		return nil, fmt.Errorf("unable to build a query: %w", err)
	}

	if isStronglyConsistent(ctx) {
		return storage.GetFrameworkEvent(ctx, eventQuery)
	}
	return storageAsync.GetFrameworkEvent(ctx, eventQuery)
}

// NewFrameworkEventEmitter creates a new Emitter object for framework events
func NewFrameworkEventEmitter() FrameworkEventEmitter {
	return FrameworkEventEmitter{}
}

// NewFrameworkEventFetcher creates a new Fetcher object for framework events
func NewFrameworkEventFetcher() FrameworkEventFetcher {
	return FrameworkEventFetcher{}
}

// NewFrameworkEventEmitterFetcher creates a new EmitterFetcher object for framework events
func NewFrameworkEventEmitterFetcher() FrameworkEventEmitterFetcher {
	return FrameworkEventEmitterFetcher{
		FrameworkEventEmitter{},
		FrameworkEventFetcher{},
	}
}
