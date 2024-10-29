package vodmover

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/andreykaipov/goobs/api/events/subscriptions"
	"github.com/facebookincubator/go-belt/tool/logger"
)

type config = Config
type VODMover struct {
	serveCount atomic.Uint32
	config
}

func New(
	cfg Config,
) *VODMover {
	return &VODMover{
		config: cfg,
	}
}

func (m *VODMover) Serve(ctx context.Context) error {
	if m.serveCount.Add(1) > 1 {
		return fmt.Errorf("Serve could be used only once")
	}
	client, err := m.connectOBS()
	if err != nil {
		return fmt.Errorf("unable to connect to OBS: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
		case ev, ok := <-client.IncomingEvents:
			if !ok {
				_ = client.Disconnect()
				logger.Errorf(ctx, "the events channel was closed, reconnecting to OBS")
				client, err = m.reconnectOBS()
				if err != nil {
					return fmt.Errorf("unable to reconnect to OBS")
				}
				continue
			}

			err := m.processEvent(ctx, ev)
			if err != nil {
				logger.Errorf(ctx, "unable to process the event: %v", err)
			}
		}
	}
}

func (m *VODMover) processEvent(
	ctx context.Context,
	ev any,
) error {
	logger.Debugf(ctx, "processEventRecordStateChanged: %T", ev)
	switch ev := ev.(type) {
	case *events.RecordStateChanged:
		return m.processEventRecordStateChanged(ctx, ev)
	}
	return nil
}

func (m *VODMover) processEventRecordStateChanged(
	ctx context.Context,
	ev *events.RecordStateChanged,
) error {
	logger.Debugf(ctx, "processEventRecordStateChanged: %#+v", ev)
	defer logger.Debugf(ctx, "/processEventRecordStateChanged")

	switch ev.OutputState {
	case "OBS_WEBSOCKET_OUTPUT_STOPPED":
		return m.processRecordedFile(ctx, ev.OutputPath)
	}
	return nil
}

func (m *VODMover) processRecordedFile(
	ctx context.Context,
	filePath string,
) error {
	fileName := path.Base(filePath)

	logger.Debugf(ctx, "we have %d rules", len(m.config.MoveVODs))
	for _, rule := range m.config.MoveVODs {
		doesMatch := wildcard.Match(rule.PatternWildcard, fileName)
		if !doesMatch {
			logger.Debugf(ctx, "rule %#+v does NOT match (filename '%s')", fileName)
			continue
		}
		logger.Debugf(ctx, "rule %#+v DOES match (filename '%s')", fileName)
		err := m.moveFile(ctx, filePath, rule.Destination)
		if err != nil {
			return fmt.Errorf("unable to move file '%s' to destination '%s': %w", filePath, rule.Destination, err)
		}
		return nil
	}

	return fmt.Errorf("no rule defined for a file with name '%s' (%s)", fileName, filePath)
}

func (m *VODMover) moveFile(
	_ context.Context,
	filePath string,
	destination string,
) error {
	return os.Rename(filePath, destination)
}

func (m *VODMover) connectOBS() (*goobs.Client, error) {
	client, err := goobs.New(
		m.OBS.Address,
		goobs.WithPassword(m.OBS.Password),
		goobs.WithEventSubscriptions(subscriptions.Outputs), // look for RecordStateChanged in https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.json
	)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to OBS: %w", err)
	}
	return client, nil
}

func (m *VODMover) reconnectOBS() (*goobs.Client, error) {
	for {
		client, err := m.connectOBS()
		if err != nil {
			logger.Default().Errorf("unable to connect to OBS: %w", err)
			time.Sleep(time.Second)
			continue
		}
		return client, nil
	}
}
