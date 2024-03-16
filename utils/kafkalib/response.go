package kafkalib

type ResponseInterface interface {
	Commit()
}

type Response struct {
	commit bool
}

func (w *Response) Commit() {
	w.commit = true
}

func NewWriter() ResponseInterface {
	return &Response{
		commit: false,
	}
}
