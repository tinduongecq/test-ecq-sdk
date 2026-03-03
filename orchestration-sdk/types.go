// Package orchestrationsdk provides a Go SDK client for the Orchestration Engine API.
// It simplifies creating and managing VMs, containers, and emulators.
package orchestrationsdk

import "time"

// ProvisionType represents the type of provisioning job
type ProvisionType string

const (
	ProvisionTypeVM        ProvisionType = "vm_provision"
	ProvisionTypeContainer ProvisionType = "container_provision"
	ProvisionTypeEmulator  ProvisionType = "emulator_provision"
)

// Configuration represents the infrastructure configuration
type Configuration string

const (
	ConfigProxmox   Configuration = "proxmox"
	ConfigOpenStack Configuration = "openstack"
	ConfigESXi      Configuration = "esxi"
	ConfigAWS       Configuration = "aws"
	ConfigAzure     Configuration = "azure"
	ConfigGCP       Configuration = "gcp"
)

// Priority represents the job priority
type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

// Schedule represents the job schedule type
type Schedule string

const (
	ScheduleImmediate Schedule = "immediate"
	ScheduleScheduled Schedule = "scheduled"
)

// TemplateType represents template types
type TemplateType string

const (
	TemplateTypeVM        TemplateType = "vm"
	TemplateTypeContainer TemplateType = "container"
	TemplateTypeEmulator  TemplateType = "emulator"
)

// RegistryType represents container registry types
type RegistryType string

const (
	RegistryTypeDockerHub RegistryType = "docker_hub"
	RegistryTypePrivate   RegistryType = "private"
)

// OSType represents operating system types
type OSType string

const (
	OSTypeLinux   OSType = "linux"
	OSTypeWindows OSType = "windows"
	OSTypeMacOS   OSType = "macos"
	OSTypeAndroid OSType = "android"
	OSTypeiOS     OSType = "ios"
)

// Resource represents the resource allocation for a VM/container
type Resource struct {
	CPU     int    `json:"cpu"`
	Memory  string `json:"memory"`
	Disk    string `json:"disk"`
	Network string `json:"network,omitempty"`
}

// NetworkInterface represents a network interface configuration
type NetworkInterface struct {
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
	Subnet    string `json:"subnet,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	VLAN      int    `json:"vlan,omitempty"`
	Available bool   `json:"available,omitempty"`
	Model     string `json:"model,omitempty"`
	Bridge    string `json:"bridge,omitempty"`
	Firewall  bool   `json:"firewall,omitempty"`
	MACAddr   string `json:"macaddr,omitempty"`
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	Interfaces []NetworkInterface `json:"interfaces,omitempty"`
	IPAddress  string             `json:"ip_address,omitempty"`
	Gateway    string             `json:"gateway,omitempty"`
}

// FileConfig represents a file to be injected into the VM
type FileConfig struct {
	URL             string `json:"url"`
	DestinationPath string `json:"destination_path"`
	Filename        string `json:"filename"`
}

// OrchestrationRequest represents the request to create an orchestration
type OrchestrationRequest struct {
	Name             string            `json:"name"`
	JobProvisionType ProvisionType     `json:"job_provision_type"`
	TemplateID       string            `json:"template_id"`
	Configuration    Configuration     `json:"configuration"`
	Resources        Resource          `json:"resources"`
	Network          *NetworkConfig    `json:"network,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	Priority         Priority          `json:"priority"`
	Schedule         Schedule          `json:"schedule"`
	ResourcePoolID   string            `json:"resource_pool_id"`
	VMEngineID       int               `json:"vm_engine_id,omitempty"`
	File             *FileConfig       `json:"file,omitempty"`
}

// TemplateInfo represents template information in the response
type TemplateInfo struct {
	ID           string `json:"id"`
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	ImageName    string `json:"image_name"`
	Type         string `json:"type"`
	OSType       string `json:"os_type"`
	Architecture string `json:"architecture"`
	DiskSize     uint64 `json:"disk_size"`
}

// ResourcePoolInfo represents resource pool information in the response
type ResourcePoolInfo struct {
	ID              string `json:"id"`
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	AvailableCPU    int    `json:"available_cpu"`
	AvailableMemory string `json:"available_memory"`
	AvailableDisk   string `json:"available_disk"`
}

// OrchestrationResponse represents the response from creating an orchestration
type OrchestrationResponse struct {
	OrchestrationID string            `json:"orchestration_id"`
	Status          string            `json:"status"`
	JobID           string            `json:"job_id"`
	Message         string            `json:"message"`
	ResourcePool    *ResourcePoolInfo `json:"resource_pool,omitempty"`
	Template        *TemplateInfo     `json:"template,omitempty"`
	CreatedAt       string            `json:"created_at"`
}

// OrchestrationStatusResponse represents orchestration status
type OrchestrationStatusResponse struct {
	OrchestrationID string            `json:"orchestration_id"`
	Status          string            `json:"status"`
	Progress        int               `json:"progress"`
	Message         string            `json:"message"`
	JobID           string            `json:"job_id"`
	ResourcePool    *ResourcePoolInfo `json:"resource_pool,omitempty"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// OrchestrationListResponse represents a list of orchestrations
type OrchestrationListResponse struct {
	Items      []*OrchestrationStatusResponse `json:"items"`
	Pagination PaginationInfo                 `json:"pagination"`
}

// ExecuteCommandRequest represents a command execution request
type ExecuteCommandRequest struct {
	Command []string `json:"command"`
}

// ExecuteCommandResponse represents a command execution response
type ExecuteCommandResponse struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

// APIResponse represents the generic API response wrapper
type APIResponse struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ListOptions represents options for listing resources
type ListOptions struct {
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// Template represents a full template response
type Template struct {
	ID             string            `json:"id"`
	UUID           string            `json:"uuid"`
	Name           string            `json:"name"`
	ImageName      string            `json:"image_name"`
	RegistryType   RegistryType      `json:"registry_type"`
	Description    string            `json:"description"`
	Type           TemplateType      `json:"type"`
	OSType         OSType            `json:"os_type"`
	Architecture   string            `json:"architecture"`
	DiskSize       uint64            `json:"disk_size"`
	CPU            string            `json:"cpu"`
	Memory         int               `json:"memory"`
	Metadata       map[string]string `json:"metadata"`
	ResourcePoolID string            `json:"resource_pool_id"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// TemplateListResponse represents a paginated list of templates
type TemplateListResponse struct {
	Items      []*Template    `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

// TemplateListOptions represents options for listing templates
type TemplateListOptions struct {
	ListOptions
	Type         TemplateType `json:"type,omitempty"`
	OSType       OSType       `json:"os_type,omitempty"`
	Architecture string       `json:"architecture,omitempty"`
	Name         string       `json:"name,omitempty"`
}
