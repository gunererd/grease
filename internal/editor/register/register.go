package register

type Register struct {
	data string
}

func NewRegister() *Register {
	return &Register{data: ""}
}

func (r *Register) Set(text string) {
	r.data = text
}

func (r *Register) Get() string {
	return r.data
}
