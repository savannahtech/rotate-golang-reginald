package client

import "github.com/osquery/osquery-go/gen/osquery"

type Client interface {
	Query(query string) (*osquery.ExtensionResponse, error)
	Close()
}
