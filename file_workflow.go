package main

import (
	"errors"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"time"
)

func registerImageOperationWorkflow(worker worker.Worker) {
	worker.RegisterWorkflowWithOptions(transformImageWorkflow, workflow.RegisterOptions{Name: "FlipImage"})
	worker.RegisterActivity(downloadImage)
	worker.RegisterActivity(transformImage)
}

type FileOperationArgs struct {
	URL string	`json:"url"`
	OutputPath string `json:"output_path"`
}

func transformImageWorkflow(ctx workflow.Context, args FileOperationArgs) error {
	if args.URL == "" {
		return errors.New("URL is empty")
	}
	if args.OutputPath == "" {
		return errors.New("OutputPath is empty")
	}

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

	future := workflow.ExecuteActivity(sessionCtx, downloadImage, args.URL, args.OutputPath)
	if err := future.Get(sessionCtx, nil); err != nil {
		return err
	}

	future = workflow.ExecuteActivity(sessionCtx, transformImage, args.OutputPath)
	if err = future.Get(sessionCtx, nil); err != nil {
		return err
	}

	return nil
}

func downloadImage(url, filepath string) error {
	return downloadFile(url, filepath)
}

func transformImage(filepath string) error {
	return imageUpsideDown(filepath)
}
