package v1alpha1

import (
	"github.com/redpanda-data/console/backend/pkg/connect"
	"github.com/redpanda-data/console/backend/pkg/filesystem"
	"github.com/redpanda-data/console/backend/pkg/git"
	"github.com/redpanda-data/console/backend/pkg/msgpack"
	"github.com/redpanda-data/console/backend/pkg/proto"
)

// Alias types are defined here because of controller-gen "make generate" does not generate DeepCopyInto function on these types.
// REF https://github.com/kubernetes-sigs/controller-tools/issues/583

type (
	// ConnectConfig is an alias to connect.Config
	ConnectConfig connect.Config

	// MsgpackConfig is an alias to msgpack.Config
	MsgpackConfig msgpack.Config

	// GitConfig is an alias to git.Config
	GitConfig git.Config

	// FilesystemConfig is an alias to filesystem.Config
	FilesystemConfig filesystem.Config
)

// ConsoleConfig is a redefinition of console.Config but struct fields defined locally
type ConsoleConfig struct {
	TopicDocumentation ConfigTopicDocumentation `json:"topicDocumentation"`
}

// ConfigTopicDocumentation is a redefinition of console.ConfigTopicDocumentation but struct fields defined locally
type ConfigTopicDocumentation struct {
	Enabled bool      `json:"enabled"`
	Git     GitConfig `json:"git"`
}

// ProtoConfig is a redefinition of console.ProtoConfig but struct fields defined locally
type ProtoConfig struct {
	Enabled bool `json:"enabled"`

	// The required proto definitions can be provided via SchemaRegistry, Git or Filesystem
	SchemaRegistry proto.SchemaRegistryConfig `json:"schemaRegistry"`
	Git            GitConfig                  `json:"git"`
	FileSystem     FilesystemConfig           `json:"fileSystem"`

	// Mappings define what proto types shall be used for each Kafka topic. If SchemaRegistry is used, no mappings are required.
	Mappings []proto.ConfigTopicMapping `json:"mappings"`
}
