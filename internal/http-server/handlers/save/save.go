package save

type Request struct {
	Genre string `json:"genre"`
	Name  string `json:"name"`
	Size  string `json:"size"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	// Name   string `json:"name"`
}

// вынести в отдельную библиотеку response и const, так как они будут одинаковые во всех ответах handlerов
const (
	StatusOK="OK"
	StatusError="Error"
)

type PictureSaver interface{
	SavePicture(r *Request) (int64,error)
}

