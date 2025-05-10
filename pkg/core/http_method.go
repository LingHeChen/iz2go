package iz2go

type Post struct{}

func (p *Post) GetMethod() string {
	return "POST"
}

type Get struct{}

func (g *Get) GetMethod() string {
	return "GET"
}

type Put struct{}

func (p *Put) GetMethod() string {
	return "PUT"
}

type Delete struct{}

func (d *Delete) GetMethod() string {
	return "DELETE"
}

type Patch struct{}

func (p *Patch) GetMethod() string {
	return "PATCH"
}

type Options struct{}

func (o *Options) GetMethod() string {
	return "OPTIONS"
}

type Head struct{}

func (h *Head) GetMethod() string {
	return "HEAD"
}

type Trace struct{}

func (t *Trace) GetMethod() string {
	return "TRACE"
}

type Connect struct{}

func (c *Connect) GetMethod() string {
	return "CONNECT"
}
