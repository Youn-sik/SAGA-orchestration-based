package workflow

import (
	"encoding/json"
	"net/url"
	"saga/saga"
	"saga/utils"
)

func Execute(qs url.Values, body []byte) error {
	// saga 객체 만들기
	s := saga.NewSaga(utils.GenerateTransactionID(), 5)
	err := json.Unmarshal(body, &s)
	if err != nil {
		utils.FatalIfError(err)
		return err
	}

	// saga 객체의 Validate 메서드 호출
	err = s.Validate()
	if err != nil {
		return err
	}

	// saga 객체의 Run 메서드 호출 및 결과 응답 제공
	return s.Run()
}
