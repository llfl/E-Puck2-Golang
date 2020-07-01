package driverutils

//Options Options
type Options struct{
	I2CDevices  string
	I2CAddress  int
	GyroAddress int
	Gyro bool
}

//Option option func
type Option func(*Options)

func newOptions(opt ...Option) Options {
	opts := Options{}
	
	for _, o := range opt{
		o(&opts)
	}

	if len(opts.I2CDevices) == 0 {
		opts.I2CDevices = DefaultI2CDevices
	}

	if opts.I2CAddress == 0 {
		opts.I2CAddress = DefaultI2CAddress
	}

	if opts.GyroAddress == 0 {
		opts.GyroAddress = DefaultGyroAddress
	}

	return opts
}

//I2CDevice set I2CDevice port
func I2CDevice(s string) Option {
	return func (o *Options) {
		o.I2CDevices = s
	}
}

//I2CAddress set I2C address
func I2CAddress(s int) Option {
	return func (o *Options) {
		o.I2CAddress = s
	}
}

//GyroAddress set I2C address
func GyroAddress(s int) Option {
	return func (o *Options) {
		o.GyroAddress = s
	}
}

//EnableGyro enable gyro
func EnableGyro() Option {
	return func (o *Options) {
		o.Gyro = true
	}
}