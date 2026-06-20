package Client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	CMD any `json:"cmd"`
	Key any `json:"key,omitempty"`
	Val any `json:"value,omitempty"`
}

type Client struct {
	Port           int
	ConnectionType string
	ConnObj        net.Conn
}

func (c *Client) Connect() (net.Conn, error) {
	conn, err := net.Dial(c.ConnectionType, fmt.Sprintf(":%v", c.Port))
	c.ConnObj = conn
	if err != nil {
		return nil, err
	}
	return c.ConnObj, nil
}

func (c *Client) Disconnect() error {
	err := c.ConnObj.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Set(key string, value string) error {
	newRequest := Request{"SET", key, value}

	jsonRequest, err := json.Marshal(newRequest)

	if err != nil {
		fmt.Println("error while marhsalling ")
	}

	_, err = c.ConnObj.Write(append(jsonRequest, '\n'))

	if err != nil {
		fmt.Println("error while sending SET request")
		return err
	}

	response, err := bufio.NewReader(c.ConnObj).ReadString('\n')
	if err != nil {
		fmt.Println("Read error: ", err)
		return err
	}
	fmt.Print(response)
	return nil
}

func (c *Client) Get(key string) (string, error) {

	newRequest := Request{"GET", key, nil}

	jsonRequest, err := json.Marshal(newRequest)
	if err != nil {
		fmt.Println("error while marhsalling ")
		return "", err
	}

	_, err = c.ConnObj.Write(append(jsonRequest, '\n'))

	if err != nil {
		fmt.Println("error while sending GET request")
		return "", err
	}

	serverReader := bufio.NewReader(c.ConnObj)

	response, err := serverReader.ReadString('\n')

	if err != nil {
		fmt.Println("Read error")
		return "", err
	}

	return response, nil
}

func (c *Client) Compact() error {
	newRequest := Request{"COMPACT", nil, nil}

	jsonRequest, _ := json.Marshal(newRequest)

	_, err := c.ConnObj.Write(append(jsonRequest, '\n'))

	if err != nil {
		fmt.Println("error while sending COMPACT request")
		return err
	}

	serverReader := bufio.NewReader(c.ConnObj)
	response, err := serverReader.ReadString('\n')

	if err != nil {
		fmt.Println("read error")
		return err
	}

	fmt.Println(response)
	return nil
}
