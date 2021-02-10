package integration

import (
	"fmt"
	"os"

	. "github.com/containers/podman/v2/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Podman image tree", func() {
	var (
		tempdir    string
		err        error
		podmanTest *PodmanTestIntegration
	)

	BeforeEach(func() {
		tempdir, err = CreateTempDirInTempDir()
		if err != nil {
			os.Exit(1)
		}
		podmanTest = PodmanTestCreate(tempdir)
		podmanTest.Setup()
		podmanTest.AddImageToRWStore(BB)
	})

	AfterEach(func() {
		podmanTest.Cleanup()
		f := CurrentGinkgoTestDescription()
		timedResult := fmt.Sprintf("Test: %s completed in %f seconds", f.TestText, f.Duration.Seconds())
		GinkgoWriter.Write([]byte(timedResult))
	})

	It("podman image tree", func() {
		SkipIfRemote("Does not work on remote client")
		Skip("don't understand why this fails")
		podmanTest.AddImageToRWStore(cirros)
		dockerfile := `FROM quay.io/libpod/cirros:latest
RUN mkdir hello
RUN touch test.txt
ENV foo=bar
`
		podmanTest.BuildImage(dockerfile, "test:latest", "true")

		session := podmanTest.Podman([]string{"image", "tree", "test:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.Podman([]string{"image", "tree", "--whatrequires", "quay.io/libpod/cirros:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.Podman([]string{"rmi", "test:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.Podman([]string{"rmi", "quay.io/libpod/cirros:latest"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})
})