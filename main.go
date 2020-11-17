package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/crabtree/sm-xsuaa-poc/internal/params"
)

func main() {
	params, err := params.Parse()
	exitOnError(err)

	httpClient := &http.Client{}
	httpClient.Transport = &BasicAuthTransport{
		Username: params.Username,
		Password: params.Password,
		Rt:       http.DefaultTransport,
	}

	sm := smclient.NewClient(context.Background(), httpClient, params.BaseURL)
	marketplace, err := sm.Marketplace(&query.Parameters{}) //ListBrokers(&query.Parameters{})
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, offering := range marketplace.ServiceOfferings {
		fmt.Printf("name: %s, brokerID: %s, ID: %s\n", offering.Name, offering.BrokerID, offering.ID)
	}

}

func exitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type BasicAuthTransport struct {
	Username string
	Password string

	Rt http.RoundTripper
}

func (b *BasicAuthTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if b.Username != "" && b.Password != "" {
		request.SetBasicAuth(b.Username, b.Password)
	}

	return b.Rt.RoundTrip(request)
}
