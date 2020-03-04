package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"strings"
)

func GetValuesFromConfig(name string) []string {
	values := beego.AppConfig.String(name)
	if values == "" {
		log.Printf("'%v' env variable is empty.", name)
		return nil
	}

	s := strings.Split(values, ",")
	log.Printf("Fetched data from %v env variable: %v", name, s)
	return s
}

func GetVersionOrNil(value, postfix string) *string {
	if value == "" {
		return nil
	}

	if postfix == "" {
		v := value
		return &v
	}

	v := fmt.Sprintf("%v-%v", value, postfix)

	return &v
}

func EncodeStructToBytes(s interface{}) ([]byte, error) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(s)
	if err != nil {
		return nil, err
	}
	return reqBodyBytes.Bytes(), nil
}