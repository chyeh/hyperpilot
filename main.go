package main

import "fmt"

func DeployService(serviceId int, cpu float64, mem float64) {
	fmt.Printf("[DEPLOY] service[%d]: cpu[%f] mem[%f]\n", serviceId, cpu, mem)
	deploy(serviceId, cpu, mem)
	return
}

type Percentiles struct {
	FiftiethPercentile    int
	NinetyFifthPercentile int
	NinetyNinthPercentile int
}

func RunLoadTest() Percentiles {
	var ans Percentiles
	for _, u := range RunningResourceUsage {
		p := runLoad(u.CPUUsage, u.MemUsage)
		ans.FiftiethPercentile = ans.FiftiethPercentile + p.FiftiethPercentile
		ans.NinetyFifthPercentile = ans.NinetyFifthPercentile + p.NinetyFifthPercentile
		ans.NinetyNinthPercentile = ans.NinetyNinthPercentile + p.NinetyNinthPercentile
	}
	return ans
}

type ResourceUsage struct {
	CPUUsage float64
	MemUsage float64
}

func GetResourceUsagesFromLastDeployment(serverId int) ResourceUsage {
	fmt.Printf("[USAGE] service[%d]: cpu[%f] mem[%f]\n", serverId, RunningResourceUsage[serverId].CPUUsage, RunningResourceUsage[serverId].MemUsage)
	return RunningResourceUsage[serverId]
}

func isSLAMet(result Percentiles) bool {
	switch {
	case
		result.NinetyFifthPercentile > SLA.NinetyFifthPercentile,
		result.NinetyNinthPercentile > SLA.NinetyNinthPercentile:
		return false
	}
	return true
}

func improve(curr map[int]ResourceUsage, got map[int]ResourceUsage) map[int]ResourceUsage {
	currSum, gotSum := 0.0, 0.0
	for _, u := range curr {
		currSum = currSum + u.CPUUsage + u.MemUsage
	}
	for _, u := range got {
		gotSum = gotSum + u.CPUUsage + u.MemUsage
	}

	if currSum > gotSum {
		fmt.Println("== Solution Improved ==", got)
		return got
	}
	return curr
}

func initDeploy(svcs map[int]Service) {
	fmt.Println("== INI DEPLOY ==")
	for k, v := range svcs {
		DeployService(k, v.CPURscCfg.Max, v.MemRscCfg.Max)
		u := GetResourceUsagesFromLastDeployment(k)
		ResourceUpperBound[k] = setUpperBound(u, v.CPURscCfg, v.MemRscCfg)
	}
}

func setUpperBound(u ResourceUsage, cCfg CPUResourceConfig, mCfg MemoryResourceConfig) ResourceUsage {
	var ans = ResourceUsage{
		CPUUsage: cCfg.Max,
		MemUsage: mCfg.Max,
	}
	if u.CPUUsage < cCfg.Max {
		ans.CPUUsage = u.CPUUsage
	}
	if u.MemUsage < mCfg.Max {
		ans.MemUsage = u.MemUsage
	}
	return ans
}

func search(svcs map[int]Service) (cpuResult float64, memResult float64) {
	fmt.Println("== SEARCH ==")
	// Assumption: Resource configuration of all services are identical
	iCPU, iMem, jCPU, jMem := Services[1].CPURscCfg.Min, Services[1].MemRscCfg.Min, ResourceUpperBound[1].CPUUsage, ResourceUpperBound[1].MemUsage
	for iCPU < jCPU && iMem < jMem {
		hCPU := iCPU + (jCPU-iCPU)/2.0
		hMem := iMem + (jMem-iMem)/2.0

		for k, _ := range svcs {
			DeployService(k, hCPU, hMem)
			GetResourceUsagesFromLastDeployment(k)
		}
		p := RunLoadTest()
		if !isSLAMet(p) {
			iCPU, iMem = hCPU+10.0, hMem+10.0
		} else {
			jCPU, jMem = hCPU, hMem
		}
	}
	return iCPU, iMem
}

func myAlgo(svcs map[int]Service) map[int]ResourceUsage {
	initDeploy(svcs)
	r := RunLoadTest()
	if !isSLAMet(r) {
		return nil
	}

	// Use binary search to find out the first valid solution
	originCPU, originMem := search(svcs)
	fmt.Println("cpu, mem", originCPU, originMem)

	ans := make(map[int]ResourceUsage)
	for k, _ := range svcs {
		ans[k] = ResourceUsage{
			CPUUsage: originCPU,
			MemUsage: originMem,
		}
	}

	// Then try to improve the answer by reducing the memory usage
	fmt.Println("== CPU Bound ==")
	currCPU, currMem := originCPU, originMem
	for currMem > Services[1].MemRscCfg.Min && currCPU < ResourceUpperBound[1].CPUUsage+10.0 {
		nextCPU := currCPU
		nextMem := currMem - 10.0
		for k, _ := range svcs {
			DeployService(k, nextCPU, nextMem)
		}
		p := RunLoadTest()
		if !isSLAMet(p) {
			currCPU, currMem = nextCPU+10.0, nextMem+10.0
		} else {
			currCPU, currMem = nextCPU, nextMem
			got := make(map[int]ResourceUsage)
			for k, _ := range svcs {
				got[k] = ResourceUsage{
					CPUUsage: GetResourceUsagesFromLastDeployment(k).CPUUsage,
					MemUsage: GetResourceUsagesFromLastDeployment(k).MemUsage,
				}
			}
			ans = improve(ans, got)
		}
	}

	// Then try to improve the answer by reducing the CPU usage
	fmt.Println("== Memory Bound ==")
	currCPU, currMem = originCPU, originMem
	for currCPU > Services[1].CPURscCfg.Min && currMem < ResourceUpperBound[1].MemUsage+10.0 {
		nextCPU := currCPU - 10.0
		nextMem := currMem
		for k, _ := range svcs {
			DeployService(k, nextCPU, nextMem)
		}
		p := RunLoadTest()
		if !isSLAMet(p) {
			currCPU, currMem = nextCPU+10.0, nextMem+10.0
		} else {
			currCPU, currMem = nextCPU, nextMem
			got := make(map[int]ResourceUsage)
			for k, _ := range svcs {
				got[k] = ResourceUsage{
					CPUUsage: GetResourceUsagesFromLastDeployment(k).CPUUsage,
					MemUsage: GetResourceUsagesFromLastDeployment(k).MemUsage,
				}
			}
			ans = improve(ans, got)
		}
	}

	fmt.Println("ans:", ans)
	return ans
}

func main() {
	myAlgo(Services)
}
