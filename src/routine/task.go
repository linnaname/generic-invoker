package routine

type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}
