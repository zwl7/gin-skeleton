package utils

import "net/http"

type JaegerLogerHandler func(k string, v interface{})

type JaegerClientHeaderHandler func(url string, header http.Header)
