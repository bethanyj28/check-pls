package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/v43/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	status_queued      = "queued"
	status_in_progress = "in_progress"
	status_completed   = "completed"
	conclusion_success = "success"
	conclusion_failure = "failure"
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

	/*
		installationID := githubapp.GetInstallationIDFromEvent(&event)
		client, err := p.NewInstallationClient(installationID)
		if err != nil {
			return errors.Wrap(err, "failed to get installation client")
		}

		headSHA := event.GetHeadCommit().GetSHA()
		//headSHA := event.GetPullRequest().GetBase().GetSHA()
		var latestCommit *github.RepositoryCommit
		// replicating customer solution
			commits, _, err := client.PullRequests.ListCommits(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), 0, &github.ListOptions{PerPage: 100})
			if err != nil {
				return errors.Wrap(err, "failed to get commit list")
			}

			for _, commit := range commits {
				if commit.GetSHA() == headSHA {
					latestCommit = commit
					break
				}
			}

		if latestCommit == nil {
			latestCommit, _, err = client.Repositories.GetCommit(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), headSHA, &github.ListOptions{})
			if err != nil {
				return errors.Wrap(err, "failed to get commit")
			}
		}

		return doCheckRun(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), latestCommit.GetSHA(), 1, client, "push")
	*/
	return nil
}

type PullRequestHandler struct {
	githubapp.ClientCreator
}

func (pr *PullRequestHandler) Handles() []string {
	return []string{"pull_request"}
}

func (pr *PullRequestHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	zerolog.Ctx(ctx).Info().Msg("INCOMING PULL_REQUEST EVENT")

	var event github.PullRequestEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse pull_request event payload")
	}

	zerolog.Ctx(ctx).Info().Interface("event", event).Msg("PULL_REQUEST EVENT")

	if event.GetAction() != "synchronize" && event.GetAction() != "opened" && event.GetAction() != "reopened" {
		return nil
	}

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := pr.NewInstallationClient(installationID)
	if err != nil {
		return errors.Wrap(err, "failed to get installation client")
	}

	headSHA := event.GetPullRequest().GetHead().GetSHA()
	var latestCommit *github.RepositoryCommit
	// replicating customer solution
	commits, _, err := client.PullRequests.ListCommits(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), event.GetNumber(), &github.ListOptions{PerPage: 100})
	if err != nil {
		return errors.Wrap(err, "failed to get commit list")
	}

	for _, commit := range commits {
		if commit.GetSHA() == headSHA {
			latestCommit = commit
			break
		}
	}

	if latestCommit == nil {
		latestCommit, _, err = client.Repositories.GetCommit(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), headSHA, &github.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get commit")
		}
	}

	//cr, _, err := client.Checks.CreateCheckRun(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), github.CreateCheckRunOptions{Name: "test check run", HeadSHA: event.GetPullRequest().GetHead().GetSHA(), Status: &status_in_progress})
	return doCheckRun(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), latestCommit.GetSHA(), 1, client, "pr")
}

func doCheckRun(ctx context.Context, owner, repo, commitSHA string, numEventPRs int, client *github.Client, event string) error {
	cr, _, err := client.Checks.CreateCheckRun(ctx, owner, repo, github.CreateCheckRunOptions{Name: "test check run", HeadSHA: commitSHA, Status: &status_in_progress})
	if err != nil {
		return errors.Wrap(err, "failed to create check run")
	}

	zerolog.Ctx(ctx).Info().Interface("check run", cr).Msg("CHECK RUN CREATED")

	time.Sleep(10 * time.Second)

	text := fmt.Sprintf("The check suite has %d PRs", numEventPRs)
	title := "Hello"
	summary := fmt.Sprintf("I came from a %s event", event)
	cr, _, err = client.Checks.UpdateCheckRun(ctx, owner, repo, cr.GetID(), github.UpdateCheckRunOptions{Name: "test check run", Output: &github.CheckRunOutput{Title: &title, Summary: &summary, Text: &text}, Status: &status_completed, Conclusion: &conclusion_failure, CompletedAt: &github.Timestamp{Time: time.Now()}})
	if err != nil {
		return errors.Wrap(err, "failed to update check run")
	}

	zerolog.Ctx(ctx).Info().Interface("check run", cr).Msg("CHECK RUN UPDATED")

	return nil
}

type CheckSuiteHandler struct {
	githubapp.ClientCreator
}

func (cs *CheckSuiteHandler) Handles() []string {
	return []string{"check_suite"}
}

func (cs *CheckSuiteHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	zerolog.Ctx(ctx).Info().Msg("INCOMING CHECK_SUITE EVENT")

	var event github.CheckSuiteEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse check_suite event payload")
	}

	zerolog.Ctx(ctx).Info().Interface("event", event).Msg("CHECK_SUITE EVENT")
	if len(event.GetCheckSuite().PullRequests) == 0 {
		zerolog.Ctx(ctx).Warn().Msg("NO PULL REQUESTS")
	}

	/*
		installationID := githubapp.GetInstallationIDFromEvent(&event)
		client, err := cs.NewInstallationClient(installationID)
		if err != nil {
			return errors.Wrap(err, "failed to get installation client")
		}

		status := "completed"
		conclusion := "failure"

		cr, _, err := client.Checks.CreateCheckRun(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), github.CreateCheckRunOptions{Name: "test check run", HeadSHA: event.GetCheckSuite().GetHeadCommit().GetSHA(), Status: &status, Conclusion: &conclusion, CompletedAt: &github.Timestamp{Time: time.Now()}})
		if err != nil {
			return errors.Wrap(err, "failed to create check run")
		}

		zerolog.Ctx(ctx).Info().Interface("check run", cr).Msg("CHECK RUN CREATED")
	*/
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

	if event.GetAction() == "rerequested" {
		installationID := githubapp.GetInstallationIDFromEvent(&event)
		client, err := cr.NewInstallationClient(installationID)
		if err != nil {
			return errors.Wrap(err, "failed to get installation client")
		}
		if err := doCheckRun(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), event.GetCheckRun().GetHeadSHA(), len(event.GetCheckRun().GetCheckSuite().PullRequests), client, "check run rerun"); err != nil {
			return errors.Wrap(err, "failed to rerun check run")
		}
	}

	if len(event.GetCheckRun().GetCheckSuite().PullRequests) == 0 {
		zerolog.Ctx(ctx).Warn().Msg("NO PULL REQUESTS")
	}

	return nil
}
