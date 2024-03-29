package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"
	"go.uber.org/zap"
)

const (
	timeout    = 10 * time.Second
	apiVersion = "externaldata.gatekeeper.sh/v1beta1"
)

var log logr.Logger

func main() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("unable to initialize logger: %v", err))
	}
	log = zapr.NewLogger(zapLog)
	log.Info("starting server...")

	// load Gatekeeper's CA certificate
	caCert, err := os.ReadFile("/tmp/gatekeeper/ca.crt")
	if err != nil {
		log.Info("error reading file" + err.Error())
		panic(err)
	}

	clientCAs := x509.NewCertPool()
	clientCAs.AppendCertsFromPEM(caCert)

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", processTimeout(validate, timeout))

	log.Info("initialing server...")

	server := &http.Server{
		Addr:              ":8090",
		Handler:           mux,
		ReadHeaderTimeout: timeout,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  clientCAs,
			MinVersion: tls.VersionTLS13,
		},
	}

	log.Info("ListenAndServeTLS")
	if err := server.ListenAndServeTLS("/etc/ssl/certs/server.crt", "/etc/ssl/certs/server.key"); err != nil {
		log.Error(err, "ListenAndServeTLS - Error")
		panic(err)
	}
}

// A Response struct to map the Entire Response
type Response struct {
	Vulnerabilities []string `json:"vulnerabilities"`
	CertifyBads     []string `json:"certifyBads"`
	SlsaList        []string `json:"SlsaList"`
	SbomList        []string `json:"SbomList"`
}

func validate(w http.ResponseWriter, req *http.Request) {
	// only accept POST requests
	if req.Method != http.MethodPost {
		//sendResponse(nil, "only POST is allowed", w)
		return
	}

	// read request body
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		//sendResponse(nil, fmt.Sprintf("unable to read request body: %v", err), w)
		return
	}

	// parse request body
	var providerRequest externaldata.ProviderRequest
	err = json.Unmarshal(requestBody, &providerRequest)
	if err != nil {
		//sendResponse(nil, fmt.Sprintf("unable to unmarshal request body: %v", err), w)
		return
	}

	results := make([]externaldata.Item, 0)
	// iterate over all keys
	for _, key := range providerRequest.Request.Keys {
		// Providers should add a caching mechanism to avoid extra calls to external data sources.
		log.Info("Image received: " + key)

		splitImage := strings.Split(key, "@")

		if len(splitImage) != 2 {
			log.Info("split failed")
			//sendResponse(nil, "split key failed", w)
			return
		}
		log.Info("Image Digest: " + splitImage[1])
		response, err := http.Get("http://rest-api.default.svc.cluster.local:8081/query/artInfo/" + url.QueryEscape(splitImage[1]))
		if err != nil {
			log.Error(err, "response error")
			//sendResponse(nil, fmt.Sprintf("ERROR getting response from GUAC: %v", err), w)
			return
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Error(err, "response read error")
			//sendResponse(nil, fmt.Sprintf("ERROR reading response data: %v", err), w)
			return
		}

		var responseObject Response
		json.Unmarshal(responseData, &responseObject)

		if len(responseObject.Vulnerabilities) > 0 {
			log.Info("found Vulnerabilities: " + strings.Join(responseObject.Vulnerabilities, ","))
			results = append(results, externaldata.Item{
				Key:   key,
				Value: len(responseObject.Vulnerabilities),
			})
		} else if len(responseObject.CertifyBads) > 0 {
			log.Info("found CertifyBad: " + responseObject.CertifyBads[0])
			results = append(results, externaldata.Item{
				Key:   key,
				Value: len(responseObject.CertifyBads),
			})
		} else if len(responseObject.SbomList) == 0 {
			log.Info("SBOM not found for image: " + key)
			results = append(results, externaldata.Item{
				Key:   key,
				Value: 1,
			})
		} else if len(responseObject.SlsaList) == 0 {
			log.Info("SLSA not found for image: " + key)
			results = append(results, externaldata.Item{
				Key:   key,
				Value: 1,
			})
		} else {
			log.Info("guac verified image: " + key)
			// results = append(results, externaldata.Item{
			// 	Key:   key,
			// 	Value: -1,
			// })
		}
		sendResponse(&results, "", w)
	}
}

// sendResponse sends back the response to Gatekeeper.
func sendResponse(results *[]externaldata.Item, systemErr string, w http.ResponseWriter) {
	response := externaldata.ProviderResponse{
		APIVersion: apiVersion,
		Kind:       "ProviderResponse",
		Response: externaldata.Response{
			Idempotent: true,
		},
	}

	if results != nil {
		response.Response.Items = *results
	} else {
		response.Response.SystemError = systemErr
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(fmt.Errorf("failed to encode: %w", err))
	}
}

func processTimeout(h http.HandlerFunc, duration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), duration)
		defer cancel()

		r = r.WithContext(ctx)

		processDone := make(chan bool)
		go func() {
			h(w, r)
			processDone <- true
		}()

		select {
		case <-ctx.Done():
			sendResponse(nil, "operation timed out", w)
		case <-processDone:
		}
	}
}
