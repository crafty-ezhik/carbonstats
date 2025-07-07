package req

import (
	"github.com/crafty-ezhik/carbonstats/pkg/res"
	"net/http"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.JSON(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	err = IsValid[T](body)
	if err != nil {
		res.JSON(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}
