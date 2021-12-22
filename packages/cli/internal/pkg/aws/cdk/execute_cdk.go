package cdk

import (
	"bufio"
	"fmt"
    "io"
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

func executeCdkCommand(appDir string, commandArgs []string, executionName string) (ProgressStream, error) {
	return executeCdkCommandAndCleanupDirectory(appDir, commandArgs, "", executionName)
}

func executeCdkCommandAndCleanupDirectory(appDir string, commandArgs []string, tmpDir string, executionName string) (ProgressStream, error) {
	log.Debug().Msgf("executeCDKCommand(%s, %v)", appDir, commandArgs)
	cmdArgs := append([]string{"run", "cdk", "--"}, commandArgs...)
	cmd := execCommand("npm", cmdArgs...)
	cmd.Dir = appDir

    // Note that cmd won't have any access to stdin, stdout, or stderr that go
    // anywhere by default. It does not inherit our streams. This is a problem
    // because sometimes the CDK needs to dialog interactively with the user to
    // e.g. get MFA codes. So we make sure to forward along important output
    // and send input ourselves from our stdin in processCommandIO.

	progressChan, wait, err := processCommandIO(cmd, executionName)
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

func processCommandIO(cmd *exec.Cmd, executionName string) (chan ProgressEvent, *sync.WaitGroup, error) {
    stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
    stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
    if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("couldn't execute CDK deploy command: %w", err)
	}
	progressChan, wait := processOutputs(bufio.NewScanner(stdout), bufio.NewScanner(stderr), executionName)
    processInputs(stdin, executionName, wait)
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

func processInputs(stdin io.WriteCloser, executionName string, wait *sync.WaitGroup) {
	wait.Add(1)
	go func() {
        defer wait.Done()
        _, err := io.Copy(stdin, os.Stdin)
        if err != nil {
			log.Debug().Msgf("error encountered while copying stdin: %v", err)
		}
    }()
}

func processOutputs(stdout *bufio.Scanner, stderr *bufio.Scanner, executionName string) (chan ProgressEvent, *sync.WaitGroup) {
	var wait sync.WaitGroup
	wait.Add(2)
	progressChan := make(chan ProgressEvent)
	currentEvent := &ProgressEvent{
		ExecutionName: executionName,
	}
	go func() {
		defer wait.Done()
		for stdout.Scan() {
            line := stdout.Text()
			log.Debug().Msg(line)
            if strings.HasPrefix(line, "MFA token for") {
                // CDK may make MFA prompts here, so we need to forward them to the user
                fmt.Printf("%s\n", line)
            }
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
