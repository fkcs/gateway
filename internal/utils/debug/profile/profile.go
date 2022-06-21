// Package profile is for profilers
package profile

type Profile interface {
	// Start the profiler
	Start() error
	// Stop the profiler
	Stop() error
	// Name of the profiler
	String() string
}

var (
	DefaultProfile Profile = new(ProfileBase)
)

type ProfileBase struct{}

func (p *ProfileBase) Start() error {
	return nil
}

func (p *ProfileBase) Stop() error {
	return nil
}

func (p *ProfileBase) String() string {
	return "profile-base"
}

type Options struct {
	Name string
}

type Option func(o *Options)

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
