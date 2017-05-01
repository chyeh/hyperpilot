package main

// RunningResourceUsage mocks the information of the process's deployment from the system
var RunningResourceUsage = make(map[int]ResourceUsage)

var ResourceUpperBound = make(map[int]ResourceUsage)

// mock function for the performance
func runLoad(cpu float64, mem float64) Percentiles {
	// Assumption: Response time increases only if memory ls less than 512 mb
	// Assumption: Response time increases only if CPU ls less than 500 m-core
	if mem >= 512.0 {
		return Percentiles{
			FiftiethPercentile:    200 + int(1.0*(500.0-cpu)),
			NinetyFifthPercentile: 250 + int(1.2*(500.0-cpu)),
			NinetyNinthPercentile: 350 + int(1.5*(500.0-cpu)),
		}
	}
	return Percentiles{
		FiftiethPercentile:    200 + int(1.0*(500.0-cpu)) + int(30.0*(512.0-mem)),
		NinetyFifthPercentile: 250 + int(1.2*(500.0-cpu)) + int(50.0*(512.0-mem)),
		NinetyNinthPercentile: 350 + int(1.5*(500.0-cpu)) + int(105.0*(512.0-mem)),
	}
}

// mock the usage of the deployment
func deploy(sID int, cpu float64, mem float64) {
	// Assumption: All the services use 512 mb at most
	if mem > 512.0 {
		mem = 512.0
	}
	// Assumption: All the services use 2000 m-core at most
	if cpu > 2000.0 {
		cpu = 2000.0
	}
	RunningResourceUsage[sID] = ResourceUsage{
		CPUUsage: cpu,
		MemUsage: mem,
	}
}
