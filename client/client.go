package client

import (
	"net"
	"personal-http-server/parser/request"
	"personal-http-server/response"
)

type Client struct {
	request  *request.Request   //The request parser
	response *response.Response //the response instance
	conn     net.Conn
}

func NewClient(conn net.Conn) *Client {
	response, err := response.NewResponse(conn)
	if err != nil {
		return nil
	}

	return &Client{
		conn:     conn,
		request:  request.NewRequest(conn),
		response: &response,
	}
}

func (c *Client) ReadRequest() {

}
