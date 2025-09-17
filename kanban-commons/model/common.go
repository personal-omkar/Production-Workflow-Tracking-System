package model

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type RawQuery struct {
	Host  string
	Port  string
	Type  string
	Query string
}

type PaginationReq struct {
	Type     string           `json:"Type"`
	Limit    string           `json:"Limit"`
	PageNo   int              `json:"Pageno"`
	Order    string           `json:"Order"`
	Criteria []map[string]any `json:"Criteria"`
}

type PaginationResp struct {
	TotalNo int `json:"TotalNo"`
	Page    int `json:"Page"`
	Offset  int `josn:"Offset"`
}

func (r *RawQuery) RawQry(result interface{}) error {
	jsonValue, err := json.Marshal(r)
	if err != nil {
		log.Printf("RawQry: Error marshaling raw query: %v", err)
		return err
	}

	resp, err := http.Post("http://"+r.Host+":"+r.Port+"/RawQuery", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("RawQry: Error making POST request: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("RawQry: Error reading response body: %v", err)
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		slog.Error("RawQry: Error decoding response body", "error", err)
		return err
	}
	return nil
}
