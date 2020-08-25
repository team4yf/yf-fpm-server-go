package utils

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	tnet "github.com/toolkits/net"
)

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

//CheckErr panic if err is not nil
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GetLocalIP 获取本地内网IP
func GetLocalIP() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.1"
		}
	})
	return clientIP
}

//JSON2String convert the json object to string
func JSON2String(j interface{}) (str string) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return "{}"
	}
	str = (string)(bytes)
	return
}

//Interface2Struct convert the json object to struct
func Interface2Struct(j interface{}, dest interface{}) (err error) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, dest); err != nil {
		return
	}
	return
}

//StringToStruct convert the string to struct
func StringToStruct(data string, desc interface{}) (err error) {
	if err = json.Unmarshal(([]byte)(data), desc); err != nil {
		return
	}
	return
}

// GenShortID 生成一个id
func GenShortID() string {
	sid, _ := shortid.Generate()
	return sid
}

// GenUUID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
