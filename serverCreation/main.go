package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
		cancelCtx()
	}()
	go func() {
		err := serverTwo.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}

const keyServerAddr = "serverAddr"

func getRoot(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	hasFirst := request.URL.Query().Has("first")
	first := request.URL.Query().Get("first")
	hasSecond := request.URL.Query().Has("second")
	second := request.URL.Query().Get("second")

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Printf("Could not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s body:%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second, body,
	)
	io.WriteString(response, "This is my website!\n")
}

func getHello(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	fmt.Printf("%s: got /hello request \n", ctx.Value(keyServerAddr))

	myName := request.PostFormValue("myName")
	if myName == "" {
		response.Header().Set("x-missing-field", "myName")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(response, fmt.Sprintf("Hello, %s!\n", myName))
}
