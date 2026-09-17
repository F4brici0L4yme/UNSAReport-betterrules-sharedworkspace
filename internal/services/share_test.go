package services

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/UNSAReport/UNSAReport/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShareService_Execute_MissingGh(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	svc := NewShareService(
		WithShareFS(mocks.NewFileSystem(t)),
		WithShareStdout(&stdout),
		WithShareLookPath(func(string) (string, error) { return "", errors.New("not found") }),
	)

	err := svc.Execute(context.Background(), ShareOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gh CLI is required")
}

func TestShareService_Execute_DerivesNameAndCreatesPrivateRepo(t *testing.T) {
	t.Parallel()

	fs := mocks.NewFileSystem(t)
	fs.On("Getwd").Return("/tmp/mi-informe", nil)
	fs.On("FileExists", ".git").Return(true)

	var runnerCalls []string
	runner := func(_ context.Context, name string, args ...string) (string, error) {
		runnerCalls = append(runnerCalls, name+" "+strings.Join(args, " "))
		if name == "gh" && len(args) > 0 && args[0] == "repo" && args[1] == "view" {
			return "https://github.com/mi-informe", nil
		}
		return "", nil
	}

	var stdout bytes.Buffer
	svc := NewShareService(
		WithShareFS(fs),
		WithShareStdout(&stdout),
		WithShareRunner(runner),
		WithShareLookPath(func(string) (string, error) { return "/usr/bin/gh", nil }),
	)

	err := svc.Execute(context.Background(), ShareOptions{})
	require.NoError(t, err)

	assert.Contains(t, runnerCalls, "gh repo create mi-informe --source=. --push --private")
	assert.Contains(t, stdout.String(), "https://github.com/mi-informe")
}

func TestShareService_Execute_InitAndCommitWhenNotARepo(t *testing.T) {
	t.Parallel()

	fs := mocks.NewFileSystem(t)
	fs.On("Getwd").Return("/tmp/nuevo-proyecto", nil)
	fs.On("FileExists", ".git").Return(false)

	var calls []string
	runner := func(_ context.Context, name string, args ...string) (string, error) {
		calls = append(calls, name+" "+strings.Join(args, " "))
		if name == "git" && args[0] == "rev-parse" {
			return "", errors.New("no commit")
		}
		if name == "gh" && len(args) > 0 && args[0] == "repo" && args[1] == "view" {
			return "https://github.com/nuevo-proyecto", nil
		}
		return "", nil
	}

	var stdout bytes.Buffer
	svc := NewShareService(
		WithShareFS(fs),
		WithShareStdout(&stdout),
		WithShareRunner(runner),
		WithShareLookPath(func(string) (string, error) { return "/usr/bin/gh", nil }),
	)

	err := svc.Execute(context.Background(), ShareOptions{Name: "nuevo-proyecto"})
	require.NoError(t, err)

	assert.Contains(t, calls, "git init")
	assert.Contains(t, calls, "git add -A")
	assert.Contains(t, calls, "git commit -m Initial report")
	assert.Contains(t, calls, "gh repo create nuevo-proyecto --source=. --push --private")
}

func TestShareService_Execute_InvalidVisibility(t *testing.T) {
	t.Parallel()

	svc := NewShareService(
		WithShareFS(mocks.NewFileSystem(t)),
		WithShareStdout(&bytes.Buffer{}),
		WithShareRunner(func(_ context.Context, _ string, _ ...string) (string, error) { return "", nil }),
		WithShareLookPath(func(string) (string, error) { return "/usr/bin/gh", nil }),
	)

	err := svc.Execute(context.Background(), ShareOptions{Visibility: "internal"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid visibility")
}
