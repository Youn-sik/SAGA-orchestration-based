package saga

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"saga/utils"
	"time"
)

type Saga struct {
	Gid                      int64                    `validate:"required"`
	Requests                 []map[string]interface{} `json:"requests" validate:"required"`
	RequestsCompensation     []map[string]interface{} `json:"requests_compensation" validate:"required"`
	HttpRequestTimeoutSecond int
	HttpClient               *http.Client
}

var v = validator.New()

func NewSaga(gid int64, httpRequestTimeoutSecond int) *Saga {
	return &Saga{
		Gid:                      gid,
		HttpRequestTimeoutSecond: httpRequestTimeoutSecond,
		HttpClient: &http.Client{
			Timeout: time.Duration(httpRequestTimeoutSecond) * time.Second,
		},
	}
}

func (s *Saga) Validate() error {
	err := v.Struct(s)
	if err != nil {
		utils.FatalIfError(err)
		return err
	}
	if len(s.Requests) < 1 || len(s.RequestsCompensation) < 1 {
		err := errors.New("invalid length requests or requests_compensation")
		utils.FatalIfError(err)
		return err
	}
	return nil
}

func (s *Saga) Run() error {
	utils.PrintInfof("========== Starting request ==========")

	for idx, req := range s.Requests {
		for url, body := range req {
			if err := s.reqHttp(url, body); err != nil {
				utils.PrintInfof("Request failed at request %s: %v", url, err)
				return s.compensate(idx)
			}
		}
	}

	utils.PrintInfof("All requests succeeded")
	return nil
}

func (s *Saga) compensate(lastSuccessIdx int) error {
	utils.PrintInfof("Starting compensation from step %d", lastSuccessIdx)

	for i := lastSuccessIdx; i >= 0; i-- {
		for url, body := range s.RequestsCompensation[i] {
			if err := s.reqHttp(url, body); err != nil {
				utils.PrintInfof("Compensation request failed at request %s: %v", url, err)
				// 정책에 따라 [compensation 실패 시]
				// 1. return err => 실패 보상 중간에 실패 시, 가능한 과정까지 실패보상 요청.
				// 2. continue => 실패 보상 중간에 실패해도, 전체 과정에 대해 실패보상 요청.
				continue
			}
		}
	}

	utils.PrintInfof("Compensation requests completed")
	return nil
}

func (s *Saga) reqHttp(url string, body interface{}) error {

	jsonBody, err := json.Marshal(body)
	if err != nil {
		utils.PrintInfof("-> REQUEST url: %s | body: %+v", url, body)
		utils.PrintErrorf("Request body parsing failed at request %s: %v", url, err)
		return err
	}

	utils.PrintInfof("-> REQUEST url: %s | body: %s", url, string(jsonBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		utils.PrintErrorf("Request initialize failed at request %s: %v", url, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		utils.PrintErrorf("Request failed at request %s: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.PrintErrorf("Response body parsing failed at request %s: %v", url, err)
		return err
	}

	utils.PrintInfof("<- RESPONSE [%s] url: %s | body: %s", resp.Status, url, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(bodyBytes))
	}

	return nil

}
