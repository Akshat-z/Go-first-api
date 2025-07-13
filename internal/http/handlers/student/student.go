package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Akshat-z/student-api/internal/storage"
	"github.com/Akshat-z/student-api/internal/types"
	"github.com/Akshat-z/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func Create(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) { //? EOF is no input no perameter passed in body so get error in decoder.
			response.WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		//+ valaditing the req body by make required in student struct
		if err := validator.New().Struct(student); err != nil {
			validationErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationErr))
			return
		}
		creationId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusAccepted, map[string]int64{"Id": creationId})
	}
}

//_ sql injection
