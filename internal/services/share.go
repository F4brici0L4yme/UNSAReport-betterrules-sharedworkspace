package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/UNSAReport/UNSAReport/internal/ports"
)

// CommandRunner executes an external command and returns its combined output.
type CommandRunner func(ctx context.Context, name string, args ...string) (string, error)

// LookPathFunc resolves an executable name to a path, mirroring exec.LookPath.
type LookPathFunc func(name string) (string, error)

// ShareService publishes the current project directory as a new GitHub repository via the gh CLI.
type ShareService struct {
	FS       ports.FileSystem
	Stdout   io.Writer
	Stderr   io.Writer
	Runner   CommandRunner
	LookPath LookPathFunc
}

// ShareOption configures a ShareService via functional options.
type ShareOption func(*ShareService)

// WithShareFS sets the filesystem used to inspect the working directory.
func WithShareFS(fs ports.FileSystem) ShareOption {
	return func(s *ShareService) { s.FS = fs }
}

// WithShareStdout sets the writer for standard output messages.
func WithShareStdout(w io.Writer) ShareOption {
	return func(s *ShareService) { s.Stdout = w }
}

// WithShareStderr sets the writer for standard error messages.
func WithShareStderr(w io.Writer) ShareOption {
	return func(s *ShareService) { s.Stderr = w }
}

// WithShareRunner sets the command runner used to invoke gh and git.
func WithShareRunner(r CommandRunner) ShareOption {
	return func(s *ShareService) { s.Runner = r }
}

// WithShareLookPath sets the executable resolver used to check for gh.
func WithShareLookPath(f LookPathFunc) ShareOption {
	return func(s *ShareService) { s.LookPath = f }
}

// NewShareService creates a ShareService with the given functional options applied.
func NewShareService(opts ...ShareOption) *ShareService {
	s := &ShareService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ShareOptions holds the parameters for a single share execution.
type ShareOptions struct {
	Name        string
	Visibility  string // "public" or "private"
	Description string
}

// Execute ensures the directory is a git repository, then creates and pushes a new GitHub repository.
func (s *ShareService) Execute(ctx context.Context, opt ShareOptions) error {
	visibility := opt.Visibility
	if visibility == "" {
		visibility = "private"
	}
	if visibility != "public" && visibility != "private" {
		return fmt.Errorf("invalid visibility %q (must be \"public\" or \"private\")", opt.Visibility)
	}

	runner := s.Runner
	if runner == nil {
		runner = defaultCommandRunner
	}
	lookPath := s.LookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}

	if _, err := lookPath("gh"); err != nil {
		return fmt.Errorf("the gh CLI is required to share a project (https://cli.github.com): %w", err)
	}

	cwd, err := s.FS.Getwd()
	if err != nil {
		return fmt.Errorf("get cwd: %w", err)
	}

	name := strings.TrimSpace(opt.Name)
	if name == "" {
		name = filepath.Base(cwd)
	}
	if name == "" || name == "." || name == string(filepath.Separator) {
		return fmt.Errorf("could not derive a repository name from the current directory; pass one explicitly")
	}

	if err := s.ensureGitRepo(ctx, runner); err != nil {
		return err
	}

	args := []string{"repo", "create", name, "--source=.", "--push"}
	if visibility == "public" {
		args = append(args, "--public")
	} else {
		args = append(args, "--private")
	}
	if strings.TrimSpace(opt.Description) != "" {
		args = append(args, "--description", opt.Description)
	}

	if _, err := fmt.Fprintf(s.Stdout, "Creating GitHub repository %q...\n", name); err != nil {
		return fmt.Errorf("write message: %w", err)
	}
	if _, err := runner(ctx, "gh", args...); err != nil {
		return fmt.Errorf("gh repo create: %w", err)
	}

	url, err := runner(ctx, "gh", "repo", "view", name, "--json", "url", "--jq", ".url")
	if err == nil {
		url = strings.TrimSpace(url)
	}
	if url == "" {
		url = fmt.Sprintf("https://github.com/%s", name)
	}

	if _, err := fmt.Fprintf(s.Stdout, "Repository ready: %s\n", url); err != nil {
		return fmt.Errorf("write url: %w", err)
	}
	return nil
}

// ensureGitRepo initializes a git repository if none exists and creates an initial commit when needed.
func (s *ShareService) ensureGitRepo(ctx context.Context, runner CommandRunner) error {
	if !s.FS.FileExists(".git") {
		if _, err := runner(ctx, "git", "init"); err != nil {
			return fmt.Errorf("git init: %w", err)
		}
	}

	if _, err := runner(ctx, "git", "rev-parse", "HEAD"); err == nil {
		return nil
	}

	if _, err := runner(ctx, "git", "add", "-A"); err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	if _, err := runner(ctx, "git", "commit", "-m", "Initial report"); err != nil {
		return fmt.Errorf("git commit: %w (set git config user.name and user.email first)", err)
	}
	return nil
}

func defaultCommandRunner(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return buf.String(), err
	}
	return buf.String(), nil
}
