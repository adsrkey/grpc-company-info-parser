package dto

import "github.com/adsrkey/grpc-company-info-parser/parser/parserpb"

type Info struct {
	ID          string
	Inn         string
	CeoName     string
	Kpp         string
	CompanyName string
}

func (i *Info) ToParserResponse() *parserpb.ParserResponse {
	return &parserpb.ParserResponse{
		Company: &parserpb.Company{
			Inn:         i.Inn,
			Kpp:         i.Kpp,
			CompanyName: i.CompanyName,
		},
		Supervisor: &parserpb.Supervisor{
			Name: i.CeoName,
		},
	}
}
