package publishers

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/mchmarny/thingz-commons"
	ws "golang.org/x/net/websocket"
)

const (
	origin = "http://localhost/"
)

// NewWebsocketPublisher parses connection string to websocket server
// and returned a configured version of the publisher
//    wss://app.domain.com:4443/path, 123456789=
func NewWebsocketPublisher(connStr string) (Publisher, error) {

	parts := strings.Split(connStr, ",")

	if len(parts) != 2 {
		log.Printf("Expected 2 parts > %v", parts)
		return nil, errors.New("Invalid URL format. Expected 2 parts")
	}

	config, err := ws.NewConfig(strings.TrimSpace(parts[0]),
		origin)
	if err != nil {
		return nil, err
	}
	config.Header.Add("Authorization", "Bearer "+strings.TrimSpace(parts[1]))

	conn, err := ws.DialConfig(config)

	/*
		conn, err := proxyDial(strings.TrimSpace(parts[0]),
			origin,
			strings.TrimSpace(parts[1]))
		if err != nil {
			log.Printf("Error on dial: %v", err)
			return nil, err
		}
	*/

	p := WebsocketPublisher{Connection: conn}

	return p, nil
}

// WebsocketPublisher
type WebsocketPublisher struct {
	Connection *ws.Conn
}

// String
func (m *WebsocketPublisher) String() string {
	return "Websocket Publisher"
}

// Publish
func (p WebsocketPublisher) Publish(in <-chan *commons.Metric, err chan<- error) {
	go func() {
		for {
			select {
			case msg := <-in:

				// convert
				content, convErr := msg.ToBytes()
				if convErr != nil {
					log.Printf("Error on convert: %v", convErr)
					err <- convErr
				}

				// send
				_, sendErr := p.Connection.Write(content)
				if sendErr != nil {
					log.Printf("Error on write: %v", sendErr)
					err <- sendErr
				}

			} // select
		} // for
	}() // go
}

// Finalize
func (p WebsocketPublisher) Finalize() {
	log.Println("Websocket publisher is done")
}

func proxyDial(url_, origin, token string) (c *ws.Conn, err error) {
	proxyUrl := os.Getenv("HTTP_PROXY")
	if len(proxyUrl) == 0 {
		return ws.Dial(url_, "", origin)
	}

	log.Printf("Using proxy: %s", proxyUrl)
	purl, err := url.Parse(os.Getenv("HTTP_PROXY"))
	if err != nil {
		return nil, err
	}

	config, err := ws.NewConfig(url_, origin)
	if err != nil {
		return nil, err
	}
	config.Header.Add("Authorization", "Bearer "+token)

	client, err := httpConnect(purl.Host, url_)
	if err != nil {
		return nil, err
	}

	return ws.NewClient(config, client)
}

func httpConnect(proxy, url_ string) (io.ReadWriteCloser, error) {

	log.Printf("Dialing proxy: %s", proxy)
	p, err := net.Dial("tcp", proxy)
	if err != nil {
		return nil, err
	}

	log.Printf("Parsing URL: %s", url_)
	turl, err := url.Parse(url_)
	if err != nil {
		return nil, err
	}

	req := http.Request{
		Method: "CONNECT",
		URL:    &url.URL{},
		Host:   turl.Host,
	}

	cc := httputil.NewProxyClientConn(p, nil)
	cc.Do(&req)
	if err != nil && err != httputil.ErrPersistEOF {
		return nil, err
	}

	rwc, _ := cc.Hijack()

	return rwc, nil
}
