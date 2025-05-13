package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

// TODO: Implement the leveler interface for Verbose & ExtraVerbose.
type CmdLineHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *CmdLineHandler) Handle(ctx context.Context, r slog.Record) error {
	// TODO: Decide on debug/verbose message highlighting.
	var preMessage string
	if r.Level == slog.LevelDebug {
		preMessage = "DEBUG: "
	}

	// Extract Attrs into a string.
	fields := "" // TODO: Clean up trailing space for message with no Attrs.
	r.Attrs(func(a slog.Attr) bool {
		fields = fields + " " +a.Key + ": " + a.Value.String()
		// fields[a.Key] = a.Value.Any()
		return true
	})

	// Colour debug messages in the event verbose is set in config.
	if r.Level == slog.LevelDebug {
		h.l.Println(color.YellowString(preMessage + r.Message + fields))
	} else {
		h.l.Println(preMessage + r.Message + fields)
	}

	return nil
}

func NewCmdLineHandler(out io.Writer, opts *slog.HandlerOptions) *CmdLineHandler {
	return &CmdLineHandler{
		Handler: slog.NewTextHandler(out, opts),
		l:       log.New(out, "", 0),
	}
}

func SetOutputLevel(verbose bool) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if verbose {
		opts.Level = slog.LevelDebug
	}
	handler := NewCmdLineHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
