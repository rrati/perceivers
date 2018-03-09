/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var httpResults *prometheus.CounterVec
var durationsHistogram *prometheus.HistogramVec
var errorsCounter *prometheus.CounterVec

// helpers

func RecordError(errorStage string, errorName string) {
	log.Infof("metrics record error %s, %s", errorStage, errorName)
	errorsCounter.With(prometheus.Labels{"stage": errorStage, "errorName": errorName}).Inc()
}

func RecordDuration(operation string, duration time.Duration) {
	log.Infof("record duration %s, %s", operation, duration)
	durationsHistogram.With(prometheus.Labels{"operation": operation}).Observe(duration.Seconds())
}

// recorders

func RecordHttpStats(path string, success bool) {
	log.Infof("record http stats -- %s, %t", path, success)
	httpResults.With(prometheus.Labels{"path": path, "result": fmt.Sprintf("%t", success)}).Inc()
}

// init

func init() {
	httpResults = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "perceptor",
		Subsystem: "image_perceiver",
		Name:      "http_response_status_codes",
		Help:      "success/failure responses from HTTP requests issued by image perceiver",
	},
		[]string{"path", "result"})

	durationsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "perceptor",
			Subsystem: "image_perceiver",
			Name:      "timings",
			Help:      "time durations of image perceiver operations",
			Buckets:   prometheus.ExponentialBuckets(0.25, 2, 20),
		},
		[]string{"operation"})

	errorsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "perceptor",
		Subsystem: "image_perceiver",
		Name:      "errors",
		Help:      "errors from image perceiver operations",
	}, []string{"stage", "errorName"})

	prometheus.MustRegister(errorsCounter)
	prometheus.MustRegister(durationsHistogram)
	prometheus.MustRegister(httpResults)
}
