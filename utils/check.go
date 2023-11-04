package utils

import (
	"bytes"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strconv"
	"strings"
)

var ErrorTimes = 0

type tcpingRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Check() (bool, error) {
	var ip string
	if Conf.GetBool(`GetIP`) || Conf.GetString(`SpecifiedHost`) == "" {
		raw, err := GetIP()
		if err != nil {
			return false, err
		}
		ip = strings.TrimSpace(raw)
	} else {
		ip = Conf.GetString(`SpecifiedHost`)
	}

	times, ok := 0, false
	var final error
	for times < 5 {
		resp, err := http.Get(Conf.GetString("URL") + `?ip=` + ip + `&port=443`)
		if err != nil {
			times++
			final = err
			continue
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			times++
			final = err
			continue
		}

		var ret tcpingRes
		err = jsoniter.Unmarshal(buf.Bytes(), &ret)
		if err != nil {
			times++
			final = err
			continue
		}

		ok, err = strconv.ParseBool(ret.Status)
		if err != nil {
			times++
			final = err
			continue
		}
		break
	}

	if times < 5 {
		return ok, nil
	}

	return false, final
}

func GetIP() (string, error) {
	times, ip := 0, ""
	var final error

	for times < 5 {
		resp, err := RequestIPSB()
		if err != nil {
			times++
			final = err
			continue
		}

		if resp.StatusCode != 200 {
			times++
			final = errors.New(resp.Status)
			continue
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			times++
			final = err
			continue
		}

		ip = buf.String()
		break
	}

	if times < 5 {
		return ip, nil
	}

	return "", final
}

func RequestIPSB() (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api-ipv4.ip.sb/ip", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	return client.Do(req)
}
