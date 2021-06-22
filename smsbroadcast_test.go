package smsBroadcast

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestBroadcast_Send(t *testing.T) {
	type fields struct {
		StatusCode int
		Body       string
	}
	type args struct {
		msg Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Output
		want1   *Response
		wantErr bool
	}{
		{
			name: "Send OK",
			fields: fields{
				StatusCode: 200,
				Body:       `OK:61111111111:1002248045`,
			},
			args: args{
				Message{
					To:      "61111111111",
					Message: "SDK Test",
					From:    "From a Long Name",
				},
			},
			want: &Output{
				Status:          "OK",
				ReceivingNumber: "61111111111",
				Message:         "1002248045",
			},
			want1: &Response{
				StatusCode:    200,
				Body:          `OK:61111111111:1002248045`,
				ContentLength: 0,
			},
		},
		{
			name: "Long Sender",
			fields: fields{
				StatusCode: 200,
				Body:       `ERROR: The sender cannot be longer than 11 characters.`,
			},
			args: args{
				Message{
					To:      "61111111111",
					Message: "SDK Test",
					From:    "From a Long Name",
				},
			},
			want: &Output{
				Status:  "ERROR",
				Message: "ERROR: The sender cannot be longer than 11 characters.",
			},
			want1: &Response{
				StatusCode:    200,
				Body:          `ERROR: The sender cannot be longer than 11 characters.`,
				ContentLength: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testClient := NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: tt.fields.StatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(tt.fields.Body)),
					Header:     make(http.Header),
				}
			})

			u, _ := url.Parse("")
			b := &Broadcast{
				client:  testClient,
				BaseURL: u,
			}
			got, got1, err := b.Send(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Broadcast.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Broadcast.Send() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Broadcast.Send() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
