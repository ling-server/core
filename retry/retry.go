package retry

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jpillora/backoff"
	"github.com/ling-server/core/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	ErrorRetryTimeout = errors.New("retry timeout")
)

type abort struct {
	cause error
}

func (a *abort) Error() string {
	if a.cause != nil {
		return fmt.Sprintf("retry abort, error: %v", a.cause)
	}

	return "retry abort"
}

// Abort wrap error to stop the Retry function
func Abort(err error) error {
	return &abort{cause: err}
}

// Options for the retry functions
type Options struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Timeout         time.Duration
	Callback        func(err error, sleep time.Duration)
	Backoff         bool
}

type Option func(*Options)

func InitialInterval(initial time.Duration) Option {
	return func(opts *Options) {
		opts.InitialInterval = initial
	}
}

func MaxInterval(max time.Duration) Option {
	return func(opts *Options) {
		opts.MaxInterval = max
	}
}

func Timeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func Callback(callback func(err error, sleep time.Duration)) Option {
	return func(opts *Options) {
		opts.Callback = callback
	}
}

func Backoff(backoff bool) Option {
	return func(opts *Options) {
		opts.Backoff = backoff
	}
}

func Retry(f func() error, options ...Option) error {
	opts := &Options{
		Backoff: true,
	}

	for _, o := range options {
		o(opts)
	}

	if opts.InitialInterval <= 0 {
		opts.InitialInterval = time.Millisecond * 100
	}

	if opts.MaxInterval <= 0 {
		opts.MaxInterval = time.Second
	}

	if opts.Timeout <= 0 {
		opts.Timeout = time.Minute
	}

	var b *backoff.Backoff

	if opts.Backoff {
		b = &backoff.Backoff{
			Min:    opts.InitialInterval,
			Max:    opts.MaxInterval,
			Factor: 2,
			Jitter: true,
		}
	}

	var err error
	timeout := time.After(opts.Timeout)
	for {
		select {
		case <-timeout:
			return errors.New(ErrorRetryTimeout)
		default:
			err = f()
			if err == nil {
				return nil
			}

			var ab *abort
			if errors.As(err, &ab) {
				return ab.cause
			}

			var sleep time.Duration
			if opts.Backoff {
				sleep = b.Duration()
			}

			if opts.Callback != nil {
				opts.Callback(err, sleep)
			}

			time.Sleep(sleep)
		}
	}
}
