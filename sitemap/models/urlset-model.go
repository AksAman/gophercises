package models

type URLSet struct {
	URLs []URL `json:"url" xml:"url"`
}

type URL struct {
	Loc        string  `json:"loc" xml:"loc"`
	Lastmod    *string `json:"lastmod,omitempty" xml:"lastmod,omitempty"`
	Changefreq *string `json:"changefreq,omitempty" xml:"changefreq,omitempty"`
	Priority   float64 `json:"priority" xml:"priority"`
	// exclude depth from xml
	Depth int `json:"depth" xml:"-"`
}

func NewURL(location string, depth int) *URL {
	return &URL{
		Loc:      location,
		Priority: 0.5,
		Depth:    depth,
	}
}

func NewURLWithPriority(location string, priority float64, depth int) *URL {
	return &URL{
		Loc:      location,
		Priority: priority,
		Depth:    depth,
	}
}

// sort.Interface implementation
// by Priority
type URLsByPriority []URL

func (a URLsByPriority) Len() int           { return len(a) }
func (a URLsByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a URLsByPriority) Less(i, j int) bool { return a[i].Priority < a[j].Priority }

// sort.Interface implementation
// by Depth

type URLsByDepth []URL

func (a URLsByDepth) Len() int           { return len(a) }
func (a URLsByDepth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a URLsByDepth) Less(i, j int) bool { return a[i].Depth < a[j].Depth }

// sort.Interface implementation
// by DepthAndPriority

type URLsByDepthAndPriority []URL

func (a URLsByDepthAndPriority) Len() int      { return len(a) }
func (a URLsByDepthAndPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a URLsByDepthAndPriority) Less(i, j int) bool {
	if a[i].Depth == a[j].Depth {
		return a[i].Priority < a[j].Priority
	}
	return a[i].Depth < a[j].Depth
}
