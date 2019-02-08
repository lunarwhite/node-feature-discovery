/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nfdworker_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	w "sigs.k8s.io/node-feature-discovery/pkg/nfd-worker"
)

func TestNewNfdWorker(t *testing.T) {
	Convey("When initializing new NfdWorker instance", t, func() {
		Convey("When one of --cert-file, --key-file or --ca-file is missing", func() {
			_, err := w.NewNfdWorker(w.Args{CertFile: "crt", KeyFile: "key"})
			_, err2 := w.NewNfdWorker(w.Args{KeyFile: "key", CaFile: "ca"})
			_, err3 := w.NewNfdWorker(w.Args{CertFile: "crt", CaFile: "ca"})
			Convey("An error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err2, ShouldNotBeNil)
				So(err3, ShouldNotBeNil)
			})
		})
	})
}
