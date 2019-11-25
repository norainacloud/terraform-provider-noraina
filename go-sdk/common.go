package go_sdk

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *NorainaApiClient) getResource(id string, resourceRoute string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, norainaDomain+resourceRoute+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-access-token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] Get %v/%v with Status code %v", resourceRoute, id, resp.StatusCode)
	}

	log.Printf("[DEBUG] Get %v/%v -> response Status : %s", resourceRoute, id, resp.Status)
	log.Printf("[DEBUG] Get %v/%v -> response Headers : %s", resourceRoute, id, resp.Header)
	log.Printf("[DEBUG] Get %v/%v -> response Headers : %v", resourceRoute, id, resp.Body)

	return resp.Body, nil
}

func (c *NorainaApiClient) deleteResource(id string, resourceRoute string) error {
	req, err := http.NewRequest(http.MethodDelete, norainaDomain+resourceRoute+"/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Add("x-access-token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Print(fmt.Println("[DEBUG] DELETE -> response Status : ", resp.Status))
	log.Print(fmt.Println("[DEBUG] DELETE -> response Headers : ", resp.Header))
	log.Print(fmt.Println("[DEBUG] DELETE -> response Body : ", string(respBody)))

	return nil
}
