package admin

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	mu        sync.Mutex
	processes map[string]*Process
}

type Process struct {
	Name      string
	Cmd       *exec.Cmd
	Running   bool
	Logs      []string
	LogLimit  int
	StartArgs []string
	Dir       string
}

func NewService() *Service {
	return &Service{
		processes: make(map[string]*Process),
	}
}

func (s *Service) RegisterProcess(name string, dir string, command string, args ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.processes[name] = &Process{
		Name:      name,
		Running:   false,
		Logs:      make([]string, 0),
		LogLimit:  1000,
		StartArgs: append([]string{command}, args...),
		Dir:       dir,
	}
}

func (s *Service) Start(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	proc, ok := s.processes[name]
	if !ok {
		return fmt.Errorf("process %s not found", name)
	}

	if proc.Running {
		return fmt.Errorf("process %s is already running", name)
	}

	cmd := exec.Command(proc.StartArgs[0], proc.StartArgs[1:]...)
	cmd.Dir = proc.Dir
	proc.Cmd = cmd

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	proc.Running = true
	s.log(name, fmt.Sprintf("Process started with PID %d", cmd.Process.Pid))

	// Stream logs
	go s.streamLogs(name, stdout)
	go s.streamLogs(name, stderr)

	// Monitor exit
	go func() {
		err := cmd.Wait()
		s.mu.Lock()
		defer s.mu.Unlock()
		proc.Running = false
		if err != nil {
			s.log(name, fmt.Sprintf("Process exited with error: %v", err))
		} else {
			s.log(name, "Process exited successfully")
		}
	}()

	return nil
}

func (s *Service) Stop(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	proc, ok := s.processes[name]
	if !ok {
		return fmt.Errorf("process %s not found", name)
	}

	if !proc.Running || proc.Cmd == nil || proc.Cmd.Process == nil {
		return fmt.Errorf("process %s is not running", name)
	}

	// On Windows, Process.Kill() only kills the parent process, not child
	// processes (e.g., uvicorn workers). Use taskkill /T /F to kill the
	// entire process tree so the port is actually freed.
	var killErr error
	if runtime.GOOS == "windows" {
		pid := strconv.Itoa(proc.Cmd.Process.Pid)
		kill := exec.Command("taskkill", "/T", "/F", "/PID", pid)
		// Capture output instead of piping to os.Stdout to avoid
		// interleaving with logrus output.
		out, err := kill.CombinedOutput()
		if len(out) > 0 {
			s.log(name, strings.TrimSpace(string(out)))
		}
		killErr = err
	} else {
		// On Unix, send SIGKILL to the process group
		killErr = proc.Cmd.Process.Kill()
	}

	if killErr != nil {
		return killErr
	}

	proc.Running = false
	s.log(name, "Process stopped by user")
	return nil
}

func (s *Service) Restart(name string) error {
	// Best effort stop
	_ = s.Stop(name)
	// Wait a bit
	time.Sleep(1 * time.Second)
	return s.Start(name)
}

func (s *Service) GetLogs(name string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	proc, ok := s.processes[name]
	if !ok {
		return nil, fmt.Errorf("process %s not found", name)
	}

	// Return copy
	logs := make([]string, len(proc.Logs))
	copy(logs, proc.Logs)
	return logs, nil
}

func (s *Service) streamLogs(name string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		text := scanner.Text()
		s.mu.Lock()
		s.log(name, text)
		s.mu.Unlock()
	}
}

func (s *Service) log(name string, msg string) {
	proc := s.processes[name]
	timestamp := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[%s] %s", timestamp, msg)

	proc.Logs = append(proc.Logs, logLine)
	if len(proc.Logs) > proc.LogLimit {
		proc.Logs = proc.Logs[len(proc.Logs)-proc.LogLimit:]
	}
	logrus.Infof("[%s] %s", name, msg)
}
