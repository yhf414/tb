package soap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var licence = ""

type CertSoapIn struct {
	XMLName      xml.Name `xml:"https://ws.nciic.org.cn/nciic_ws/services/NciicServices nciicCheck"`
	Action       string   `xml:"-"`
	InLicense    string   `xml:"inLicense"`
	InConditions string   `xml:"inConditions"`
}

func (si CertSoapIn) GetAction() string {
	return si.Action
}

type CertSoapOut struct {
	XMLName xml.Name `xml:"https://api.nciic.org.cn/NciicServices nciicCheckResponse"`
	Return  string   `xml:"out"`
}

func NewCertSoapIn(licence string, name string, cid string) CertSoapIn {
	si := CertSoapIn{}
	si.Action = ""
	si.InLicense = licence
	si.InConditions = fmt.Sprintf(`<ROWS><INFO><SBM>中国工商银行</SBM></INFO><ROW><GMSFHM>公民身份号码</GMSFHM><XM>姓名</XM></ROW><ROW FSD="100600" YWLX="个人贷款"><GMSFHM>%s</GMSFHM><XM>%s</XM></ROW></ROWS>`, cid, name)
	return si
}

type AuthService struct {
	Url string
}

func NewAuthService() *AuthService {
	s := &AuthService{}
	s.Url = "https://ws.nciic.org.cn/nciic_ws/services/NciicServices"

	return s
}

//
func (s AuthService) nciicCheck(si CertSoapIn) (r string, err error) {
	sr, err := CallService(si, s.Url)
	if err != nil {
		return "", err
	}

	// monta a estrutura de retorno
	var so CertSoapOut
	err = xml.Unmarshal([]byte(sr.Body.Content), &so)
	if err != nil {
		return "", err
	}

	return so.Return, nil
}

type SF_Item struct {
	Gmsfhm     string `xml:"gmsfhm"`
	Res_xm     string `xml:"result_xm"`
	Res_gmsfhm string `xml:"result_gmsfhm"`
	Error      string `xml:"errormesage"`
}

func Certify(name string, cid string) (bool, error) {
	s := NewAuthService()
	if licence == "" {
		f, _ := os.Open("licence.txt")
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return false, err
		}
		licence = string(buf)
	}
	soapin := NewCertSoapIn(licence, name, cid)
	res, err := s.nciicCheck(soapin)
	if err != nil {
		return false, err
	}
	resXml := struct {
		XMLName xml.Name `xml:ROWS`
		Row     struct {
			XMLName   xml.Name `xml:"ROW"`
			ErrorCode int      `xml:"ErrorCode"`
			ErrorMsg  string   `xml:"ErrorMsg"`
			No        int      `xml:"no,attr"`
			Input     struct {
				XMLName xml.Name `xml:"INPUT"`
				Gmsfhm  string   `xml:"gmsfhm"`
				Xm      string   `xml:"xm"`
			}
			OutPut struct {
				XMLName xml.Name  `xml:"OUTPUT"`
				Items   []SF_Item `xml:"ITEM"`
			}
		}
	}{}
	err = xml.Unmarshal([]byte(res), &resXml)
	if err != nil {
		return false, err
	}
	if len(resXml.Row.OutPut.Items) < 2 {
		errRes := struct {
			XMLName xml.Name `xml:RESPONSE`
			Rows    struct {
				XMLName xml.Name `xml:"ROWS"`
				Row     struct {
					XMLName   xml.Name `xml:"ROW"`
					ErrorCode int      `xml:"ErrorCode"`
					ErrorMsg  string   `xml:"ErrorMsg"`
				}
			}
		}{}
		err = xml.Unmarshal([]byte(res), &errRes)
		if errRes.Rows.Row.ErrorCode != 0 {
			return false, errors.New(fmt.Sprintf("code:%d, msg:%s", errRes.Rows.Row.ErrorCode, errRes.Rows.Row.ErrorMsg))
		}
		return false, errors.New(res)
	}
	for _, v := range resXml.Row.OutPut.Items {
		if v.Error != "" {
			return false, errors.New(v.Error)
		} else if v.Res_gmsfhm == "不一致" || v.Res_xm == "不一致" {
			return false, errors.New("姓名与身份证号码不一致")
		}

	}
	return true, nil

}
