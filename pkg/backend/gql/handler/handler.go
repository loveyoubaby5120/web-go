package handler

// Forked from graphql-go/handler.

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"fmt"
	"jdy/pkg/util/errs"
)

const (
	contentTypeJSON           = "application/json"
	contentTypeGraphQL        = "application/graphql"
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	contentTypeDataForm       = "multipart/form-data"
)

type RequestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

// a workaround for getting`variables` as a JSON string
type requestOptionsCompatibility struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

func getFromForm(values url.Values) (*RequestOptions, error) {
	query := values.Get("query")
	if query == "" {
		return nil, errs.InvalidArgument.New("no query")
	}
	// get variables map
	var variables map[string]interface{}
	variablesStr := values.Get("variables")
	if variablesStr != "" {
		err := json.Unmarshal([]byte(variablesStr), &variables)
		if err != nil {
			return nil, errs.InvalidArgument.New("invalid variables json: %s, err: %v", variablesStr, err)
		}
	}

	return &RequestOptions{
		Query:         query,
		Variables:     variables,
		OperationName: values.Get("operationName"),
	}, nil
}

// NewRequestOptions Parses a http.Request into GraphQL request options struct
func NewRequestOptions(r *http.Request) (*RequestOptions, error) {
	if r.URL.Query().Get("query") != "" {
		return getFromForm(r.URL.Query())
	}
	if r.Method == "GET" {
		return nil, errs.InvalidArgument.New("missing query")
	}

	if r.Method != "POST" {
		return nil, errs.InvalidArgument.New("not POST or GET")
	}

	if r.Body == nil {
		return nil, errs.InvalidArgument.New("no body")
	}

	// TODO: improve Content-Type handling
	contentTypeStr := r.Header.Get("Content-Type")
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := contentTypeTokens[0]

	switch contentType {
	case contentTypeGraphQL:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errs.InvalidArgument.Wrap(err)
		}
		return &RequestOptions{
			Query: string(body),
		}, nil
	case contentTypeFormURLEncoded:
		if err := r.ParseForm(); err != nil {
			return nil, errs.InvalidArgument.Wrap(err)
		}

		return getFromForm(r.PostForm)
	case contentTypeDataForm:
		if err := r.ParseMultipartForm(64 * 1024); err != nil {
			return nil, err
		}
		return getFromForm(r.PostForm)
	case contentTypeJSON:
		fallthrough
	default:
		var opts RequestOptions
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errs.InvalidArgument.Wrap(err)
		}
		// fmt.Println(format.ByteToString(body))
		// fmt.Println(&opts)
		err = json.Unmarshal(body, &opts)
		fmt.Println(err)
		if err != nil {
			// Probably `variables` was sent as a string instead of an object.
			// So, we try to be polite and try to parse that as a JSON string
			var optsCompatible requestOptionsCompatibility
			err = json.Unmarshal(body, &optsCompatible)
			if err != nil {
				return nil, errs.InvalidArgument.New("invalid json")
			}
			err = json.Unmarshal([]byte(optsCompatible.Variables), &opts.Variables)
			if err != nil {
				return nil, errs.InvalidArgument.New("invalid json")
			}
		}
		// fmt.Println(&opts)
		return &opts, nil
	}
}
