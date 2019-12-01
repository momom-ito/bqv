// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/googleapis/google-cloud-go-testing/bigquery/bqiface"
	"github.com/pkg/errors"
	"github.com/rerost/bq-table-validator/domain/bqquery"
	"github.com/rerost/bq-table-validator/domain/tablemock"
	"github.com/rerost/bq-table-validator/domain/validator"
	"github.com/rerost/bqv/domain/annotateparser"
	"github.com/rerost/bqv/domain/viewmanager"
	"github.com/rerost/bqv/domain/viewservice"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

// Injectors from wire.go:

func InitializeCmd(ctx context.Context, cfg Config) (*cobra.Command, error) {
	viewService := viewservice.NewService()
	bqClient, err := NewBQClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	bqManager := viewmanager.NewBQManager(bqClient)
	fileManager := NewFileManager(cfg)
	extractor := annotateparser.NewExtractor()
	parser := annotateparser.NewParser()
	manifests := annotateparser.NewManifests(extractor, parser)
	client, err := NewRawBQClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	middleware := NewBQMiddleware(client)
	time := NewTime()
	tableMock := tablemock.NewTableMock(time)
	validatorValidator := validator.NewValidator(middleware, tableMock)
	command := NewCmdRoot(ctx, viewService, bqManager, fileManager, manifests, validatorValidator)
	return command, nil
}

// wire.go:

func NewBQClient(ctx context.Context, cfg Config) (viewmanager.BQClient, error) {
	c, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return viewmanager.BQClient(bqiface.AdaptClient(c)), nil
}

func NewFileManager(cfg Config) viewmanager.FileManager {
	return viewmanager.NewFileManager(cfg.Dir)
}

func NewRawBQClient(ctx context.Context, cfg Config) (bqiface.Client, error) {
	zap.L().Debug("Create BQ Client", zap.String("ProjectID", cfg.ProjectID))
	bqClient, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bqiface.AdaptClient(bqClient), nil
}

func NewBQMiddleware(bqClient bqiface.Client) validator.Middleware {
	return bqquery.NewBQQuery(bqClient)
}

func NewTime() time.Time {
	return time.Now()
}
