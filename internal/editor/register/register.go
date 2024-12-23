package register

import "log"

type Register struct {
	data string
}

func NewRegister() *Register {
	return &Register{data: ""}
}

func (r *Register) Set(text string) {
	r.data = text
	log.Printf("Save to register: '%s'\n", r.data)
}

func (r *Register) Get() string {
	log.Printf("Get from register: '%s'\n", r.data)
	return r.data
}
