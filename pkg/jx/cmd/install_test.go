package cmd_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"fmt"

	configio "github.com/jenkins-x/jx/pkg/io"
	"github.com/jenkins-x/jx/pkg/jx/cmd"
	"github.com/jenkins-x/jx/pkg/jx/cmd/opts"
	"github.com/jenkins-x/jx/pkg/util"

	//. "github.com/petergtz/pegomock"
	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	t.Parallel()
	testDir := path.Join("test_data", "install_jenkins_x_versions")
	_, err := os.Stat(testDir)
	assert.NoError(t, err)

	configStore := configio.NewFileStore()
	version, err := cmd.LoadVersionFromCloudEnvironmentsDir(testDir, configStore)
	assert.NoError(t, err)

	assert.Equal(t, "0.0.3321", version, "For stable version in dir %s", testDir)
}

func TestGenerateProwSecret(t *testing.T) {
	fmt.Println(util.RandStringBytesMaskImprSrc(41))
}

func TestGetSafeUsername(t *testing.T) {
	t.Parallel()
	username := `Your active configuration is: [cloudshell-16392]
tutorial@bamboo-depth-206411.iam.gserviceaccount.com`
	assert.Equal(t, opts.GetSafeUsername(username), "tutorial@bamboo-depth-206411.iam.gserviceaccount.com")

	username = `tutorial@bamboo-depth-206411.iam.gserviceaccount.com`
	assert.Equal(t, opts.GetSafeUsername(username), "tutorial@bamboo-depth-206411.iam.gserviceaccount.com")
}

func TestCheckFlags(t *testing.T) {

	var tests = []struct {
		name           string
		in             *cmd.InstallFlags
		nextGeneration bool
		tekton         bool
		prow           bool
		staticJenkins  bool
		knativeBuild   bool
		err            error
	}{
		{
			name:           "default",
			in:             &cmd.InstallFlags{},
			nextGeneration: false,
			tekton:         false,
			prow:           false,
			staticJenkins:  true,
			knativeBuild:   false,
			err:            nil,
		},
		{
			name: "next_generation",
			in: &cmd.InstallFlags{
				NextGeneration: true,
			},
			nextGeneration: true,
			tekton:         true,
			prow:           true,
			staticJenkins:  false,
			knativeBuild:   false,
			err:            nil,
		},
		{
			name: "prow",
			in: &cmd.InstallFlags{
				Prow: true,
			},
			nextGeneration: false,
			tekton:         true,
			prow:           true,
			staticJenkins:  false,
			knativeBuild:   false,
			err:            nil,
		},
		{
			name: "prow_and_knative",
			in: &cmd.InstallFlags{
				Prow:         true,
				KnativeBuild: true,
			},
			nextGeneration: false,
			tekton:         false,
			prow:           true,
			staticJenkins:  false,
			knativeBuild:   true,
			err:            nil,
		},
		{
			name: "next_generation_and_static_jenkins",
			in: &cmd.InstallFlags{
				NextGeneration: true,
				StaticJenkins:  true,
			},
			err: errors.New("Incompatible options '--ng' and '--static-jenkins'. Please pick only one of them. We recommend --ng as --static-jenkins is deprecated"),
		},
		{
			name: "tekton_and_static_jenkins",
			in: &cmd.InstallFlags{
				Tekton:        true,
				StaticJenkins: true,
			},
			err: errors.New("Incompatible options '--tekton' and '--static-jenkins'. Please pick only one of them. We recommend --tekton as --static-jenkins is deprecated"),
		},
		{
			name: "tekton_and_knative",
			in: &cmd.InstallFlags{
				Tekton:       true,
				KnativeBuild: true,
			},
			err: errors.New("Incompatible options '--knative-build' and '--tekton'. Please pick only one of them. We recommend --tekton as --knative-build is deprecated"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := cmd.InstallOptions{
				CommonOptions: &opts.CommonOptions{
					BatchMode: true,
				},
				Flags:       *tt.in,
				InitOptions: cmd.InitOptions{},
			}

			err := opts.CheckFlags()
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.nextGeneration, opts.Flags.NextGeneration, "NextGeneration flag is not as expected")
				assert.Equal(t, tt.tekton, opts.Flags.Tekton, "Tekton flag is not as expected")
				assert.Equal(t, tt.prow, opts.Flags.Prow, "Prow flag is not as expected")
				assert.Equal(t, tt.staticJenkins, opts.Flags.StaticJenkins, "StaticJenkins flag is not as expected")
				assert.Equal(t, tt.knativeBuild, opts.Flags.KnativeBuild, "KnativeBuild flag is not as expected")
			}
		})
	}
}

func TestInstallRun(t *testing.T) {
	// Create mocks...
	//factory := cmd_mocks.NewMockFactory()
	//kubernetesInterface := kube_mocks.NewSimpleClientset()
	//// Override CreateKubeClient to return mock Kubernetes interface
	//When(factory.CreateKubeClient()).ThenReturn(kubernetesInterface, "jx-testing", nil)

	//options := cmd.CreateInstallOptions(factory, os.Stdin, os.Stdout, os.Stderr)

	//err := options.Run()

	//assert.NoError(t, err, "Should not error")
}
