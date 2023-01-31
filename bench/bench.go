package bench

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type client struct {
	c   *http.Client
	url string
	key string
}

type Event struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
}

func makeHTTPRequest[Req, Resp any](ctx context.Context, c *client, method string, url string, key string, reqBody Req) (Resp, error) {
	httpRequestBody, err := json.Marshal(reqBody)
	if err != nil {
		return *new(Resp), err
	}
	reader := bytes.NewReader(httpRequestBody)
	httpRequest, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return *new(Resp), err
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("u21-key", key)

	httpResponse, err := c.c.Do(httpRequest)
	if err != nil {
		return *new(Resp), err
	}
	defer httpResponse.Body.Close()
	respBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return *new(Resp), err
	}
	if httpResponse.StatusCode != 200 {
		return *new(Resp), fmt.Errorf("http status code: %d, body: %s", httpResponse.StatusCode, respBody)
	}

	var temp Resp
	err = json.Unmarshal(respBody, &temp)
	if err != nil {
		return *new(Resp), err
	}
	return temp, nil
}

func (c *client) GetEventGeneric(ctx context.Context, id int) (Event, error) {
	return makeHTTPRequest[any, Event](ctx, c, http.MethodGet, c.url+"/"+strconv.Itoa(id), c.key, nil)
}

func (c *client) CreateEventGeneric(ctx context.Context, e Event) (Event, error) {
	return makeHTTPRequest[Event, Event](ctx, c, http.MethodPost, c.url+"/events", c.key, e)
}

func (c *client) GetEvent(ctx context.Context, id int) (Event, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/"+strconv.Itoa(id), nil)
	if err != nil {
		return Event{}, fmt.Errorf("req: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("u21-key", c.key)

	resp, err := c.c.Do(req)
	if err != nil {
		return Event{}, fmt.Errorf("do: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Event{}, fmt.Errorf("readall: %v", err)
	}

	if resp.StatusCode != 200 {
		return Event{}, fmt.Errorf("code: %v, body: %v", resp.StatusCode, string(body))
	}

	var e Event
	if err := json.Unmarshal(body, &e); err != nil {
		return Event{}, fmt.Errorf("unmarshal: %v", err)
	}
	return e, nil
}

func (c *client) CreateEvent(ctx context.Context, e Event) (Event, error) {
	reqBody, err := json.Marshal(e)
	if err != nil {
		return Event{}, fmt.Errorf("marshal: %v", err)
	}
	reader := bytes.NewReader(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/events", reader)
	if err != nil {
		return Event{}, fmt.Errorf("req: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("u21-key", c.key)

	resp, err := c.c.Do(req)
	if err != nil {
		return Event{}, fmt.Errorf("do: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Event{}, fmt.Errorf("readall: %v", err)
	}

	if resp.StatusCode != 200 {
		return Event{}, fmt.Errorf("code: %v, body: %v", resp.StatusCode, string(body))
	}

	var ev Event
	if err := json.Unmarshal(body, &ev); err != nil {
		return Event{}, fmt.Errorf("unmarshal: %v", err)
	}
	return ev, nil
}
