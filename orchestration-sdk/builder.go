package orchestrationsdk

// RequestBuilder is the base builder for orchestration requests
type RequestBuilder struct {
	request *OrchestrationRequest
	errors  []error
}

// VMRequestBuilder builds VM provisioning requests
type VMRequestBuilder struct {
	RequestBuilder
}

// ContainerRequestBuilder builds Container provisioning requests
type ContainerRequestBuilder struct {
	RequestBuilder
}

// EmulatorRequestBuilder builds Emulator provisioning requests
type EmulatorRequestBuilder struct {
	RequestBuilder
}

// ===== VM Request Builder =====

// NewVMRequest creates a new VM request builder with the given name
func NewVMRequest(name string) *VMRequestBuilder {
	return &VMRequestBuilder{
		RequestBuilder: RequestBuilder{
			request: &OrchestrationRequest{
				Name:             name,
				JobProvisionType: ProvisionTypeVM,
				Configuration:    DefaultConfiguration,
				Priority:         DefaultPriority,
				Schedule:         DefaultSchedule,
				Resources:        DefaultVMResources,
			},
		},
	}
}

// NewDevVMRequest creates a development VM with small resources and low priority
func NewDevVMRequest(name string) *VMRequestBuilder {
	return &VMRequestBuilder{
		RequestBuilder: RequestBuilder{
			request: &OrchestrationRequest{
				Name:             name,
				JobProvisionType: ProvisionTypeVM,
				Configuration:    DefaultConfiguration,
				Priority:         PriorityLow,
				Schedule:         DefaultSchedule,
				Resources:        DevVMResources,
				Metadata:         map[string]string{"environment": "development"},
			},
		},
	}
}

// NewProdVMRequest creates a production VM with larger resources and high priority
func NewProdVMRequest(name string) *VMRequestBuilder {
	return &VMRequestBuilder{
		RequestBuilder: RequestBuilder{
			request: &OrchestrationRequest{
				Name:             name,
				JobProvisionType: ProvisionTypeVM,
				Configuration:    DefaultConfiguration,
				Priority:         PriorityHigh,
				Schedule:         DefaultSchedule,
				Resources:        ProdVMResources,
				Metadata:         map[string]string{"environment": "production"},
			},
		},
	}
}

// WithTemplate sets the template ID
func (b *VMRequestBuilder) WithTemplate(templateID string) *VMRequestBuilder {
	b.request.TemplateID = templateID
	return b
}

// WithResourcePool sets the resource pool ID
func (b *VMRequestBuilder) WithResourcePool(resourcePoolID string) *VMRequestBuilder {
	b.request.ResourcePoolID = resourcePoolID
	return b
}

// WithResources sets the resource allocation (CPU cores, Memory, Disk)
func (b *VMRequestBuilder) WithResources(cpu int, memory, disk string) *VMRequestBuilder {
	b.request.Resources = Resource{
		CPU:    cpu,
		Memory: memory,
		Disk:   disk,
	}
	return b
}

// WithResourcePreset sets resources using a preset configuration
func (b *VMRequestBuilder) WithResourcePreset(preset ResourcePreset) *VMRequestBuilder {
	b.request.Resources = GetResourcePreset(preset)
	return b
}

// WithConfiguration sets the infrastructure configuration
func (b *VMRequestBuilder) WithConfiguration(config Configuration) *VMRequestBuilder {
	b.request.Configuration = config
	return b
}

// WithPriority sets the job priority
func (b *VMRequestBuilder) WithPriority(priority Priority) *VMRequestBuilder {
	b.request.Priority = priority
	return b
}

// WithSchedule sets the job schedule
func (b *VMRequestBuilder) WithSchedule(schedule Schedule) *VMRequestBuilder {
	b.request.Schedule = schedule
	return b
}

// WithNetwork sets the network configuration
func (b *VMRequestBuilder) WithNetwork(network NetworkConfig) *VMRequestBuilder {
	b.request.Network = &network
	return b
}

// WithNetworkPreset sets network using a preset configuration
func (b *VMRequestBuilder) WithNetworkPreset(preset NetworkPreset) *VMRequestBuilder {
	b.request.Network = GetNetworkPreset(preset)
	return b
}

// WithMetadata sets custom metadata
func (b *VMRequestBuilder) WithMetadata(metadata map[string]string) *VMRequestBuilder {
	b.request.Metadata = metadata
	return b
}

// AddMetadata adds a single metadata key-value pair
func (b *VMRequestBuilder) AddMetadata(key, value string) *VMRequestBuilder {
	if b.request.Metadata == nil {
		b.request.Metadata = make(map[string]string)
	}
	b.request.Metadata[key] = value
	return b
}

// WithVMEngineID sets the VM engine ID (for Proxmox)
func (b *VMRequestBuilder) WithVMEngineID(id int) *VMRequestBuilder {
	b.request.VMEngineID = id
	return b
}

// WithFile sets the file to be injected into the VM
func (b *VMRequestBuilder) WithFile(url, destinationPath, filename string) *VMRequestBuilder {
	b.request.File = &FileConfig{
		URL:             url,
		DestinationPath: destinationPath,
		Filename:        filename,
	}
	return b
}

// Build validates and returns the orchestration request
func (b *VMRequestBuilder) Build() (*OrchestrationRequest, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}
	return b.request, nil
}

// MustBuild returns the orchestration request, panics on validation error
func (b *VMRequestBuilder) MustBuild() *OrchestrationRequest {
	req, err := b.Build()
	if err != nil {
		panic(err)
	}
	return req
}

func (b *VMRequestBuilder) validate() error {
	if b.request.Name == "" {
		return ErrNameRequired
	}
	if b.request.TemplateID == "" {
		return ErrTemplateIDRequired
	}
	if b.request.ResourcePoolID == "" {
		return ErrResourcePoolIDRequired
	}
	if b.request.Resources.CPU <= 0 {
		return ErrInvalidCPU
	}
	if b.request.Resources.Memory == "" {
		return ErrInvalidMemory
	}
	return nil
}

// ===== Container Request Builder =====

// NewContainerRequest creates a new Container request builder
func NewContainerRequest(name string) *ContainerRequestBuilder {
	return &ContainerRequestBuilder{
		RequestBuilder: RequestBuilder{
			request: &OrchestrationRequest{
				Name:             name,
				JobProvisionType: ProvisionTypeContainer,
				Configuration:    DefaultConfiguration,
				Priority:         DefaultPriority,
				Schedule:         DefaultSchedule,
				Resources:        DefaultContainerResources,
			},
		},
	}
}

// WithTemplate sets the template ID
func (b *ContainerRequestBuilder) WithTemplate(templateID string) *ContainerRequestBuilder {
	b.request.TemplateID = templateID
	return b
}

// WithResourcePool sets the resource pool ID
func (b *ContainerRequestBuilder) WithResourcePool(resourcePoolID string) *ContainerRequestBuilder {
	b.request.ResourcePoolID = resourcePoolID
	return b
}

// WithResources sets the resource allocation
func (b *ContainerRequestBuilder) WithResources(cpu int, memory, disk string) *ContainerRequestBuilder {
	b.request.Resources = Resource{
		CPU:    cpu,
		Memory: memory,
		Disk:   disk,
	}
	return b
}

// WithConfiguration sets the infrastructure configuration
func (b *ContainerRequestBuilder) WithConfiguration(config Configuration) *ContainerRequestBuilder {
	b.request.Configuration = config
	return b
}

// WithPriority sets the job priority
func (b *ContainerRequestBuilder) WithPriority(priority Priority) *ContainerRequestBuilder {
	b.request.Priority = priority
	return b
}

// WithNetwork sets the network configuration
func (b *ContainerRequestBuilder) WithNetwork(network NetworkConfig) *ContainerRequestBuilder {
	b.request.Network = &network
	return b
}

// WithMetadata sets custom metadata
func (b *ContainerRequestBuilder) WithMetadata(metadata map[string]string) *ContainerRequestBuilder {
	b.request.Metadata = metadata
	return b
}

// AddMetadata adds a single metadata key-value pair
func (b *ContainerRequestBuilder) AddMetadata(key, value string) *ContainerRequestBuilder {
	if b.request.Metadata == nil {
		b.request.Metadata = make(map[string]string)
	}
	b.request.Metadata[key] = value
	return b
}

// Build validates and returns the orchestration request
func (b *ContainerRequestBuilder) Build() (*OrchestrationRequest, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}
	return b.request, nil
}

// MustBuild returns the orchestration request, panics on validation error
func (b *ContainerRequestBuilder) MustBuild() *OrchestrationRequest {
	req, err := b.Build()
	if err != nil {
		panic(err)
	}
	return req
}

func (b *ContainerRequestBuilder) validate() error {
	if b.request.Name == "" {
		return ErrNameRequired
	}
	if b.request.TemplateID == "" {
		return ErrTemplateIDRequired
	}
	if b.request.ResourcePoolID == "" {
		return ErrResourcePoolIDRequired
	}
	return nil
}

// ===== Emulator Request Builder =====

// NewEmulatorRequest creates a new Emulator request builder
func NewEmulatorRequest(name string) *EmulatorRequestBuilder {
	return &EmulatorRequestBuilder{
		RequestBuilder: RequestBuilder{
			request: &OrchestrationRequest{
				Name:             name,
				JobProvisionType: ProvisionTypeEmulator,
				Configuration:    DefaultConfiguration,
				Priority:         DefaultPriority,
				Schedule:         DefaultSchedule,
				Resources:        DefaultEmulatorResources,
			},
		},
	}
}

// WithTemplate sets the template ID
func (b *EmulatorRequestBuilder) WithTemplate(templateID string) *EmulatorRequestBuilder {
	b.request.TemplateID = templateID
	return b
}

// WithResourcePool sets the resource pool ID
func (b *EmulatorRequestBuilder) WithResourcePool(resourcePoolID string) *EmulatorRequestBuilder {
	b.request.ResourcePoolID = resourcePoolID
	return b
}

// WithResources sets the resource allocation
func (b *EmulatorRequestBuilder) WithResources(cpu int, memory, disk string) *EmulatorRequestBuilder {
	b.request.Resources = Resource{
		CPU:    cpu,
		Memory: memory,
		Disk:   disk,
	}
	return b
}

// WithConfiguration sets the infrastructure configuration
func (b *EmulatorRequestBuilder) WithConfiguration(config Configuration) *EmulatorRequestBuilder {
	b.request.Configuration = config
	return b
}

// WithPriority sets the job priority
func (b *EmulatorRequestBuilder) WithPriority(priority Priority) *EmulatorRequestBuilder {
	b.request.Priority = priority
	return b
}

// WithMetadata sets custom metadata
func (b *EmulatorRequestBuilder) WithMetadata(metadata map[string]string) *EmulatorRequestBuilder {
	b.request.Metadata = metadata
	return b
}

// Build validates and returns the orchestration request
func (b *EmulatorRequestBuilder) Build() (*OrchestrationRequest, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}
	return b.request, nil
}

// MustBuild returns the orchestration request, panics on validation error
func (b *EmulatorRequestBuilder) MustBuild() *OrchestrationRequest {
	req, err := b.Build()
	if err != nil {
		panic(err)
	}
	return req
}

func (b *EmulatorRequestBuilder) validate() error {
	if b.request.Name == "" {
		return ErrNameRequired
	}
	if b.request.TemplateID == "" {
		return ErrTemplateIDRequired
	}
	if b.request.ResourcePoolID == "" {
		return ErrResourcePoolIDRequired
	}
	return nil
}
