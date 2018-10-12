package requests

func (r *Request) ReqHook(option ReqOption) {
	r.reqOptions = append(r.reqOptions, option)
}

func (r *Request) RespHook(option RespOption) {
	r.respOptions = append(r.respOptions, option)
}
