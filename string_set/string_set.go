package string_set

type Collection struct {
	items map[string]bool
}

func (c *Collection) Add(v string) {
	if c.items == nil {
		c.items = make(map[string]bool)
	}
	c.items[v] = true
}

func (c *Collection) Delete(v string) {
	if c.items == nil {
		c.items = make(map[string]bool)
	}
	delete(c.items, v)
}

func (c Collection) Exists(v string) bool {
	if c.items == nil {
		return false
	}
	_, ok := c.items[v]
	return ok
}

func (c Collection) Union(o Collection) Collection {
	out := make(map[string]bool)
	for k, _ := range c.items {
		out[k] = true
	}
	for k, _ := range o.items {
		out[k] = true
	}
	c.items = out
	return c
}

func (c Collection) Subtract(o Collection) Collection {
	out := make(map[string]bool)
	for k, _ := range c.items {
		out[k] = true
	}
	for k, _ := range o.items {
		delete(out, k)
	}
	c.items = out
	return c
}

func (c Collection) Intersection(o Collection) Collection {
	out := make(map[string]bool)
	for k, _ := range c.items {
		if o.Exists(k) {
			out[k] = true
		}
	}
	c.items = out
	return c
}

func (c Collection) ToSlice() (out []string) {
	out = make([]string, 0, len(c.items))
	for s, _ := range c.items {
		out = append(out, s)
	}
	return
}
