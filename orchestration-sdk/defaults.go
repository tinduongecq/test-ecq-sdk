package orchestrationsdk

import "time"

// Default client configuration values
const (
	DefaultBaseURL   = "http://localhost:8080"
	DefaultTimeout   = 30 * time.Second
	DefaultRetries   = 3
	DefaultRetryWait = 1 * time.Second
	DefaultUserAgent = "orchestration-sdk-go/1.0"
)

// Default resource configurations
var (
	// DefaultVMResources provides default resources for a VM
	DefaultVMResources = Resource{
		CPU:    2,
		Memory: "4Gi",
		Disk:   "50Gi",
	}

	// DefaultContainerResources provides default resources for a container
	DefaultContainerResources = Resource{
		CPU:    1,
		Memory: "512Mi",
		Disk:   "10Gi",
	}

	// DefaultEmulatorResources provides default resources for an emulator
	DefaultEmulatorResources = Resource{
		CPU:    2,
		Memory: "4Gi",
		Disk:   "20Gi",
	}

	// DevVMResources provides small resources for development VMs
	DevVMResources = Resource{
		CPU:    1,
		Memory: "2Gi",
		Disk:   "20Gi",
	}

	// ProdVMResources provides larger resources for production VMs
	ProdVMResources = Resource{
		CPU:    4,
		Memory: "8Gi",
		Disk:   "100Gi",
	}

	// HighPerformanceVMResources provides high-performance resources
	HighPerformanceVMResources = Resource{
		CPU:    8,
		Memory: "16Gi",
		Disk:   "200Gi",
	}
)

// Default API paths
const (
	PathOrchestration       = "/api/v1/orchestration"
	PathOrchestrationStatus = "/api/v1/orchestration/%s"
	PathOrchestrationCancel = "/api/v1/orchestration/%s/cancel"
	PathExecuteCommand      = "/api/v1/orchestration/%s/execute-command"
	PathStartVM             = "/api/v1/orchestration/%s/start"
	PathStopVM              = "/api/v1/orchestration/%s/stop"
	PathRemoveVM            = "/api/v1/orchestration/%s/remove"
	PathTemplates           = "/api/v1/templates"
	PathTemplateByID        = "/api/v1/templates/%s"
)

// Default orchestration settings
const (
	DefaultConfiguration = ConfigProxmox
	DefaultPriority      = PriorityMedium
	DefaultSchedule      = ScheduleImmediate
)

// ResourcePreset represents a predefined resource configuration
type ResourcePreset string

const (
	PresetDev             ResourcePreset = "dev"
	PresetProd            ResourcePreset = "prod"
	PresetHighPerformance ResourcePreset = "high_performance"
	PresetContainer       ResourcePreset = "container"
	PresetEmulator        ResourcePreset = "emulator"
)

// GetResourcePreset returns the Resource configuration for a preset
func GetResourcePreset(preset ResourcePreset) Resource {
	switch preset {
	case PresetDev:
		return DevVMResources
	case PresetProd:
		return ProdVMResources
	case PresetHighPerformance:
		return HighPerformanceVMResources
	case PresetContainer:
		return DefaultContainerResources
	case PresetEmulator:
		return DefaultEmulatorResources
	default:
		return DefaultVMResources
	}
}

// NetworkPreset represents a predefined network configuration
type NetworkPreset string

const (
	NetworkPresetBridge NetworkPreset = "bridge"
	NetworkPresetNAT    NetworkPreset = "nat"
	NetworkPresetHost   NetworkPreset = "host"
)

// GetNetworkPreset returns a NetworkConfig for a preset
func GetNetworkPreset(preset NetworkPreset) *NetworkConfig {
	switch preset {
	case NetworkPresetBridge:
		return &NetworkConfig{
			Interfaces: []NetworkInterface{
				{
					Type:   "bridge",
					Bridge: "vmbr0",
					Model:  "virtio",
				},
			},
		}
	case NetworkPresetNAT:
		return &NetworkConfig{
			Interfaces: []NetworkInterface{
				{
					Type:  "nat",
					Model: "virtio",
				},
			},
		}
	case NetworkPresetHost:
		return &NetworkConfig{
			Interfaces: []NetworkInterface{
				{
					Type: "host",
				},
			},
		}
	default:
		return nil
	}
}
