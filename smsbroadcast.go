package smsBroadcast

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// A Client manages communication with the smsBroadcast API.
type Broadcast struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.  Defaults to https://www.yourdomain.com/billing/, but can be
	// set to a domain endpoint to use with your billing at your enterprise.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL  *url.URL
	username string
	password string
}

// NewClient returns a new smsBroadcast API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(user string, pass string, httpClient *http.Client) *Broadcast {
	if httpClient == nil {
		httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	}
	baseURL, _ := url.Parse("https://api.smsbroadcast.com.au/api-adv.php")

	b := &Broadcast{client: httpClient, BaseURL: baseURL, username: user, password: pass}

	return b
}

// A WRequest manages communication with the smsBroadcast API.
type WRequest struct {
	data *url.Values
	url  *url.URL
}

// Response is a smsBroadcast API response.  This wraps the standard http.Response
// returned from smsBroadcast and provides convenient access to things like
// pagination links.
type Response struct {
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	Body          string
	ContentLength int64
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	body, _ := ioutil.ReadAll(r.Body)
	response := &Response{
		Status:        r.Status,
		StatusCode:    r.StatusCode,
		Body:          string(body),
		ContentLength: r.ContentLength,
	}
	return response
}

// newRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (b *Broadcast) newRequest(dat map[string]string, action string) (*WRequest, error) {
	rel, err := url.Parse("")
	if err != nil {
		return nil, err
	}
	u := b.BaseURL.ResolveReference(rel)

	if len(strings.TrimSpace(action)) > 0 {
		dat["username"] = b.username
		dat["password"] = b.password
	}

	return &WRequest{url: u, data: addFormValues(dat)}, nil
}

// addFormValues adds the parameters in opt as URL values parameters.
func addFormValues(opt map[string]string) *url.Values {
	uv := url.Values{}
	for k, v := range opt {
		uv.Set(k, v)
	}
	return &uv
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func checkResponse(r *Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	return errors.New(r.Body)
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (b *Broadcast) doRequest(req WRequest, v interface{}) (*Response, error) {

	resp, err := b.client.PostForm(req.url.String(), *req.data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)
	err = checkResponse(response)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	return response, err
}

// Params specifies the optional parameters to various List methods that
// support pagination.
type Params struct {
	parms map[string]string
	u     string
}

func do(b *Broadcast, p Params, a interface{}) (*Response, error) {
	req, err := b.newRequest(p.parms, p.u)
	if err != nil {
		return nil, err
	}

	resp, err := b.doRequest(*req, a)
	if err != nil {
		return resp, err
	}

	return resp, err
}

type Message struct {
	To       string // The receiving mobile number(s)
	From     string // The sender ID for the messages. Can be a mobile number or letters, up to 11 characters and should not contain punctuation or spaces. Leave blank to use SMS Broadcast’s 2-way number.
	Message  string // The content of the SMS message. Must not be longer than 160 characters unless the maxsplit parameter is used. Must be URL encoded.
	Ref      string // Your reference number for the message to help you track the message status. This parameter is optional and can be up to 20 characters.
	MaxSplit int    // Determines the maximum length of your SMS message
	Delay    int    // Number of minutes to delay the message. Use this to schedule messages for later delivery.
}

type Output struct {
	Status          string // Will show if your messages have been accepted by the API
	ReceivingNumber string // The receiving mobile number
	Message         string // Will display our reference number for the SMS message, or the reason for or Error Message a failed SMS message.
}

// Send message function can be used to send a single.
func (b *Broadcast) Send(msg Message) (*Output, *Response, error) {
	a := new(Output)
	parms := map[string]string{"message": msg.Message, "to": msg.To, "from": msg.From}

	resp, err := do(b, Params{parms: parms, u: "api-adv.php"}, a)
	if err != nil {
		return nil, resp, err
	}

	// TODO Work out output
	// log.Print([]byte(resp.Body))

	return a, resp, err
}
