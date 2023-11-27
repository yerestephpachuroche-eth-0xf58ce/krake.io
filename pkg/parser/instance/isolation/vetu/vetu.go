package vetu

import (
	"github.com/cirruslabs/cirrus-ci-agent/api"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/instance/resources"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/nameable"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/node"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/parseable"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/parserkit"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/schema"
	jsschema "github.com/lestrrat-go/jsschema"
	"strconv"
)

type Vetu struct {
	proto *api.Isolation_Vetu_

	parseable.DefaultParser
}

func New(mergedEnv map[string]string, parserKit *parserkit.ParserKit) *Vetu {
	vetu := &Vetu{
		proto: &api.Isolation_Vetu_{
			Vetu: &api.Isolation_Vetu{},
		},
	}

	vmSchema := schema.String("Source VM image (or name) to clone the new VM from.")
	vetu.OptionalField(nameable.NewSimpleNameable("image"), vmSchema, func(node *node.Node) error {
		image, err := node.GetExpandedStringValue(mergedEnv)
		if err != nil {
			return err
		}

		vetu.proto.Vetu.Image = image

		return nil
	})

	userSchema := schema.String("SSH username.")
	vetu.OptionalField(nameable.NewSimpleNameable("user"), userSchema, func(node *node.Node) error {
		user, err := node.GetExpandedStringValue(mergedEnv)
		if err != nil {
			return err
		}

		vetu.proto.Vetu.User = user

		return nil
	})

	passwordSchema := schema.String("SSH password.")
	vetu.OptionalField(nameable.NewSimpleNameable("password"), passwordSchema, func(node *node.Node) error {
		password, err := node.GetExpandedStringValue(mergedEnv)
		if err != nil {
			return err
		}

		vetu.proto.Vetu.Password = password

		return nil
	})

	cpuSchema := schema.Number("Number of VM CPUs.")
	vetu.OptionalField(nameable.NewSimpleNameable("cpu"), cpuSchema, func(node *node.Node) error {
		cpu, err := node.GetExpandedStringValue(mergedEnv)
		if err != nil {
			return err
		}
		cpuParsed, err := strconv.ParseUint(cpu, 10, 32)
		if err != nil {
			return node.ParserError("%s", err.Error())
		}
		vetu.proto.Vetu.Cpu = uint32(cpuParsed)
		return nil
	})

	memorySchema := schema.Memory()
	memorySchema.Description = "VM memory size in megabytes."
	vetu.OptionalField(nameable.NewSimpleNameable("memory"), memorySchema, func(node *node.Node) error {
		memory, err := node.GetExpandedStringValue(mergedEnv)
		if err != nil {
			return err
		}
		memoryParsed, err := resources.ParseMegaBytes(memory)
		if err != nil {
			return node.ParserError("%s", err.Error())
		}
		vetu.proto.Vetu.Memory = uint32(memoryParsed)
		return nil
	})

	return vetu
}

func (vetu *Vetu) Parse(node *node.Node, parserKit *parserkit.ParserKit) error {
	return vetu.DefaultParser.Parse(node, parserKit)
}

func (vetu *Vetu) Proto() *api.Isolation_Vetu_ {
	return vetu.proto
}

func (vetu *Vetu) Schema() *jsschema.Schema {
	modifiedSchema := vetu.DefaultParser.Schema()

	modifiedSchema.Type = jsschema.PrimitiveTypes{jsschema.ObjectType}
	modifiedSchema.Description = "Vetu VM isolation."

	return modifiedSchema
}
