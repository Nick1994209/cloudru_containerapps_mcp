package domain

// Credentials represents the authentication credentials for Cloud.ru
type Credentials struct {
	RegistryName string
	KeyID        string
	KeySecret    string
}

// DockerImage represents a Docker image to be built and pushed
type DockerImage struct {
	RegistryName     string
	RepositoryName   string
	ImageVersion     string
	DockerfilePath   string
	DockerfileTarget string
	DockerfileFolder string
}

// ContainerApp represents a Cloud.ru Container App
type ContainerApp struct {
	ProjectID     string `json:"projectId"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	CreatedBy     string `json:"createdBy"`
	UpdatedBy     string `json:"updatedBy"`
	Configuration struct {
		Ingress struct {
			PubliclyAccessible     bool          `json:"publiclyAccessible"`
			PublicUri              string        `json:"publicUri"`
			InternalUri            string        `json:"internalUri"`
			AdditionalPortMappings []interface{} `json:"additionalPortMappings"`
		} `json:"ingress"`
		AutoDeployments struct {
			Enabled bool   `json:"enabled"`
			Pattern string `json:"pattern"`
		} `json:"autoDeployments"`
		Privileged bool `json:"privileged"`
	} `json:"configuration"`
	Template struct {
		Timeout     string `json:"timeout"`
		IdleTimeout string `json:"idleTimeout"`
		Protocol    string `json:"protocol"`
		Scaling     struct {
			MinInstanceCount int `json:"minInstanceCount"`
			MaxInstanceCount int `json:"maxInstanceCount"`
			Rule             struct {
				Type  string `json:"type"`
				Value struct {
					Soft int `json:"soft"`
					Hard int `json:"hard"`
				} `json:"value"`
			} `json:"rule"`
		} `json:"scaling"`
		Containers []struct {
			Name      string `json:"name"`
			Image     string `json:"image"`
			Resources struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"resources"`
			ContainerPort int `json:"containerPort"`
			Env           []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Type  string `json:"type,omitempty"`
			} `json:"env"`
			Command      []interface{} `json:"command"`
			Args         []interface{} `json:"args"`
			VolumeMounts []struct {
				Name      string `json:"name"`
				MountPath string `json:"mountPath"`
				ReadOnly  bool   `json:"readOnly"`
			} `json:"volumeMounts"`
		} `json:"containers"`
		InitContainers []interface{} `json:"initContainers"`
		Volumes        []struct {
			Name             string `json:"name"`
			Type             string `json:"type"`
			VolumeAttributes struct {
				BucketName string `json:"bucketName"`
				TenantId   string `json:"tenantId"`
				Region     string `json:"region"`
				ReadOnly   string `json:"readOnly"`
				Entrypoint string `json:"entrypoint"`
			} `json:"volumeAttributes"`
		} `json:"volumes"`
	} `json:"template"`
}
