package types

// WasmConfig is the extra config required for wasm
type WasmConfig struct {
	SmartQueryGasLimit uint64 `mapstructure:"query_gas_limit"`
	CacheSize          uint64 `mapstructure:"lru_size"`
}

type TrustlessMsg struct {
	CodeHash []byte
	Msg      []byte
}

type ParseAuto struct {
	AutoMsg struct {
	} `json:"auto_msg"`
}

type EstimateResult struct {
	Estimation struct {
		Status     string `json:"status"`
		TotalCount int    `json:"total_count"`
		//ReadyForReveal bool   `json:"ready_for_reveal"`
		//AmountEstimation string `json:"amount_estimation"`
	} `json:"estimation"`
}

type RevealResult struct {
	RevealEstimation struct {
		Status         string   `json:"status"`
		Message        string   `json:"message"`
		BestEstimation int      `json:"best_estimation"`
		Comments       []string `json:"comments"`
		BestEstimator  string   `json:"best_estimator"`
		EstimationList []int    `json:"estimation_list"`
	} `json:"reveal_estimation"`
}

type ParseReveal struct {
	RevealEstimation struct {
	} `json:"reveal_estimation"`
}

type ParseTransferable struct {
	Transferable struct {
	} `json:"transferable"`
}

type TransferableResult struct {
	Transferable struct {
		Status string `json:"status"`
	} `json:"transferable"`
}

type StatusResult struct {
	StatusOnly struct {
		Status string `json:"status"`
	} `json:"status_only"`
}

type ParseFlag struct {
	Flag struct {
	} `json:"flag"`
}

type ParseDelete struct {
	RetractEstimation struct {
	} `json:"retract_estimation"`
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWasmConfig() WasmConfig {
	return WasmConfig{
		SmartQueryGasLimit: defaultQueryGasLimit,
		CacheSize:          defaultLRUCacheSize,
	}
}

func (m TrustlessMsg) Serialize() []byte {
	return append(m.CodeHash, m.Msg...)
}

const defaultLRUCacheSize = uint64(0)
const defaultQueryGasLimit = uint64(3000000)
