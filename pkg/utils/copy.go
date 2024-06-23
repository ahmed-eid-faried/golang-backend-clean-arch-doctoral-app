package utils

import (
	"encoding/json"
	"fmt"
)

func Copy(dest interface{}, src interface{}) {
	data, err := json.Marshal(src)
	fmt.Print(data)

	if err != nil {
		fmt.Print(err)
		return
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		fmt.Print(err)
		return
	}

}
