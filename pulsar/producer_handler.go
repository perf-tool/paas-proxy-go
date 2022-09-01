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

package pulsar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func produceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenant := vars["tenant"]
	namespace := vars["namespace"]
	topic := vars["topic"]
	var req ProduceMsgReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cost, err := produce(tenant, namespace, topic, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := ProduceMsgResp{Cost: cost}
	payload, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	logrus.Error("write payload failed ", err)
}

func produce(tenant, namespace, topic string, req ProduceMsgReq) (int64, error) {
	startTime := time.Now().Unix()
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               fmt.Sprintf("pulsar://%s:%d", req.Host, req.Port),
		ConnectionTimeout: time.Duration(30) * time.Second,
		OperationTimeout:  time.Duration(30) * time.Second,
	})
	if err != nil {
		logrus.Error("Failed to create client ", err)
		return 0, err
	}
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: fmt.Sprintf("persistent://%s/%s/%s", tenant, namespace, topic),
	})
	if err != nil {
		logrus.Error("Failed to create producer ", err)
		return 0, err
	}
	defer producer.Close()
	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte(req.Msg),
	})
	if err != nil {
		logrus.Error("Failed to publish message ", err)
		return 0, err
	}
	return time.Now().Unix() - startTime, nil
}
