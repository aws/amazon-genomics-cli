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
	// and do our own prompting in processCommandIO.

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
	progressChan, wait := processOutputs(bufio.NewScanner(stdout), bufio.NewScanner(stderr), stdin, executionName)
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

func processOutputs(stdout *bufio.Scanner, stderr *bufio.Scanner, stdin io.WriteCloser, executionName string) (chan ProgressEvent, *sync.WaitGroup) {
	var wait sync.WaitGroup
	wait.Add(2)
	progressChan := make(chan ProgressEvent)
	currentEvent := &ProgressEvent{
		ExecutionName: executionName,
	}
	go func() {
		defer wait.Done()
		// We can't just scan through lines, because the MFA prompt is a
		// partial line and we need to see it. So scan runes instead and do our
		// own line buffering.
		stdout.Split(bufio.ScanRunes)
		line := ""
		for stdout.Scan() {
			line += stdout.Text()
			if strings.HasSuffix(line, "\n") {
				// We got a whole line at this character
				log.Debug().Msg(line[:len(line)-1])
				line = ""
			}
			if strings.HasSuffix(line, ": ") && strings.HasPrefix(line, "MFA token for") {
				// CDK may make MFA prompts here, so we need to forward them to the user.
				// We also need to make sure to drop down a couple lines
				// because if there's a progress spinner going it will just
				// immediately clobber our prompt.
				fmt.Printf("\n%s\n\n", line)
				line = ""

				oldStepDescription := currentEvent.StepDescription
				currentEvent.StepDescription = "Waiting for MFA..."
				progressChan <- *currentEvent

				// And we need to read and pass along a code.
				var reply string
				fmt.Scanln(&reply)

				currentEvent.StepDescription = oldStepDescription
				progressChan <- *currentEvent

				_, err := stdin.Write([]byte(reply + "\n"))
				if err != nil {
					log.Error().Msgf("error encountered while forwarding MFA code: %v", err)
				} else {
					log.Debug().Msg("Sent MFA code")
				}
				// We only need to send at most one MFA code, and if we don't
				// close its standard input we get stuck when the CDK is done
				// with its work.
				err = stdin.Close()
				if err != nil {
					log.Error().Msgf("error encountered while closing CDK input stream: %v", err)
				}
			}
		}
		if line != "" {
			// Handle any last unterminated line
			log.Debug().Msg(line)
		}
		err := stdout.Err()
		if err != nil {
			log.Error().Msgf("error encountered while scanning stdout: %v", err)
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
			log.Error().Msgf("error encountered while scanning stderr: %v", err)
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
