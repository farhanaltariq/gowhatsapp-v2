package main_test

import (
	"github.com/farhanaltariq/fiberplate/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encryption", Ordered, func() {

	Context("Encryption And Decryption", func() {
		It("Encrypt and decrypt lowercase", func() {
			result, salt := utils.Encrypt("test")
			decrypt, _ := utils.Decrypt(result, salt)
			Expect(decrypt).To(Equal("test"))
		})

		It("Encrypt and decrypt uppercase", func() {
			result, salt := utils.Encrypt("TEST")
			decrypt, _ := utils.Decrypt(result, salt)
			Expect(decrypt).To(Equal("TEST"))
		})

		It("Encrypt and decrypt number", func() {
			result, salt := utils.Encrypt("12345")
			decrypt, _ := utils.Decrypt(result, salt)
			Expect(decrypt).To(Equal("12345"))
		})

		It("Encrypt and decrypt mixed character", func() {
			result, salt := utils.Encrypt("sEcr3t")
			decrypt, _ := utils.Decrypt(result, salt)
			Expect(decrypt).To(Equal("sEcr3t"))
		})

		It("Encrypt and decrypt special character", func() {
			result, salt := utils.Encrypt("sEcr3t!@129484weruoqrioeufp[q[woi[irwq[fdk]]]]")
			decrypt, _ := utils.Decrypt(result, salt)
			Expect(decrypt).To(Equal("sEcr3t!@129484weruoqrioeufp[q[woi[irwq[fdk]]]]"))
		})
	})
})
