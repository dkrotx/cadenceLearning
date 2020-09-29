package main

import (
	"errors"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

func registerImageOperationWorkflow(worker worker.Worker) {
	worker.RegisterWorkflowWithOptions(transformImageWorkflow, workflow.RegisterOptions{Name: "FlipImage"})
	worker.RegisterWorkflow(transformImageSubWorkflow)
	worker.RegisterActivity(downloadImage)
	worker.RegisterActivity(transformImage)
}

type FileOperationArgs struct {
	URL        string `json:"url"`
	OutputPath string `json:"output_path"`
}

func transformImageSubWorkflow(ctx workflow.Context, url, outputPath string) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Second * 5,
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	so := &workflow.SessionOptions{
		CreationTimeout:  time.Minute,
		ExecutionTimeout: time.Minute,
	}
	sessionCtx, err := workflow.CreateSession(ctx, so)
	if err != nil {
		return err
	}
	defer workflow.CompleteSession(sessionCtx)

	future := workflow.ExecuteActivity(sessionCtx, downloadImage, url, outputPath)
	if err := future.Get(sessionCtx, nil); err != nil {
		return err
	}

	future = workflow.ExecuteActivity(sessionCtx, transformImage, outputPath)
	if err = future.Get(sessionCtx, nil); err != nil {
		return err
	}

	return nil
}

func transformImageWorkflow(ctx workflow.Context, args FileOperationArgs) error {
	if args.URL == "" {
		return errors.New("URL is empty")
	}
	if args.OutputPath == "" {
		return errors.New("OutputPath is empty")
	}

	childWFID := "img-transform::" + workflow.GetInfo(ctx).WorkflowExecution.ID

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: childWFID,
		ExecutionStartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithChildOptions(ctx, cwo)

	workflow.GetLogger(ctx).Info("Starting child workflow", zap.String("workflow_id", childWFID))
	future := workflow.ExecuteChildWorkflow(ctx, transformImageSubWorkflow, args.URL, args.OutputPath)
	return future.Get(ctx, nil)
}

func downloadImage(url, filepath string) error {
	return downloadFile(url, filepath)
}

func transformImage(filepath string) error {
	return imageUpsideDown(filepath)
}
