// batch.go 提供批量仿真运行和参数扫描功能。
package sim

// SweepResult 参数扫描结果。
type SweepResult struct {
	ParamName  string       `json:"param_name"`
	ParamValues []float64   `json:"param_values"`
	Results    []Stats      `json:"results"`
}

// SweepLambda 对到达率做参数扫描。
func SweepLambda(base Config, lambdas []float64) (*SweepResult, error) {
	res := &SweepResult{ParamName: "lambda", ParamValues: lambdas}
	for _, l := range lambdas {
		c := base
		c.Lambda = l
		var st Stats
		var err error
		if c.Servers == 1 {
			st, err = RunMM1(c)
		} else {
			st, err = RunMMc(c, c.Servers)
		}
		if err != nil {
			st = Stats{} // skip failures but record zero
		}
		res.Results = append(res.Results, st)
	}
	return res, nil
}

// SweepServers 对服务台数做参数扫描。
func SweepServers(base Config, serverCounts []int) (*SweepResult, error) {
	res := &SweepResult{ParamName: "servers"}
	for _, s := range serverCounts {
		res.ParamValues = append(res.ParamValues, float64(s))
		c := base
		c.Servers = s
		st, err := RunMMc(c, s)
		if err != nil {
			st = Stats{}
		}
		res.Results = append(res.Results, st)
	}
	return res, nil
}

// CompareConfigs 对比多组配置的仿真结果。
type CompareResult struct {
	Configs []Config `json:"configs"`
	Stats   []Stats  `json:"stats"`
}

// BatchRun 批量运行多组配置。
func BatchRun(configs []Config) (*CompareResult, error) {
	cr := &CompareResult{Configs: configs}
	for _, c := range configs {
		var st Stats
		var err error
		if c.Servers == 1 {
			st, err = RunMM1(c)
		} else {
			st, err = RunMMc(c, c.Servers)
		}
		if err != nil {
			st = Stats{}
		}
		cr.Stats = append(cr.Stats, st)
	}
	return cr, nil
}

// WhatIf 模拟「如果增加 n 个服务台」的效果。
func WhatIf(base Config, addServers int) (Stats, Stats, error) {
	before, err := RunMMc(base, base.Servers)
	if err != nil {
		return Stats{}, Stats{}, err
	}
	c := base
	c.Servers += addServers
	after, err := RunMMc(c, c.Servers)
	if err != nil {
		return before, Stats{}, err
	}
	return before, after, nil
}
