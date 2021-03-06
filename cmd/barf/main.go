package main

import (
	"fmt"
	"os"

	"barf/internal/cli"
	"barf/internal/cli/actions"
	"barf/internal/cli/ui"
	"barf/internal/config"
	"barf/internal/daemon"

	mcli "github.com/jawher/mow.cli"
)

func wrapAction(action cli.Action, args map[string]interface{}) {
	err := action(args)

	if err != nil {
		fmt.Println(err)
		os.Exit(255)
	}
}

func main() {
	daemon.CheckDaemon()

	app := mcli.App(config.Name, config.Description)
	app.Version("v version", fmt.Sprintf("%s\n%s\n%s", config.Version, config.BuildChecksum, config.BuildTime))

	width := app.IntOpt("width, w", 0, "terminal width to use, if not set (or zero) it will be auto detected")

	app.Before = func() {
		ui.SetWidth(*width)
	}

	app.Action = func() {
		wrapAction(actions.Monitor, map[string]interface{}{})
	}

	app.Command("list l", "list active operations", func(cmd *mcli.Cmd) {
		cmd.Action = func() {
			wrapAction(actions.List, map[string]interface{}{})
		}
	})

	app.Command("monitor m", "monitors active operations", func(cmd *mcli.Cmd) {
		cmd.LongDesc = "monitors active operations, if ids are given it will exit when those operations have finished"
		cmd.Spec = "[IDS...]"
		ids := cmd.StringsArg("IDS", nil, "IDs to monitor")

		cmd.Action = func() {
			wrapAction(actions.Monitor, map[string]interface{}{
				"ids": ids,
			})
		}
	})

	app.Command("abort a", "aborts an active operation", func(cmd *mcli.Cmd) {
		cmd.Spec = "ID"
		id := cmd.StringArg("ID", "", "ID to abort")

		cmd.Action = func() {
			wrapAction(actions.Abort, map[string]interface{}{
				"id": id,
			})
		}
	})

	app.Command("copy cp", "copies files or directories", func(cmd *mcli.Cmd) {
		cmd.Spec = "SRC... DST"
		src := cmd.StringsArg("SRC", nil, "Source files or directories to copy")
		dst := cmd.StringArg("DST", "", "Destination to copy to")

		cmd.Action = func() {
			wrapAction(actions.Copy, map[string]interface{}{
				"src": src,
				"dst": dst,
			})
		}
	})

	app.Command("move mv", "moves files or directories", func(cmd *mcli.Cmd) {
		cmd.Spec = "SRC... DST"
		src := cmd.StringsArg("SRC", nil, "Source files or directories to move")
		dst := cmd.StringArg("DST", "", "Destination to move to")

		cmd.Action = func() {
			wrapAction(actions.Move, map[string]interface{}{
				"src": src,
				"dst": dst,
			})
		}
	})

	app.Command("push", "mirrors source directory in destination directory", func(cmd *mcli.Cmd) {
		cmd.Spec = "SRC DST"
		src := cmd.StringArg("SRC", "", "Source directory to read from")
		dst := cmd.StringArg("DST", "", "Destination directory to update")

		cmd.Action = func() {
			wrapAction(actions.Push, map[string]interface{}{
				"src": src,
				"dst": dst,
			})
		}
	})

	app.Command("pull", "mirrors dst directory in src directory", func(cmd *mcli.Cmd) {
		cmd.Spec = "SRC DST"
		src := cmd.StringArg("SRC", "", "Source directory to update")
		dst := cmd.StringArg("DST", "", "Destination directory to to read from")

		cmd.Action = func() {
			wrapAction(actions.Pull, map[string]interface{}{
				"src": src,
				"dst": dst,
			})
		}
	})

	if !config.IsProduction() {
		app.Command("dummy", "starts dummy operations", func(cmd *mcli.Cmd) {
			cmd.Spec = "[ITER]"
			iterations := cmd.StringArg("ITER", "10", "Iterations to run")

			cmd.Action = func() {
				wrapAction(actions.Dummy, map[string]interface{}{
					"iterations": iterations,
				})
			}
		})
	}

	app.Command("stop s", "stop background process", func(cmd *mcli.Cmd) {
		cmd.Action = func() {
			wrapAction(actions.Stop, map[string]interface{}{})
		}
	})

	app.Command("update u", "check for updates", func(cmd *mcli.Cmd) {
		cmd.Action = func() {
			wrapAction(actions.Update, map[string]interface{}{})
		}
	})

	app.Run(os.Args)
}
