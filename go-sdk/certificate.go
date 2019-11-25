package go_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type CertificateCreateRequest struct {
	Name  string `json:"name"`
	Cert  string `json:"cert"`
	Key   string `json:"key"`
	Chain string `json:"chain,omitempty"`
}

type CertificateCreateResponse struct {
	Status string `json:"status"`
	Data   CertificateCreateResponseData
}

type CertificateCreateResponseData struct {
	Name   string `json:"name"`
	CertId string `json:"cert_id"`
}

type CertificateGetResponse struct {
	Status string `json:"status"`
	Data   CertificateGetResponseData
}

type CertificateGetResponseData struct {
	Name        string `json:"name"`
	CreatedDate string `json:"created_date"`
}

func (c *NorainaApiClient) CreateCertificate(iCert CertificateCreateRequest) (*CertificateCreateResponse, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	err := appendFileToWriter("cert", iCert.Cert, w)
	if err != nil {
		return nil, err
	}

	err = appendFileToWriter("key", iCert.Key, w)
	if err != nil {
		return nil, err
	}

	if iCert.Chain != "" {
		err = appendFileToWriter("chain", iCert.Chain, w)
		if err != nil {
			return nil, err
		}
	}

	err = w.WriteField("name", iCert.Name)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", norainaDomain+certificateRoute, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-access-token", c.Token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Create Certificate Error with Status code %v", resp.StatusCode)
	}

	certificateResponse := &CertificateCreateResponse{}
	err = json.NewDecoder(resp.Body).Decode(certificateResponse)
	if err != nil {
		return nil, err
	}

	return certificateResponse, nil
}

func (c *NorainaApiClient) GetCertificate(id string) (*CertificateGetResponse, error) {
	body, err := c.getResource(id, certificateRoute)
	if err != nil {
		return nil, err
	}

	defer body.Close()

	certificateResponse := &CertificateGetResponse{}
	err = json.NewDecoder(body).Decode(certificateResponse)
	if err != nil {
		return nil, err
	}

	return certificateResponse, nil
}

func (c *NorainaApiClient) DeleteCertificate(id string) error {
	return c.deleteResource(id, certificateRoute)
}

func appendFileToWriter(name string, path string, w *multipart.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := w.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}
