package runners

import (
	"math"

	"barf/internal/cmd"
	"barf/internal/op"
	"barf/internal/rsync"
	"barf/internal/typeconv"
)

type copyRunner struct {
	operation     *op.Operation
	rsync         *rsync.Rsync
	stdoutHandler cmd.LogHandler
	stderrHandler cmd.LogHandler
	statusHandler statushandler
}

func (r *copyRunner) init(operation *op.Operation) {
	r.operation = operation
	r.rsync = rsync.NewRsync()
	r.rsync.OnStdout(r.handleStdout)
	r.rsync.OnStderr(r.handleStderr)
	r.rsync.OnStatus(r.handleStatus)
}

func (r *copyRunner) Start() {
	args := []string{}
	fromArray, _ := typeconv.ToArray(r.operation.Args["from"])
	from := typeconv.ToStringArray(fromArray)

	for _, value := range from {
		args = append(args, value)
	}

	args = append(args, r.operation.Args["to"].(string))

	r.rsync.Copy(args)
}

func (r *copyRunner) Abort() error {
	return r.rsync.Abort()
}

func (r *copyRunner) OperationID() op.OperationID {
	return r.operation.ID
}

func (r *copyRunner) OnStdout(handler cmd.LogHandler) {
	r.stdoutHandler = handler
}

func (r *copyRunner) OnStderr(handler cmd.LogHandler) {
	r.stderrHandler = handler
}

func (r *copyRunner) OnStatus(handler statushandler) {
	r.statusHandler = handler
}

func (r *copyRunner) handleStdout(line string) {
	r.stdoutHandler(line)
}

func (r *copyRunner) handleStderr(line string) {
	r.stderrHandler(line)
}

func (r *copyRunner) handleStatus(status *rsync.RsyncStatus) {
	r.statusHandler(&op.OperationStatus{
		Message:        status.Message,
		BytesDiffTotal: status.BytesDiffTotal,
		BytesTotal:     status.BytesTotal,
		BytesDone:      status.BytesDoneTotal,
		Progress:       math.Round(status.Progress*100) / 100,
		Speed:          math.Round(status.Speed*100) / 100,
		FilesDiffTotal: status.FilesDiffTotal,
		FilesTotal:     status.FilesTotal,
		FilesDone:      status.CurrentFileIndex,
		SecondsLeft:    status.SecondsLeft,
		FileName:       status.CurrentFileName,
		Finished:       status.Finished,
		ExitCode:       status.ExitCode,
	})
}
