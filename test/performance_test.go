// main_test.go
package main_test

import (
	"fmt"
	"net/http"
	"sync"

	// Assuming CustomFormatter is in utils package

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Performance Test", Ordered, func() {
	Context("Single request", func() {
		testRequest := func() {
			It("Should handle 1 request per second", func() {
				resp, err := http.Get(fullURL)
				Expect(err).Should(BeNil())
				Expect(resp.StatusCode).Should(Equal(200))
			})
		}
		testRequest()
	})

	Context("Multiple requests", func() {
		testRequests := func(numRequests int) {
			It(fmt.Sprintf("Should handle %d requests per second", numRequests), func() {
				requests := make(chan *http.Response, numRequests)
				errors := make(chan error, numRequests)

				// Use a WaitGroup to wait for all goroutines to finish
				var wg sync.WaitGroup
				wg.Add(numRequests)

				// Launch goroutines to make HTTP requests
				for i := 0; i < numRequests; i++ {
					go func() {
						defer wg.Done()

						resp, err := http.Get(fullURL)
						if err != nil {
							errors <- err
							return
						}
						requests <- resp
					}()
				}

				// Wait for all goroutines to finish
				go func() {
					wg.Wait()
					close(requests)
					close(errors)
				}()

				// Collect results and make assertions
				successes := 0
				for {
					select {
					case resp, ok := <-requests:
						if !ok {
							requests = nil
						} else {
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
							successes++
						}
					case err, ok := <-errors:
						if !ok {
							errors = nil
						} else {
							logrus.Infoln("Number of successful requests:", successes)
							Fail(fmt.Sprintf("Unexpected error: %v", err))
						}
					}

					if requests == nil && errors == nil {
						break
					}
				}

				// Ensure the number of successful requests meets your expectations
				Expect(successes).To(Equal(numRequests))
			})
		}

		// Run the test with different numbers of requests
		for _, num := range []int{100, 125, 130} {
			testRequests(num)
		}
	})

})
