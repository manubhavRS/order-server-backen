package Health

import (
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
	w.WriteHeader(http.StatusOK)
	return
}
