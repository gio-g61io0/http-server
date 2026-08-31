package client

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"personal-http-server/parser/request"
	"personal-http-server/response"
	"strconv"
)

type Client struct {
	buffer   []byte
	request  *request.Request   //The request parser
	response *response.Response //the response instance
	conn     net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:     conn,
		request:  request.NewRequest(conn),
		response: nil,
		buffer:   make([]byte, 2048),
	}
}

func (c *Client) ReadRequest() {
	read, err := c.request.ReadRequest(c.buffer)

	if err != nil {
		log.Println("Something went wrong parsing the request")
		return
	}
	log.Printf("Read %d of bytes from request", read)
	c.request.InspectRequest()

	data := bytes.NewReader([]byte("Hello world. Gwapo si Gio"))

	res, err := response.NewResponse(c.conn, response.HttpVersion(), response.BuildStatusLine(response.StatusCode(response.OK)), response.BuildBody(data))

	res.HeaderSet("content-length", strconv.Itoa(len("Hello world. Gwapo si Gio")))
	res.HeaderSet("Connection", "close")
	c.response = &res
}

func (c *Client) WriteResponse() (int, error) {
	if c.response == nil {
		return 0, fmt.Errorf("Cannot write response! No response set")
	}
	c.conn.Write(c.response.ResponseLine())
	c.conn.Write([]byte(c.response.WriteHeader()))
	c.conn.Write(c.response.Body)

	return 0, nil

}
