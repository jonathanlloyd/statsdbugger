package statsd

type Metric interface {
	Name() string
	Tags() map[string]string
}

type Counter struct {
	CName       string
	CValue      float64
	CSampleRate float64
	CTags       map[string]string
}

func (c Counter) Name() string {
	return c.CName
}

func (c Counter) Tags() map[string]string {
	return c.CTags
}

type Gauge struct {
	GName  string
	GValue float64
	GTags  map[string]string
}

func (g Gauge) Name() string {
	return g.GName
}

func (g Gauge) Tags() map[string]string {
	return g.GTags
}
