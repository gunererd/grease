package handler

var defaultRegister = &Register{data: ""}

type Register struct {
	data string
}

func (r *Register) Set(text string) {
	r.data = text
}

func (r *Register) Get() string {
	return r.data
}
