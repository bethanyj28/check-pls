package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/go-github/v43/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type PushHandler struct {
	githubapp.ClientCreator
}

func (p *PushHandler) Handles() []string {
	return []string{"push"}
}

func (p *PushHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	zerolog.Ctx(ctx).Info().Msg("INCOMING PUSH EVENT")

	var event github.PushEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse push event payload")
	}

	zerolog.Ctx(ctx).Info().Interface("event", event).Msg("PUSH EVENT")

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := p.NewInstallationClient(installationID)
	if err != nil {
		return errors.Wrap(err, "failed to get installation client")
	}

	/*
		cs, _, err := client.Checks.CreateCheckSuite(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), github.CreateCheckSuiteOptions{HeadSHA: event.GetHead()})
		if err != nil {
			return errors.Wrap(err, "failed to create check suite")
		}

		zerolog.Ctx(ctx).Info().Interface("check suite", cs).Msg("CHECK SUITE CREATED")
	*/

	status := "completed"
	conclusion := "failure"

	cr, _, err := client.Checks.CreateCheckRun(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), github.CreateCheckRunOptions{Name: "test check run", HeadSHA: event.GetHeadCommit().GetID(), Status: &status, Conclusion: &conclusion, CompletedAt: &github.Timestamp{Time: time.Now()}})
	if err != nil {
		return errors.Wrap(err, "failed to create check run")
	}

	zerolog.Ctx(ctx).Info().Interface("check run", cr).Msg("CHECK RUN CREATED")
	return nil
}

type CheckRunHandler struct {
	githubapp.ClientCreator
}

func (cr *CheckRunHandler) Handles() []string {
	return []string{"check_run"}
}

func (cr *CheckRunHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	zerolog.Ctx(ctx).Info().Msg("INCOMING CHECK_RUN EVENT")

	var event github.CheckRunEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse check run event payload")
	}

	zerolog.Ctx(ctx).Info().Interface("event", event).Msg("CHECK_RUN EVENT")

	return nil
}
