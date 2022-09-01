// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"paas-proxy-go/pulsar"
	"time"
)

func main() {
	r := mux.NewRouter()
	pulsar.Handlers(r.PathPrefix("/v1/pulsar").Subrouter())
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:20001",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
