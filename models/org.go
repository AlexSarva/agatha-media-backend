package models

type INNRequest struct {
	INN string `json:"inn"`
}

type Organization struct {
	OGRN      string `json:"ogrn" ch:"ogrn"`
	INN       string `json:"inn" ch:"inn"`
	KPP       string `json:"kpp" ch:"kpp"`
	ShortName string `json:"short_name" ch:"short_name"`
	FullName  string `json:"full_name" ch:"full_name"`
	RegDate   string `json:"reg_date" ch:"reg_date"`
	EndDate   string `json:"end_date" ch:"end_date"`
	OKVED     string `json:"okved_id" ch:"okved_id"`
	Capital   string `json:"capital" ch:"capital"`
	RegionID  int8   `json:"region_id" ch:"region_id"`
	Address   string `json:"address" ch:"address"`
}

type Test struct {
	INN string `json:"inn" ch:"inn"`
}
