//
// This file is automatically generated by scripts/generate-plugins.go -- Do not edit!
//

package command

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/plugin"

	filebuilder "github.com/hashicorp/packer/builder/file"
	nullbuilder "github.com/hashicorp/packer/builder/null"
	hcppackerartifactdatasource "github.com/hashicorp/packer/datasource/hcp-packer-artifact"
	hcppackerimagedatasource "github.com/hashicorp/packer/datasource/hcp-packer-image"
	hcppackeriterationdatasource "github.com/hashicorp/packer/datasource/hcp-packer-iteration"
	hcppackerversiondatasource "github.com/hashicorp/packer/datasource/hcp-packer-version"
	httpdatasource "github.com/hashicorp/packer/datasource/http"
	nulldatasource "github.com/hashicorp/packer/datasource/null"
	artificepostprocessor "github.com/hashicorp/packer/post-processor/artifice"
	checksumpostprocessor "github.com/hashicorp/packer/post-processor/checksum"
	compresspostprocessor "github.com/hashicorp/packer/post-processor/compress"
	manifestpostprocessor "github.com/hashicorp/packer/post-processor/manifest"
	shelllocalpostprocessor "github.com/hashicorp/packer/post-processor/shell-local"
	breakpointprovisioner "github.com/hashicorp/packer/provisioner/breakpoint"
	fileprovisioner "github.com/hashicorp/packer/provisioner/file"
	hcp_sbomprovisioner "github.com/hashicorp/packer/provisioner/hcp_sbom"
	powershellprovisioner "github.com/hashicorp/packer/provisioner/powershell"
	shellprovisioner "github.com/hashicorp/packer/provisioner/shell"
	shelllocalprovisioner "github.com/hashicorp/packer/provisioner/shell-local"
	sleepprovisioner "github.com/hashicorp/packer/provisioner/sleep"
	windowsrestartprovisioner "github.com/hashicorp/packer/provisioner/windows-restart"
	windowsshellprovisioner "github.com/hashicorp/packer/provisioner/windows-shell"
)

type ExecuteCommand struct {
	Meta
}

var Builders = map[string]packersdk.Builder{
	"file": new(filebuilder.Builder),
	"null": new(nullbuilder.Builder),
}

var Provisioners = map[string]packersdk.Provisioner{
	"breakpoint":      new(breakpointprovisioner.Provisioner),
	"file":            new(fileprovisioner.Provisioner),
	"hcp_sbom":        new(hcp_sbomprovisioner.Provisioner),
	"powershell":      new(powershellprovisioner.Provisioner),
	"shell":           new(shellprovisioner.Provisioner),
	"shell-local":     new(shelllocalprovisioner.Provisioner),
	"sleep":           new(sleepprovisioner.Provisioner),
	"windows-restart": new(windowsrestartprovisioner.Provisioner),
	"windows-shell":   new(windowsshellprovisioner.Provisioner),
}

var PostProcessors = map[string]packersdk.PostProcessor{
	"artifice":    new(artificepostprocessor.PostProcessor),
	"checksum":    new(checksumpostprocessor.PostProcessor),
	"compress":    new(compresspostprocessor.PostProcessor),
	"manifest":    new(manifestpostprocessor.PostProcessor),
	"shell-local": new(shelllocalpostprocessor.PostProcessor),
}

var Datasources = map[string]packersdk.Datasource{
	"hcp-packer-artifact":  new(hcppackerartifactdatasource.Datasource),
	"hcp-packer-image":     new(hcppackerimagedatasource.Datasource),
	"hcp-packer-iteration": new(hcppackeriterationdatasource.Datasource),
	"hcp-packer-version":   new(hcppackerversiondatasource.Datasource),
	"http":                 new(httpdatasource.Datasource),
	"null":                 new(nulldatasource.Datasource),
}

var pluginRegexp = regexp.MustCompile("packer-(builder|post-processor|provisioner|datasource)-(.+)")

func (c *ExecuteCommand) Run(args []string) int {
	// This is an internal call (users should not call this directly) so we're
	// not going to do much input validation. If there's a problem we'll often
	// just crash. Error handling should be added to facilitate debugging.
	log.Printf("args: %#v", args)
	if len(args) != 1 {
		c.Ui.Error(c.Help())
		return 1
	}

	// Plugin will match something like "packer-builder-amazon-ebs"
	parts := pluginRegexp.FindStringSubmatch(args[0])
	if len(parts) != 3 {
		c.Ui.Error(c.Help())
		return 1
	}
	pluginType := parts[1] // capture group 1 (builder|post-processor|provisioner)
	pluginName := parts[2] // capture group 2 (.+)

	server, err := plugin.Server()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error starting plugin server: %s", err))
		return 1
	}

	switch pluginType {
	case "builder":
		builder, found := Builders[pluginName]
		if !found {
			c.Ui.Error(fmt.Sprintf("Could not load builder: %s", pluginName))
			return 1
		}
		server.RegisterBuilder(builder)
	case "provisioner":
		provisioner, found := Provisioners[pluginName]
		if !found {
			c.Ui.Error(fmt.Sprintf("Could not load provisioner: %s", pluginName))
			return 1
		}
		server.RegisterProvisioner(provisioner)
	case "post-processor":
		postProcessor, found := PostProcessors[pluginName]
		if !found {
			c.Ui.Error(fmt.Sprintf("Could not load post-processor: %s", pluginName))
			return 1
		}
		server.RegisterPostProcessor(postProcessor)
	case "datasource":
		datasource, found := Datasources[pluginName]
		if !found {
			c.Ui.Error(fmt.Sprintf("Could not load datasource: %s", pluginName))
			return 1
		}
		server.RegisterDatasource(datasource)
	}

	server.Serve()

	return 0
}

func (*ExecuteCommand) Help() string {
	helpText := `
Usage: packer execute PLUGIN

  Runs an internally-compiled version of a plugin from the packer binary.

  NOTE: this is an internal command and you should not call it yourself.
`

	return strings.TrimSpace(helpText)
}

func (c *ExecuteCommand) Synopsis() string {
	return "internal plugin command"
}
