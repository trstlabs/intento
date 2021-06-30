package types

// WasmConfig is the extra config required for wasm
type WasmConfig struct {
	SmartQueryGasLimit uint64 `mapstructure:"query_gas_limit"`
	CacheSize          uint64 `mapstructure:"lru_size"`
}

type SecretMsg struct {
	CodeHash []byte
	Msg      []byte
}

type EstimateResult struct {
	Estimation struct {
		Status         string `json:"status"`
		TotalCount     int    `json:"total_count"`
		ReadyForReveal bool   `json:"ready_for_reveal"`
		//AmountEstimation string `json:"amount_estimation"`
	} `json:"estimation"`
}

type RevealResult struct {
	RevealEstimation struct {
		Status         string   `json:"status"`
		Message        string   `json:"message"`
		Bestestimation int      `json:"best_estimation"`
		Comments       []string `json:"comments"`
		Bestestimator  string   `json:"best_estimator"`
	} `json:"reveal_estimation"`
}

type ParseReveal struct {
	RevealEstimation struct {
	} `json:"reveal_estimation"`
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWasmConfig() WasmConfig {
	return WasmConfig{
		SmartQueryGasLimit: defaultQueryGasLimit,
		CacheSize:          defaultLRUCacheSize,
	}
}

func (m SecretMsg) Serialize() []byte {
	return append(m.CodeHash, m.Msg...)
}

const defaultLRUCacheSize = uint64(0)
const defaultQueryGasLimit = uint64(3000000)
