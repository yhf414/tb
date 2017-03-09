package soap

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type SoapIn interface {
	GetAction() string
}

func CallService(si SoapIn, url string) (sr *Envelope, err error) {
	// cria o soap envelope
	se := NewEnvelope()

	// gerar o conteúdo do corpo em xml
	bsi, err := xml.Marshal(&si)
	if err != nil {
		return nil, err
	}
	// associa o corpo da requisição
	se.Body.Content = string(bsi)

	// gerar o xml da requisição
	bse, err := xml.Marshal(&se)
	if err != nil {
		return nil, err
	}

	// cria um reader para o corpo da requisição
	br := strings.NewReader(string(bse))

	// cria a requisição
	req, err := http.NewRequest("POST", url, br)
	if err != nil {
		return nil, err
	}

	// adiciona os cabeçalhos http necessários do soap
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", si.GetAction())
	client := &http.Client{}
	if strings.HasPrefix(url, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// executa a requisição
		client = &http.Client{Transport: tr}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// le o conteudo do retorno
	bsr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// gerar a estrutura de retorno
	err = xml.Unmarshal(bsr, &sr)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 500 {
		var sf Fault
		err = xml.Unmarshal([]byte(sr.Body.Content), &sf)
		if err != nil {
			return nil, errors.New(resp.Status)
		}
		return nil, errors.New(sf.FaultString)
	}
	return sr, nil
}
