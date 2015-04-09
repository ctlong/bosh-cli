package templatescompiler_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"

	bmproperty "github.com/cloudfoundry/bosh-init/common/property"
	bmreljob "github.com/cloudfoundry/bosh-init/release/job"
	bmerbrenderer "github.com/cloudfoundry/bosh-init/templatescompiler/erbrenderer"

	fakebmrender "github.com/cloudfoundry/bosh-init/templatescompiler/erbrenderer/fakes"

	. "github.com/cloudfoundry/bosh-init/templatescompiler"
)

var _ = Describe("JobRenderer", func() {
	var (
		jobRenderer      JobRenderer
		fakeERBRenderer  *fakebmrender.FakeERBRenderer
		job              bmreljob.Job
		context          bmerbrenderer.TemplateEvaluationContext
		fs               *fakesys.FakeFileSystem
		jobProperties    bmproperty.Map
		globalProperties bmproperty.Map
		srcPath          string
		dstPath          string
	)

	BeforeEach(func() {
		srcPath = "fake-src-path"
		dstPath = "fake-dst-path"
		jobProperties = bmproperty.Map{
			"fake-property-key": "fake-job-property-value",
		}

		globalProperties = bmproperty.Map{
			"fake-property-key": "fake-global-property-value",
		}

		job = bmreljob.Job{
			Templates: map[string]string{
				"director.yml.erb": "config/director.yml",
			},
			ExtractedPath: srcPath,
		}

		logger := boshlog.NewLogger(boshlog.LevelNone)

		context = NewJobEvaluationContext(job, jobProperties, globalProperties, "fake-deployment-name", logger)

		fakeERBRenderer = fakebmrender.NewFakeERBRender()

		fs = fakesys.NewFakeFileSystem()
		jobRenderer = NewJobRenderer(fakeERBRenderer, fs, logger)

		fakeERBRenderer.SetRenderBehavior(
			filepath.Join(srcPath, "templates/director.yml.erb"),
			filepath.Join(dstPath, "config/director.yml"),
			context,
			nil,
		)

		fakeERBRenderer.SetRenderBehavior(
			filepath.Join(srcPath, "monit"),
			filepath.Join(dstPath, "monit"),
			context,
			nil,
		)

		fs.TempDirDir = dstPath
	})

	AfterEach(func() {
		err := fs.RemoveAll(dstPath)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Render", func() {
		It("renders job templates", func() {
			renderedjob, err := jobRenderer.Render(job, jobProperties, globalProperties, "fake-deployment-name")
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeERBRenderer.RenderInputs).To(Equal([]fakebmrender.RenderInput{
				{
					SrcPath: filepath.Join(srcPath, "templates/director.yml.erb"),
					DstPath: filepath.Join(renderedjob.Path(), "config/director.yml"),
					Context: context,
				},
				{
					SrcPath: filepath.Join(srcPath, "monit"),
					DstPath: filepath.Join(renderedjob.Path(), "monit"),
					Context: context,
				},
			}))
		})

		Context("when rendering fails", func() {
			BeforeEach(func() {
				fakeERBRenderer.SetRenderBehavior(
					filepath.Join(srcPath, "templates/director.yml.erb"),
					filepath.Join(dstPath, "config/director.yml"),
					context,
					bosherr.Error("fake-template-render-error"),
				)
			})

			It("returns an error", func() {
				_, err := jobRenderer.Render(job, jobProperties, globalProperties, "fake-deployment-name")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-template-render-error"))
			})
		})
	})
})
