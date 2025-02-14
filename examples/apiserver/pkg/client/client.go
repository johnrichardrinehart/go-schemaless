package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"encoding/json"

	"github.com/rbastic/go-schemaless/examples/apiserver/pkg/api"
	"github.com/rbastic/go-schemaless/models"

	"io/ioutil"
	"net/http"
)

const contentTypeJSON = "application/json"

type Client struct {
	Address string
}

func New() *Client {
	return &Client{}
}

func (c *Client) WithAddress(addr string) *Client {
	c.Address = addr
	return c
}

func (c *Client) Get(ctx context.Context, tblName string, rowKey string, columnKey string, refKey int64) (cell models.Cell, found bool, err error) {
	postURL := c.Address + "/api/get"

	// TODO: make the context part of the request

	var getRequest api.GetRequest
	getRequest.Table = tblName
	getRequest.RowKey = rowKey
	getRequest.ColumnKey = columnKey
	getRequest.RefKey = refKey

	getRequestMarshal, err := json.Marshal(getRequest)
	if err != nil {
		return models.Cell{}, false, err
	}

	request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(getRequestMarshal))
	if err != nil {
		return models.Cell{}, false, err
	}
	request.Header.Set("Content-Type", contentTypeJSON)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return models.Cell{}, false, err
	}
	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return models.Cell{}, false, err
	}
	var gr api.GetResponse

	err = json.Unmarshal(responseBody, &gr)
	if err != nil {
		return models.Cell{}, false, err
	}
	if gr.Error != "" {
		return models.Cell{}, false, errors.New(gr.Error)
	}

	return *gr.Cell, gr.Found, nil
}

func (c *Client) GetLatest(ctx context.Context, tblName string, rowKey string, columnKey string) (cell models.Cell, found bool, err error) {
	postURL := c.Address + "/api/getLatest"

	var getLatestRequest api.GetRequest
	getLatestRequest.Table = tblName
	getLatestRequest.RowKey = rowKey
	getLatestRequest.ColumnKey = columnKey

	getLatestRequestMarshal, err := json.Marshal(getLatestRequest)
	if err != nil {
		return models.Cell{}, false, err
	}

	request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(getLatestRequestMarshal))
	if err != nil {
		return models.Cell{}, false, err
	}
	request.Header.Set("Content-Type", contentTypeJSON)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return models.Cell{}, false, err
	}
	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return models.Cell{}, false, err
	}

	var glr api.GetLatestResponse

	fmt.Printf("(client) GETLATEST RESPONSEBODY:'%s'\n", responseBody)

	err = json.Unmarshal(responseBody, &glr)
	if err != nil {
		return models.Cell{}, false, err
	}
	if glr.Error != "" {
		return models.Cell{}, false, errors.New(glr.Error)
	}

	return *glr.Cell, glr.Found, nil
}

func (c *Client) PartitionRead(ctx context.Context, tblName string, partitionNumber int, location string, value uint64, limit int) (cells []models.Cell, found bool, err error) {
	postURL := c.Address + "/api/partitionRead"

	// TODO: add context

	var partitionReadRequest api.PartitionReadRequest
	partitionReadRequest.Table = tblName
	partitionReadRequest.PartitionNumber = partitionNumber
	partitionReadRequest.Location = location
	partitionReadRequest.Value = value
	partitionReadRequest.Limit = limit

	partitionReadRequestMarshal, err := json.Marshal(partitionReadRequest)
	if err != nil {
		return nil, false, err
	}

	request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(partitionReadRequestMarshal))
	if err != nil {
		return nil, false, err
	}
	request.Header.Set("Content-Type", contentTypeJSON)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, false, err
	}
	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, false, err
	}

	var prr api.PartitionReadResponse
	err = json.Unmarshal(responseBody, &prr)
	if err != nil {
		return nil, false, err
	}

	if prr.Error != "" {
		return nil, false, errors.New(prr.Error)
	}

	return prr.Cells, prr.Found, nil
}

func (c *Client) Put(ctx context.Context, tblName string, rowKey string, columnKey string, refKey int64, body string) (*api.PutResponse, error) {
	postURL := c.Address + "/api/put"

	// TODO: make the context part of the request

	var putRequest api.PutRequest
	putRequest.Table = tblName
	putRequest.RowKey = rowKey
	putRequest.ColumnKey = columnKey
	putRequest.RefKey = refKey
	putRequest.Body = body

	putRequestMarshal, err := json.Marshal(putRequest)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(putRequestMarshal))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentTypeJSON)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	pr := new(api.PutResponse)
	err = json.Unmarshal(responseBody, pr)
	return pr, err
}
