package stella

import "stella/sync"

type Ctx struct {
	values *sync.Map[string, any]
}

func NewCtx() *Ctx {
	return &Ctx{
		values: sync.NewMap[string, any](),
	}
}

func (self *Ctx) Values() map[string]any {
	return self.values.Map()
}

func (self *Ctx) Set(name string, value any) {
	self.values.Set(name, value)
}

func (self *Ctx) Get(name string) any {
	if !self.values.Has(name) {
		return nil
	}

	return self.values.Get(name)
}
