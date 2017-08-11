package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/post", toPost)
	http.HandleFunc("/put", toPut)
	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func toPost(w http.ResponseWriter, r *http.Request) {
	info, err := getRequestInfo(r)
	if err != nil {
		fmt.Println(err)
	}
	info.method = "POST"
	doRequest(info, w)
}

func toPut(w http.ResponseWriter, r *http.Request) {
	info, err := getRequestInfo(r)
	if err != nil {
		fmt.Println(err)
	}
	info.method = "PUT"
	doRequest(info, w)
}

func doRequest(info requestInfo, w http.ResponseWriter) error {
	if info.debug {
		info.print(w)
	}

	httpClient := getHTTPClient(info)

	urlString := getURLString(info)
	if info.debug {
		w.Write([]byte("urlString: " + urlString + "\n"))
	}

	var bodyReader io.Reader
	if info.body != "" {
		bodyReader = bytes.NewBuffer([]byte(info.body))
	}

	request, err := http.NewRequest(info.method, urlString, bodyReader)
	if err != nil {
		w.Write([]byte(err.Error() + "\n"))
		return nil
	}

	for k, v := range info.headers {
		request.Header.Add(k, v)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		w.Write([]byte(err.Error() + "\n"))
		return nil
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(result)
	}

	return nil
}

func getRequestInfo(r *http.Request) (requestInfo, error) {
	q := r.URL.Query()
	info := requestInfo{
		debug:           returnFlagAndDelAll(q, "debug"),
		ignoreSslErrors: returnFlagAndDelAll(q, "ignoreSslErrors"),
		protocol:        returnFirstAndDelAll(q, "protocol"),
		host:            returnFirstAndDelAll(q, "host"),
		port:            returnFirstAndDelAll(q, "port"),
		path:            returnFirstAndDelAll(q, "path"),
		body:            returnFirstAndDelAll(q, "body"),
		headers:         returnAndDelAllHeaders(q),
	}
	info.query = q
	return info, nil
}

const headerParamPrefix = "_"

func returnFlagAndDelAll(q url.Values, key string) bool {
	_, ok := q[key]
	if ok {
		q.Del(key)
	}
	return ok
}

func returnAndDelAllHeaders(q url.Values) map[string]string {
	h := make(map[string]string)
	for k, v := range q {
		if k == headerParamPrefix+"Authorization_Basic" {
			h["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(v[0]))
			q.Del(headerParamPrefix + "Authorization_Basic")
		} else if strings.HasPrefix(k, headerParamPrefix) {
			h[k[len(headerParamPrefix):]] = strings.Join(v, ",")
			q.Del(k)
		}
	}
	return h
}

func returnFirstAndDelAll(q url.Values, key string) string {
	val := ""
	items := q[key]
	if len(items) > 0 {
		val = items[0]
	}
	q.Del(key)
	return val
}

type requestInfo struct {
	debug           bool
	ignoreSslErrors bool
	method          string
	protocol        string
	host            string
	port            string
	path            string
	query           url.Values
	body            string
	headers         map[string]string
}

func (i requestInfo) print(w http.ResponseWriter) {
	w.Write([]byte("ignoreSslErrors: " + strconv.FormatBool(i.ignoreSslErrors) + "\n"))
	w.Write([]byte("method: " + i.method + "\n"))
	w.Write([]byte("protocol: " + i.protocol + "\n"))
	w.Write([]byte("host: " + i.host + "\n"))
	w.Write([]byte("port: " + i.port + "\n"))
	w.Write([]byte("path: " + i.path + "\n"))
	w.Write([]byte("query: " + i.query.Encode() + "\n"))
	w.Write([]byte("body: " + i.body + "\n"))
	for k, v := range i.headers {
		w.Write([]byte("headers: " + k + ": " + v + "\n"))
	}
}

func getHTTPClient(info requestInfo) *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	if info.protocol == "https" && info.ignoreSslErrors {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return client
}

func getURLString(info requestInfo) string {
	urlString := info.protocol + "://" + info.host
	if info.port != "" {
		urlString = urlString + ":" + info.port
	}
	urlString = urlString + info.path
	queryString := info.query.Encode()
	if queryString != "" {
		urlString += "?" + queryString
	}
	return urlString
}
