package noraina

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	certificateRoute = "api/certificate"
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

func (c *Client) CreateCertificate(ctx context.Context, iCert CertificateCreateRequest) (*CertificateCreateResponse, error) {
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

	req, err := http.NewRequest(http.MethodPost, defaultBaseURL+certificateRoute, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-access-token", c.Token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	certificateResponse := new(CertificateCreateResponse)
	err = c.Do(ctx, req, certificateResponse)
	if err != nil {
		return nil, err
	}

	return certificateResponse, nil
}

func (c *Client) GetCertificate(ctx context.Context, id string) (*CertificateGetResponse, error) {
	path := fmt.Sprintf("%s/%s", certificateRoute, id)
	req, err := c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	certificateResponse := new(CertificateGetResponse)
	err = c.Do(ctx, req, certificateResponse)
	if err != nil {
		return nil, err
	}

	return certificateResponse, nil
}

func (c *Client) DeleteCertificate(ctx context.Context, id string) error {
	return c.deleteResource(ctx, id, certificateRoute)
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
