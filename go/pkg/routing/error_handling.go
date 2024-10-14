// func write_error()
package routing

import (
	"encoding/json"
	"io"
)

type HttpErrorBody struct {
	Err string `json:"error"`
}

func WriteHttpError(writer io.Writer, errorMessage string) error {
	body := HttpErrorBody{
		Err: errorMessage,
	}

	serializedError, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = writer.Write(serializedError)
	if err != nil {
		return err
	}

	return nil
}
