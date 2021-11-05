package cdk

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/rs/zerolog/log"
)

var ExecuteCdkCommand = executeCdkCommand
var execCommand = exec.Command
var progressRegex = regexp.MustCompile(`^.*\|\s*([0-9]+/[0-9]+)\s*\|(.*)`)

var osRemoveAll = os.RemoveAll

func executeCdkCommand(appDir string, commandArgs []string, uniqueKey string) (ProgressStream, error) {
	return executeCdkCommandAndCleanupDirectory(appDir, commandArgs, "", uniqueKey)
}

func executeCdkCommandAndCleanupDirectory(appDir string, commandArgs []string, tmpDir string, uniqueKey string) (ProgressStream, error) {
	log.Debug().Msgf("executeCDKCommand(%s, %v)", appDir, commandArgs)
	cmdArgs := append([]string{"run", "cdk", "--"}, commandArgs...)
	cmd := execCommand("npm", cmdArgs...)
	cmd.Dir = appDir

	progressChan, wait, err := processCommandOutputs(cmd, uniqueKey)
	if err != nil {
		deleteCDKOutputDir(tmpDir)
		return nil, err
	}

	go func() {
		defer close(progressChan)
		defer deleteCDKOutputDir(tmpDir)
		wait.Wait()
		err = cmd.Wait()
		if err != nil {
			progressChan <- ProgressEvent{Err: err}
		}
	}()
	return progressChan, nil
}

func processCommandOutputs(cmd *exec.Cmd, uniqueKey string) (chan ProgressEvent, *sync.WaitGroup, error) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("couldn't execute CDK deploy command: %w", err)
	}
	progressChan, wait := processOutputs(bufio.NewScanner(stdout), bufio.NewScanner(stderr), uniqueKey)
	return progressChan, wait, nil
}

func deleteCDKOutputDir(cdkOutputDir string) {
	if cdkOutputDir == "" {
		return
	}
	if err := osRemoveAll(cdkOutputDir); err != nil {
		log.Error().Err(err).Msgf("tried to delete output from cdk from location '%s' but failed", cdkOutputDir)
	}
}

func processOutputs(stdout *bufio.Scanner, stderr *bufio.Scanner, uniqueKey string) (chan ProgressEvent, *sync.WaitGroup) {
	var wait sync.WaitGroup
	wait.Add(2)
	progressChan := make(chan ProgressEvent)
	currentEvent := &ProgressEvent{
		UniqueKey: uniqueKey,
	}
	go func() {
		defer wait.Done()
		for stdout.Scan() {
			log.Debug().Msg(stdout.Text())
		}
		err := stdout.Err()
		if err != nil {
			log.Debug().Msgf("error encountered while scanning stdout: %v", err)
		}
	}()
	go func() {
		defer wait.Done()
		for stderr.Scan() {
			line := stderr.Text()
			progressChan <- updateEvent(currentEvent, line)
		}
		err := stderr.Err()
		if err != nil {
			log.Debug().Msgf("error encountered while scanning stderr: %v", err)
		}
	}()
	return progressChan, &wait
}

func updateEvent(event *ProgressEvent, line string) ProgressEvent {
	event.Outputs = append(event.Outputs, line)
	event.LastOutput = line
	match := progressRegex.FindStringSubmatch(line)
	if len(match) == 3 {
		stepParts := strings.Split(match[1], "/")
		event.StepDescription = match[2]
		if currentStep, err := strconv.Atoi(stepParts[0]); err != nil {
			log.Debug().Msgf("Unable to convert current step '%s' to int: %v", stepParts[0], err)
		} else {
			event.CurrentStep = currentStep
		}
		if totalSteps, err := strconv.Atoi(stepParts[1]); err != nil {
			log.Debug().Msgf("Unable to convert total steps '%s' to int: %v", stepParts[1], err)
		} else {
			event.TotalSteps = totalSteps
		}
		event.StepDescription = match[2]
	}
	return *event
}
