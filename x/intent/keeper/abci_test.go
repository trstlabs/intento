package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestShouldRecur(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	k := keepers.IntentKeeper

	now := time.Now().UTC()
	interval := 5 * time.Minute

	tests := []struct {
		name           string
		setupFlow      func() types.Flow
		errorString    string
		hasHistory     bool
		hasHistoryErr  bool
		expectedResult bool
	}{
		{
			name: "balance too low - should not recur",
			setupFlow: func() types.Flow {
				return types.Flow{
					ID:       4,
					ExecTime: now,
					EndTime:  now.Add(1 * time.Hour),
					Interval: interval,
					Configuration: &types.ExecutionConfiguration{
						StopOnSuccess: false,
						StopOnFailure: false,
					},
				}
			},
			errorString:    types.ErrBalanceTooLow,
			expectedResult: false,
		},
		{
			name: "no conditions - should recur",
			setupFlow: func() types.Flow {
				return types.Flow{
					ID:       5,
					ExecTime: now,
					EndTime:  now.Add(1 * time.Hour),
					Interval: interval,
					Configuration: &types.ExecutionConfiguration{
						StopOnSuccess: false,
						StopOnFailure: false,
					},
				}
			},
			errorString:    "",
			expectedResult: true,
		},
		{
			name: "stop on success - with error - should recur",
			setupFlow: func() types.Flow {
				return types.Flow{
					ExecTime: now,
					EndTime:  now.Add(1 * time.Hour),
					Interval: interval,
					Configuration: &types.ExecutionConfiguration{
						StopOnSuccess: true,
						StopOnFailure: false,
					},
				}
			},
			errorString:    "some error",
			hasHistory:     true,
			hasHistoryErr:  true,
			expectedResult: true,
		},
		{
			name: "stop on success - no error - should not recur",
			setupFlow: func() types.Flow {
				return types.Flow{
					ID:       1,
					ExecTime: now,
					EndTime:  now.Add(1 * time.Hour),
					Interval: interval,
					Configuration: &types.ExecutionConfiguration{
						StopOnSuccess: true,
						StopOnFailure: false,
					},
				}
			},
			errorString:    "",
			hasHistory:     false,
			hasHistoryErr:  false,
			expectedResult: false, // Should not recur because there's no error and StopOnSuccess is true
		},
		{
			name: "stop on failure - with error - should not recur",
			setupFlow: func() types.Flow {
				return types.Flow{
					ExecTime: now,
					EndTime:  now.Add(1 * time.Hour),
					Interval: interval,
					Configuration: &types.ExecutionConfiguration{
						StopOnSuccess: false,
						StopOnFailure: true,
					},
				}
			},
			errorString:    "some error",
			hasHistory:     true,
			hasHistoryErr:  true,
			expectedResult: false,
		},
		{
			name: "stop on failure - no error - should recur",
			setupFlow: func() types.Flow {
				return types.Flow{
					ID:       2,
					ExecTime: now,
					EndTime:  now.Add(1 * time.Hour),
					Interval: interval,
					Configuration: &types.ExecutionConfiguration{
						StopOnSuccess: false,
						StopOnFailure: true,
					},
				}
			},
			errorString:    "",
			hasHistory:     false,
			hasHistoryErr:  false,
			expectedResult: true, // Should recur because there's no error and StopOnFailure is true
		},
		{
			name: "end time reached - should not recur",
			setupFlow: func() types.Flow {
				// Set ExecTime and EndTime so that the next execution would be after EndTime
				nearEndTime := now.Add(-1 * time.Minute)
				return types.Flow{
					ID:       3,
					ExecTime: nearEndTime,
					EndTime:  now,
					Interval: 2 * time.Minute, // Next execution would be at now + 1min, which is after EndTime
					Configuration: &types.ExecutionConfiguration{
						StopOnSuccess: false,
						StopOnFailure: false,
					},
				}
			},
			errorString:    "",
			expectedResult: false, // Should not recur because next execution would be after EndTime
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flow := tc.setupFlow()

			// Set up history if needed
			if tc.hasHistory {
				errMsg := ""
				if tc.hasHistoryErr {
					errMsg = "history error"
				}
				history := types.FlowHistoryEntry{
					Errors: []string{errMsg},
				}
				k.SetFlowHistoryEntry(ctx, flow.ID, &history)
			}

			result := k.shouldRecur(ctx, flow, tc.errorString)
			require.Equal(t, tc.expectedResult, result, "test case: %s", tc.name)
		})
	}
}
