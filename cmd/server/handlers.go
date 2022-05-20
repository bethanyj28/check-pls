package main

import (
	"context"

	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
)

type PushHandler struct {
	githubapp.ClientCreator
}

func (p *PushHandler) Handles() []string {
	return []string{"push"}
}

func (p *PushHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	zerolog.Ctx(ctx).Info().Msg("PUSH IT")
	return nil
}

type CheckRunHandler struct {
	githubapp.ClientCreator
}

func (cr *CheckRunHandler) Handles() []string {
	return []string{"check_run"}
}

func (cr *CheckRunHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	zerolog.Ctx(ctx).Info().Msg("CHECK IT")
	return nil
}
