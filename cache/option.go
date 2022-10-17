package cache

import "time"

type Options struct {
	Address    string
	Codec      Codec
	Expiration time.Duration
	Prefix     string
}

type Option func(*Options)

func (opts *Options) Key(key string) string {
	return opts.Prefix + key
}

func newOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// Address sets the address
func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// Expiration sets the default expiration
func Expiration(d time.Duration) Option {
	return func(o *Options) {
		o.Expiration = d
	}
}

// Prefix sets the prefix
func Prefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}
