package redis

type Option struct {
	Host     string
	Port     int
	DB       int
	Password string
}

func New(options ...func(*Option)) *Option {
	svr := &Option{}
	for _, o := range options {
		o(svr)
	}
	return svr
}

func NewHost(host string) func(*Option) {
	return func(o *Option) {
		o.Host = host
	}
}

func NewPort(port int) func(*Option) {
	return func(o *Option) {
		o.Port = port
	}
}

func NewDB(db int) func(*Option) {
	return func(o *Option) {
		o.DB = db
	}
}

func NewPassword(password string) func(*Option) {
	return func(o *Option) {
		o.Password = password
	}
}
