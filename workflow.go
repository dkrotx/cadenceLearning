package main

import (
	"go.uber.org/cadence"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

func defaultActivityOptions() workflow.ActivityOptions {
	retryPolicy := &cadence.RetryPolicy{
		InitialInterval:          time.Second,
		BackoffCoefficient:       2,
		MaximumInterval:          10 * time.Second,
		ExpirationInterval:       time.Minute * 5,
		NonRetriableErrorReasons: []string{"stop-retry"},
	}

	return workflow.ActivityOptions{
		TaskList:               TaskListName,
		ScheduleToCloseTimeout: time.Second * 5,
		ScheduleToStartTimeout: time.Second * 5,
		StartToCloseTimeout:    time.Second * 5,
		HeartbeatTimeout:       time.Second * 5,
		WaitForCancellation:    false,
		RetryPolicy:            retryPolicy,
	}
}

func confirmAge(ctx workflow.Context) bool {
	var signalVal string

	logger := workflow.GetLogger(ctx)

	logger.Info("Waiting for age confirmation")
	signalChan := workflow.GetSignalChannel(ctx, "age-confirmation")

	s := workflow.NewSelector(ctx)
	s.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Received age confirmation status", zap.String("value", signalVal))
	})
	s.Select(ctx)

	return signalVal == "confirmed"
}

// Asynchronously check for the stopping signal
func receivedStopSignal(ctx workflow.Context) bool {
	var signalVal string

	if workflow.GetSignalChannel(ctx, "stop").ReceiveAsync(&signalVal) {
		return true
	}

	return false
}

func SimpleWorkflowEx(ctx workflow.Context, value string) error {
	ctx = workflow.WithActivityOptions(ctx, defaultActivityOptions())

	future := workflow.ExecuteActivity(ctx, TransformNameActivity, value)
	var newName string
	if err := future.Get(ctx, &newName); err != nil {
		return err
	}

	future = workflow.ExecuteActivity(ctx, SayHelloActivity, newName)
	var result string
	if err := future.Get(ctx, &result); err != nil {
		return err
	}

	workflow.GetLogger(ctx).Info("Done", zap.String("result", result))
	return nil
}

func CheckDrivingLicenceWorkflow(ctx workflow.Context, driverName string) error {
	ctx = workflow.WithActivityOptions(ctx, defaultActivityOptions())
	logger := workflow.GetLogger(ctx)

	future := workflow.ExecuteActivity(ctx, FindUserInDatabase, driverName)
	var userInfo UserInformation
	if err := future.Get(ctx, &userInfo); err != nil {
		return err
	}

	//userInfo.Confirmed = confirmAge(ctx)
	future = workflow.ExecuteActivity(ctx, AllowedToDriveCar, userInfo)
	var result bool
	if err := future.Get(ctx, &result); err != nil {
		return err
	}

	logger.Info("Done", zap.Bool("result", result))

	if receivedStopSignal(ctx) {
		logger.Info("Stop signal is received, finishing workflow")
		return nil
	}

	logger.Info("Stop signal isn't received, running again")
	workflow.Sleep(ctx, time.Second) // prevent burning CPU
	return workflow.NewContinueAsNewError(ctx, CheckDrivingLicenceWorkflow, driverName)
}

func registerWorkflows(worker worker.Worker) {
	worker.RegisterWorkflowWithOptions(SimpleWorkflowEx, workflow.RegisterOptions{
		Name: "WelcomeUser",
	})

	worker.RegisterWorkflowWithOptions(CheckDrivingLicenceWorkflow, workflow.RegisterOptions{
		Name: "CheckDrivingLicence",
	})
}
