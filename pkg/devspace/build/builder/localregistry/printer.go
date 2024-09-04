package localregistry

import (
	"context"
	"github.com/containerd/console"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/moby/buildkit/util/progress/progresswriter"

	"io"
)

type printer struct {
	status chan *client.SolveStatus
	done   <-chan struct{}
	err    error
}

func (p *printer) Done() <-chan struct{} {
	return p.done
}

func (p *printer) Err() error {
	return p.err
}

func (p *printer) Status() chan *client.SolveStatus {
	if p == nil {
		return nil
	}
	return p.status
}

func NewPrinter(ctx context.Context, out io.Writer) (progresswriter.Writer, error) {
	statusCh := make(chan *client.SolveStatus)
	doneCh := make(chan struct{})
	pw := &printer{
		status: statusCh,
		done:   doneCh,
	}
	cons := console.Current()
	go func() {
		_, err := progressui.DisplaySolveStatus(ctx, cons, out, statusCh)
		if err != nil {
			pw.err = err
		}
		close(doneCh)
	}()
	return pw, nil
}
