package main

type Service struct {
	CPURscCfg CPUResourceConfig
	MemRscCfg MemoryResourceConfig
}

type CPUResourceConfig struct {
	Min float64
	Max float64
}

type MemoryResourceConfig struct {
	Min float64
	Max float64
}

var Services = map[int]Service{
	1: {
		CPURscCfg: CPUResourceConfig{
			Min: 0.0,
			Max: 500.0,
		},
		MemRscCfg: MemoryResourceConfig{
			Min: 0.0,
			Max: 1024.0,
		},
	},
	2: {
		CPURscCfg: CPUResourceConfig{
			Min: 0.0,
			Max: 500.0,
		},
		MemRscCfg: MemoryResourceConfig{
			Min: 0.0,
			Max: 1024.0,
		},
	},
}

type SLAConfig struct {
	NinetyFifthPercentile int
	NinetyNinthPercentile int
}

var SLA = SLAConfig{
	NinetyFifthPercentile: 750,
	NinetyNinthPercentile: 1000,
}
