package dto

import "net/http"

type Request struct {
	BPM int
	RPM int
	w   http.ResponseWriter
	r   *http.Request
}

func NewRequest(BPM, RPM int, w http.ResponseWriter, r *http.Request) *Request {
	return &Request{
		BPM: BPM,
		RPM: RPM,
		w:   w,
		r:   r,
	}
}
