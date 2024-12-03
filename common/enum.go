package common

import (
	"encoding/json"
	"errors"
)

type ResourceType string

const (
	CPU ResourceType = "1"
	GPU ResourceType = "2"
)

func (r *ResourceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case string(CPU), string(GPU):
		*r = ResourceType(s)
		return nil
	}
	return errors.New("invalid resource type")
}

type GpuSKUSpec struct {
	GpuCount   int32
	CpuCount   int32
	MemorySize int32
}

type GpuSKU string
type GpuCount string
type GpuMapping string

const (
	V100_16GB_PCIE    GpuSKU     = "V100_16GB_PCIE"
	A100_80GB_PCIE    GpuSKU     = "A100_80GB_PCIE"
	RTX4090_24GB_PCIE GpuSKU     = "RTX4090_24GB_PCIE"
	GPU_1             GpuCount   = "GPU_1"
	GPU_2             GpuCount   = "GPU_2"
	GPU_4             GpuCount   = "GPU_4"
	GPU_8             GpuCount   = "GPU_8"
	GPU_16            GpuCount   = "GPU_16"
	GPU_32            GpuCount   = "GPU_32"
	GPU_Mapping_1     GpuMapping = "1G_1C_1M" // g4dn.xlarge 1 GPU 4vCPU 16GB
	GPU_Mapping_2     GpuMapping = "1G_1C_2M" // g4dn.2xlarge 1 GPU 8vCPU 32GB
	GPU_Mapping_3     GpuMapping = "1G_2C_2M" // g4dn.12large 4 GPU 48vCPU 192GB
	GPU_Mapping_4     GpuMapping = "1G_2C_4M" // g4dn.12large 4 GPU 48vCPU 192GB
	GPU_Mapping_5     GpuMapping = "1G_4C_8M" // g4dn.metal 8 GPU 96vCPU 384GB
	//GPU_Mapping_1        = "1:3:14" // g4dn.xlarge 1 GPU 4vCPU 16GB
	//GPU_Mapping_2        = "1:6:25" // g4dn.2xlarge 1 GPU 8vCPU 32GB
	//GPU_Mapping_3        = "1:9:38" // g4dn.12large 4 GPU 48vCPU 192GB
	//GPU_Mapping_4        = "1:9:38" // g4dn.metal 8 GPU 96vCPU 384GB
)

func (r *GpuSKU) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "" {
		return nil
	}

	switch s {
	case string(V100_16GB_PCIE), string(A100_80GB_PCIE), string(RTX4090_24GB_PCIE):
		*r = GpuSKU(s)
		return nil
	}
	return errors.New("invalid gpu sku")
}

func (r *GpuCount) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "" {
		return nil
	}

	switch s {
	case string(GPU_1), string(GPU_2), string(GPU_4), string(GPU_8), string(GPU_16), string(GPU_32):
		*r = GpuCount(s)
		return nil
	}
	return errors.New("invalid gpu count")
}

func (r *GpuMapping) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "" {
		return nil
	}

	switch s {
	case string(GPU_Mapping_1), string(GPU_Mapping_2), string(GPU_Mapping_3), string(GPU_Mapping_4), string(GPU_Mapping_5):
		*r = GpuMapping(s)
		return nil
	}
	return errors.New("invalid gpu mapping")
}
