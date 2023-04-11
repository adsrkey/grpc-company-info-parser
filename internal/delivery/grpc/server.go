package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adsrkey/grpc-company-info-parser/internal/delivery/grpc/dto"
	"github.com/adsrkey/grpc-company-info-parser/internal/util/parser"
	"github.com/adsrkey/grpc-company-info-parser/parser/parserpb"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Server struct {
	parserpb.UnimplementedParserServiceServer
}

func (s *Server) Get(ctx context.Context, req *parserpb.ParserRequest) (*parserpb.ParserResponse, error) {
	_, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
	defer cancel()

	log.Println("get request with:", req.String())
	inn := req.GetInn()

	client := &http.Client{}
	mainGetResp, err := makeMainReq(client)
	if err != nil {
		return nil, err
	}
	defer mainGetResp.Body.Close()

	cookie := mainGetResp.Header.Get("Cookie")

	ajaxGetReq, err := makeAjaxReq(inn, cookie)
	if err != nil {
		return nil, err
	}

	ajaxResp, err := client.Do(ajaxGetReq)
	if err != nil {
		return nil, err
	}

	info := &dto.Info{}
	info.Inn = inn
	err = bindAjaxData(ajaxResp.Body, info)
	if err != nil {
		return nil, err
	}

	if info.ID == "" {
		return nil, errors.New("link is empty")
	}

	infoGetReq, err := makeInfoGetReq(info.ID, cookie)
	if err != nil {
		return nil, err
	}

	infoResp, err := client.Do(infoGetReq)
	if err != nil {
		return &parserpb.ParserResponse{}, err
	}
	defer infoResp.Body.Close()

	html, err := bindInfoHtml(infoResp.Body)
	if err != nil {
		return nil, err
	}

	p := parser.New(info, html)
	err = p.Parse()
	if err != nil {
		return nil, err
	}

	return info.ToParserResponse(), nil
}

func makeMainReq(client *http.Client) (*http.Response, error) {
	return client.Get("https://www.rusprofile.ru")
}

func generateAjaxUrl(inn string) string {
	cacheKey := fmt.Sprintf("%0.16f", rand.Float64())
	return fmt.Sprintf("https://www.rusprofile.ru/ajax.php?query=%s&action=search&cacheKey=%s", inn, cacheKey)
}

func makeAjaxReq(inn string, cookie string) (*http.Request, error) {
	url := generateAjaxUrl(inn)
	ajaxGetReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ajaxGetReq.Header.Set("Cookie", cookie)
	return ajaxGetReq, nil
}

func makeInfoGetReq(id string, cookie string) (*http.Request, error) {
	idUrl := fmt.Sprintf("https://www.rusprofile.ru%s", id)

	infoGetReq, err := http.NewRequest(http.MethodGet, idUrl, nil)
	if err != nil {
		return nil, err
	}

	infoGetReq.Header.Set("Cookie", cookie)

	return infoGetReq, nil
}

func bindAjaxData(body io.ReadCloser, info *dto.Info) error {
	defer body.Close()

	var ajaxData Data
	err := json.NewDecoder(body).Decode(&ajaxData)
	if err != nil {
		return err
	}

	if len(ajaxData.Ul) == 1 {
		info.ID = ajaxData.Ul[0].Link
		info.CeoName = ajaxData.Ul[0].CeoName
	}

	return nil
}

func bindInfoHtml(body io.ReadCloser) (string, error) {
	htmlBytes, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(htmlBytes), nil
}
