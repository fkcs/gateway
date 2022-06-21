package cache

import "time"

type Options struct {
	Expiration time.Duration
	Items      map[string]Item
}

type Option func(o *Options)

func Expiration(d time.Duration) Option {
	return func(o *Options) {
		o.Expiration = d
	}
}

func Items(i map[string]Item) Option {
	return func(o *Options) {
		o.Items = i
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Expiration: DefaultExpiration,
		Items:      make(map[string]Item),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
