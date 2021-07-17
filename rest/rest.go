package rest

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
This package includes a lot of the basic handling for the application to have
both versioning and health checks
*/

const (
	HEADER_CONTENT_TYPE = "Content-Type"

	CONTENT_TYPE_JSON = "application/json"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}

type ResponseWrapper struct {
	Status    int         `json:"status"`
	RequestId string      `json:request_id`
	Version   string      `json:"version,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Hostname  string      `json:"host"`
}

func NewResponseWrapperWithVersion(requestId string, version string) ResponseWrapper {
	return ResponseWrapper{
		Status:    http.StatusOK,
		RequestId: requestId,
		Version:   version,
		Hostname:  hostname,
	}
}

func NewResponseWrapperWithError(requestId string, err error) ResponseWrapper {
	return ResponseWrapper{
		Status:    http.StatusNotFound,
		RequestId: requestId,
		Error:     err.Error(),
		Hostname:  hostname,
	}
}

func NewResponseWrapper(requestId string) ResponseWrapper {
	return ResponseWrapper{
		Status:    http.StatusOK,
		RequestId: requestId,
		Hostname:  hostname,
	}
}

func NewResponseWrapperWithInternalError(requestId string, err error) ResponseWrapper {
	return ResponseWrapper{
		Status:    http.StatusInternalServerError,
		RequestId: requestId,
		Error:     err.Error(),
		Hostname:  hostname,
	}
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(NotFound)
}

func NotFound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add(HEADER_CONTENT_TYPE, CONTENT_TYPE_JSON)
	response.WriteHeader(http.StatusNotFound)

	err := fmt.Errorf("%s not found", request.URL.String())

	encoder := json.NewEncoder(response)
	encErr := encoder.Encode(NewResponseWrapperWithError("-", err))
	if encErr != nil {
		log.Panicln(errors.Wrap(encErr, "Error encoding JSON response"))
	}

	response.Write([]byte("\n"))

	log.Println(err)
}

func HealthCheckEndpoint(version string, quitTime time.Duration) func(http.ResponseWriter, *http.Request) {
	healthStatus := http.StatusOK

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(
			ch,
			syscall.SIGINT,
			syscall.SIGHUP,
			syscall.SIGTERM,
			syscall.SIGUSR1,
		)

		drainTime := quitTime * time.Second

		sig := <-ch
		fmt.Printf("[INFO] Got a %s, quiting in %s. \n", sig, drainTime)
		<-time.After(drainTime)

		os.Exit(0)
	}()

	return func(response http.ResponseWriter, request *http.Request) {
		pingResponse := NewResponseWrapperWithVersion("ping", version)

		msg, err := json.Marshal(pingResponse)
		if err != nil {
			panic(err)
		}

		response.Header().Add(HEADER_CONTENT_TYPE, CONTENT_TYPE_JSON)
		response.WriteHeader(healthStatus)
		response.Write(msg)

		if healthStatus == http.StatusOK {
			log.Println("PING - In Service")
		} else {
			log.Println("PING - Going Out Of Service")
		}
	}
}

func VersionEndpoint(version string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		versionResponse := NewResponseWrapperWithVersion("version", version)

		msg, err := json.Marshal(versionResponse)
		if err != nil {
			panic(err)
		}

		response.Write(msg)
		log.Printf(string(msg))
	}
}
