package keeper

import (
	"encoding/json"
	"fmt"

	"reflect"
	"strings"

	"cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/gogoproto/proto"
	osmosistwapv1beta1 "github.com/trstlabs/intento/x/intent/msg_registry/osmosis/twap/v1beta1"
	"github.com/trstlabs/intento/x/intent/types"
)

// CompareResponseValue compares the value of a response key based on the ResponseComparison
func (k Keeper) CompareResponseValue(ctx sdk.Context, flowID uint64, responses []*cdctypes.Any, comparison types.Comparison) (bool, error) {
	k.Logger(ctx).Debug("response comparison", "flowID", flowID)
	var queryCallback []byte = nil
	if comparison.ICQConfig != nil {
		queryCallback = comparison.ICQConfig.Response
	}
	if comparison.ResponseKey == "" && queryCallback == nil {
		return true, nil
	}
	if len(responses) <= int(comparison.ResponseIndex) && queryCallback == nil {
		return false, fmt.Errorf("not enough message responses to compare to, number of responses: %v", len(responses))
	}
	var valueFromResponse, previousValue interface{}
	var err error

	// Get the previous value if difference_mode is enabled
	if comparison.DifferenceMode {
		if previousValue, err = k.getPreviousResponseValue(ctx, flowID, comparison.ResponseIndex, comparison.ResponseKey, comparison.ValueType, queryCallback != nil); err != nil {
			k.Logger(ctx).Error("error getting previous response value", "error", err)
			return false, fmt.Errorf("error getting previous response value: %v", err)
		}
	}

	valueFromResponse, err = k.getValueFromResponse(
		ctx,
		comparison.ResponseIndex,
		comparison.ResponseKey,
		&comparison.ValueType,
		queryCallback,
		responses,
	)
	if err != nil {
		return false, fmt.Errorf("error getting value from response: %w", err)
	}

	operand, err := parseOperand(comparison.Operand, comparison.ValueType)
	if err != nil {
		valueType := comparison.Operand
		k.Logger(ctx).Debug("operand is not a number, trying to get value from response", "operand", comparison.Operand, "valueType", comparison.ValueType)
		// try to use operand field to get value from response
		operand, err = k.getValueFromResponse(
			ctx,
			comparison.ResponseIndex,
			comparison.ResponseKey,
			&valueType,
			queryCallback,
			responses,
		)
		if err != nil {
			return false, fmt.Errorf("error parsing operand: %v", err)
		}
	}

	if comparison.DifferenceMode {
		k.Logger(ctx).Debug("difference mode enabled", "valueFromResponse", valueFromResponse, "previousValue", previousValue)
		diffValue, err := k.calculateDifference(ctx, valueFromResponse, previousValue)
		if err != nil {
			return false, fmt.Errorf("error calculating difference: %w", err)
		}
		valueFromResponse = diffValue
	}

	// Normalize types for comparison
	valueFromResponse, operand = normalizeJSONTypes(valueFromResponse, operand)

	k.Logger(ctx).Debug("Comparing value", "valueFromResponse", valueFromResponse, "operand", operand, "operator", comparison.Operator)
	switch comparison.Operator {
	case types.ComparisonOperator_EQUAL:
		switch val := valueFromResponse.(type) {
		case math.Dec:
			operandDec, ok := operand.(math.Dec)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Dec")
			}
			return val.Equal(operandDec), nil

		case math.Int:
			operandInt, ok := operand.(math.Int)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Int")
			}
			return val.Equal(operandInt), nil

		default:
			return reflect.DeepEqual(valueFromResponse, operand), nil
		}

	case types.ComparisonOperator_NOT_EQUAL:
		switch val := valueFromResponse.(type) {
		case math.Dec:
			operandDec, ok := operand.(math.Dec)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Dec")
			}
			return !val.Equal(operandDec), nil

		case math.Int:
			operandInt, ok := operand.(math.Int)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Int")
			}
			return !val.Equal(operandInt), nil

		default:
			return !reflect.DeepEqual(valueFromResponse, operand), nil
		}
	case types.ComparisonOperator_SMALLER_THAN:
		if valDec, ok := valueFromResponse.(math.Dec); ok {
			operandDec, ok := operand.(math.Dec)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Dec")
			}
			subAmount, err := valDec.Sub(operandDec)
			if err != nil {
				return false, fmt.Errorf("error subtracting operand from value: %v", err)
			}
			return subAmount.IsNegative(), nil
		}
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a < b })

	case types.ComparisonOperator_LARGER_THAN:
		if valDec, ok := valueFromResponse.(math.Dec); ok {
			operandDec, ok := operand.(math.Dec)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Dec")
			}
			subAmount, err := valDec.Sub(operandDec)
			if err != nil {
				return false, fmt.Errorf("error subtracting operand from value: %v", err)
			}
			return subAmount.IsPositive(), nil
		}
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a > b })

	case types.ComparisonOperator_GREATER_EQUAL:
		if valDec, ok := valueFromResponse.(math.Dec); ok {
			operandDec, ok := operand.(math.Dec)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Dec")
			}
			subAmount, err := valDec.Sub(operandDec)
			if err != nil {
				return false, fmt.Errorf("error subtracting operand from value: %v", err)
			}
			return !subAmount.IsNegative(), nil
		}
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a >= b })

	case types.ComparisonOperator_LESS_EQUAL:
		if valDec, ok := valueFromResponse.(math.Dec); ok {
			operandDec, ok := operand.(math.Dec)
			if !ok {
				return false, fmt.Errorf("operand is not of type math.Dec")
			}
			subAmount, err := valDec.Sub(operandDec)
			if err != nil {
				return false, fmt.Errorf("error subtracting operand from value: %v", err)
			}
			return !subAmount.IsPositive(), nil
		}
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a <= b })
	case types.ComparisonOperator_CONTAINS:
		return contains(valueFromResponse, operand), nil
	case types.ComparisonOperator_NOT_CONTAINS:
		return !contains(valueFromResponse, operand), nil
	case types.ComparisonOperator_STARTS_WITH:
		return strings.HasPrefix(fmt.Sprintf("%v", valueFromResponse), fmt.Sprintf("%v", operand)), nil
	case types.ComparisonOperator_ENDS_WITH:
		return strings.HasSuffix(fmt.Sprintf("%v", valueFromResponse), fmt.Sprintf("%v", operand)), nil

	default:
		return false, fmt.Errorf("unsupported comparison operator: %v", comparison.Operator)
	}
}

// FeedbackLoop replaces the value in a message with the value from a response
func (k Keeper) RunFeedbackLoops(ctx sdk.Context, flowID uint64, msgs *[]*cdctypes.Any, conditions *types.ExecutionConditions) error {
	k.Logger(ctx).Debug("\n=== Starting RunFeedbackLoops ===")
	k.Logger(ctx).Debug("Total messages in flow: %d\n", len(*msgs))

	if conditions == nil || len(conditions.FeedbackLoops) == 0 {
		k.Logger(ctx).Debug("No feedback loops to process")
		return nil
	}

	// Process each feedback loop
	for _, feedbackLoop := range conditions.FeedbackLoops {
		k.Logger(ctx).Debug("\n--- Processing feedback loop ---")
		k.Logger(ctx).Debug("Message index", feedbackLoop.MsgsIndex)
		if int(feedbackLoop.MsgsIndex) < len(*msgs) {
			k.Logger(ctx).Debug("Message type URL", (*msgs)[feedbackLoop.MsgsIndex].TypeUrl)
		}
		k.Logger(ctx).Debug("Response index", feedbackLoop.ResponseIndex)
		k.Logger(ctx).Debug("Response key", feedbackLoop.ResponseKey)
		k.Logger(ctx).Debug("Message key", feedbackLoop.MsgKey)

		// Skip if the message index is out of range
		if int(feedbackLoop.MsgsIndex) >= len(*msgs) {
			k.Logger(ctx).Debug("Skipping feedback loop - message index %d out of range (max: %d)\n",
				feedbackLoop.MsgsIndex, len(*msgs)-1)
			continue
		}

		var queryCallback []byte = nil
		if feedbackLoop.ICQConfig != nil {
			queryCallback = feedbackLoop.ICQConfig.Response
		}

		k.Logger(ctx).Debug("running feedback loop", "flowID", flowID, "queryCallback", queryCallback)

		if feedbackLoop.ResponseKey == "" && queryCallback == nil {
			return nil
		}

		if feedbackLoop.FlowID != 0 {
			flowID = feedbackLoop.FlowID

		}
		var valueFromResponse, previousValue interface{}
		var err error

		// Get the previous value if difference_mode is enabled
		if feedbackLoop.DifferenceMode {
			if previousValue, err = k.getPreviousResponseValue(ctx, flowID, feedbackLoop.ResponseIndex, feedbackLoop.ResponseKey, feedbackLoop.ValueType, queryCallback != nil); err != nil {
				k.Logger(ctx).Error("error getting previous response value", "error", err)
				return fmt.Errorf("error getting previous response value: %v", err)
			}
		}

		var responses []*cdctypes.Any
		if queryCallback == nil {
			history, err := k.GetFlowHistory(ctx, flowID)
			if err != nil {
				return err
			}
			if len(history) == 0 {
				return nil
			}
			responses = history[len(history)-1].MsgResponses
			if len(responses) == 0 || int(feedbackLoop.ResponseIndex) >= len(responses) {
				continue
			}
		}

		valueFromResponse, err = k.getValueFromResponse(
			ctx,
			feedbackLoop.ResponseIndex,
			feedbackLoop.ResponseKey,
			&feedbackLoop.ValueType,
			queryCallback,
			responses,
		)
		if err != nil {
			return fmt.Errorf("error getting value from response: %w", err)
		}
		if feedbackLoop.DifferenceMode {
			diffValue, err := k.calculateDifference(ctx, valueFromResponse, previousValue)
			if err != nil {
				return fmt.Errorf("error calculating difference: %w", err)
			}
			valueFromResponse = diffValue
		}

		// Get the message to modify
		msgIndex := feedbackLoop.MsgsIndex

		// Debug log to verify correct index
		k.Logger(ctx).Debug("Processing feedback loop with MsgsIndex=%d (original=%d)\n",
			msgIndex, feedbackLoop.MsgsIndex)

		if int(msgIndex) >= len(*msgs) {
			return fmt.Errorf("message index %d out of range", msgIndex)
		}

		msgAny := (*msgs)[msgIndex]

		// Create a new Any message with the same type and value
		msgCopy := &cdctypes.Any{
			TypeUrl: msgAny.TypeUrl,
			Value:   append([]byte(nil), msgAny.Value...), // Safe copy of the value
		}

		k.Logger(ctx).Debug("Running feedback loop", "msgIndex", msgIndex, "typeUrl", msgAny.TypeUrl, "responseIndex", feedbackLoop.ResponseIndex, "responseKey", feedbackLoop.ResponseKey, "msgKey", feedbackLoop.MsgKey)

		// Handle authz wrapped messages
		var msgToInterface sdk.Msg
		var isWrapped bool
		k.Logger(ctx).Debug("Checking if message is wrapped (type: %s)...\n", msgCopy.TypeUrl)

		if msgCopy.TypeUrl == sdk.MsgTypeURL(&authztypes.MsgExec{}) {
			k.Logger(ctx).Debug("Message is wrapped in MsgExec, unwrapping...")
			msgExec := &authztypes.MsgExec{}
			if err := proto.Unmarshal(msgCopy.Value, msgExec); err != nil {
				return fmt.Errorf("failed to unmarshal MsgExec: %w", err)
			}
			if len(msgExec.Msgs) == 0 {
				return fmt.Errorf("no messages in MsgExec")
			}

			innerMsgCopy := &cdctypes.Any{
				TypeUrl: msgExec.Msgs[0].TypeUrl,
				Value:   append([]byte(nil), msgExec.Msgs[0].Value...),
			}

			var innerMsg sdk.Msg
			if err := k.cdc.UnpackAny(innerMsgCopy, &innerMsg); err != nil {
				return fmt.Errorf("failed to unpack inner message: %w", err)
			}

			msgToInterface = innerMsg
			isWrapped = true

		} else {
			k.Logger(ctx).Debug("Message is not wrapped, unpacking directly...")
			// For non-wrapped messages, unpack directly
			if err := k.cdc.UnpackAny(msgCopy, &msgToInterface); err != nil {
				return fmt.Errorf("failed to unpack message: %w", err)
			}
			k.Logger(ctx).Debug("Successfully unpacked message of type: %T\n", msgToInterface)
		}

		// Log the message type and value for debugging
		msgType := fmt.Sprintf("%T", msgToInterface)
		k.Logger(ctx).Debug("Processing message",
			"typeUrl", msgCopy.TypeUrl,
			"isWrapped", isWrapped,
			"msgType", msgType)

		// Get the reflect value of the message
		msgValue := reflect.ValueOf(msgToInterface)
		if msgValue.Kind() == reflect.Ptr {
			msgValue = msgValue.Elem()
		}
		k.Logger(ctx).Debug("Message value kind: %v, type: %v\n", msgValue.Kind(), msgValue.Type())

		// Ensure we're dealing with a struct
		if msgValue.Kind() != reflect.Struct {
			return fmt.Errorf("expected a struct, got %v", msgValue.Kind())
		}
		// Traverse the fields to find the one we need to update
		k.Logger(ctx).Debug("Looking for field: %s in message type: %T\n", feedbackLoop.MsgKey, msgToInterface)
		fieldToReplace, err := traverseFields(msgToInterface, feedbackLoop.MsgKey)
		if err != nil {
			k.Logger(ctx).Debug("Error traversing fields: %v\n", err)
			return fmt.Errorf("failed to traverse fields: %w", err)
		}

		k.Logger(ctx).Debug("Found field to replace",
			"field", fieldToReplace,
			"canSet", fieldToReplace.CanSet(),
			"currentValue", fieldToReplace.Interface(),
			"newValue", valueFromResponse)

		// Set the new value with proper type checking
		if !fieldToReplace.CanSet() {
			return fmt.Errorf("field %s cannot be set", feedbackLoop.MsgKey)
		}

		// Convert the value to the target type
		targetValue := reflect.ValueOf(valueFromResponse)
		if !targetValue.Type().AssignableTo(fieldToReplace.Type()) {
			// Try to convert the value if it's not directly assignable
			if targetValue.Type().ConvertibleTo(fieldToReplace.Type()) {
				targetValue = targetValue.Convert(fieldToReplace.Type())
			} else {
				return fmt.Errorf("cannot assign %s to %s",
					targetValue.Type(), fieldToReplace.Type())
			}
		}

		// Set the field value
		fieldToReplace.Set(targetValue)
		k.Logger(ctx).Debug("Successfully updated field", feedbackLoop.MsgKey, " to value", fieldToReplace.Interface())
		k.Logger(ctx).Debug("Updated field value",
			"field", feedbackLoop.MsgKey,
			"newValue", fieldToReplace.Interface())

		// Repack the message based on whether it was wrapped or not
		k.Logger(ctx).Debug("Repacking message", "isWrapped", isWrapped)

		if isWrapped {
			// For wrapped messages, we need to update the inner message in the MsgExec
			// Get the original MsgExec
			var msgExec authztypes.MsgExec
			if err := proto.Unmarshal(msgAny.Value, &msgExec); err != nil {
				return fmt.Errorf("failed to unpack MsgExec: %w", err)
			}

			// Update the inner message in the MsgExec
			msgExec.Msgs[0], err = cdctypes.NewAnyWithValue(msgToInterface)
			if err != nil {
				return fmt.Errorf("failed to create new Any for inner message: %w", err)
			}

			// Pack the updated MsgExec back to Any
			updatedMsgAny, err := cdctypes.NewAnyWithValue(&msgExec)
			if err != nil {
				return fmt.Errorf("failed to create new Any for MsgExec: %w", err)
			}

			// Replace the message in the slice
			(*msgs)[feedbackLoop.MsgsIndex] = updatedMsgAny
		} else {
			// For non-wrapped messages, create a new Any with the updated message
			updatedMsgAny, err := cdctypes.NewAnyWithValue(msgToInterface)
			if err != nil {
				return fmt.Errorf("failed to create new Any for message: %w", err)
			}
			(*msgs)[feedbackLoop.MsgsIndex] = updatedMsgAny
		}

		// Log the updated message for debugging
		k.Logger(ctx).Debug("Updated message", feedbackLoop.MsgsIndex, "TypeURL", (*msgs)[feedbackLoop.MsgsIndex].TypeUrl, "Value", (*msgs)[feedbackLoop.MsgsIndex].Value)
	}
	return nil
}

// getValueFromResponse extracts a value from either a regular response or an ICQ response
func (k Keeper) getValueFromResponse(
	ctx sdk.Context,
	responseIndex uint32,
	responseKey string,
	valueType *string,
	queryCallback []byte,
	responses []*cdctypes.Any,
) (interface{}, error) {
	var valueFromResponse interface{}
	var err error

	if queryCallback != nil {
		// Handle ICQ response
		valueFromResponse, err = parseICQResponse(queryCallback, valueType)
		if err != nil {
			var respAny cdctypes.Any
			err := k.cdc.Unmarshal(queryCallback, &respAny)
			if err != nil {
				var jsonObj map[string]interface{}
				if err := json.Unmarshal(queryCallback, &jsonObj); err != nil {
					k.Logger(ctx).Debug("response comparison: could not parse query callback data", "CallbackData", queryCallback)
					return nil, fmt.Errorf("response comparison: error parsing query callback data: %v, valueType: %s, error: %v", queryCallback, *valueType, err)
				}
				if responseKey != "" {
					valueFromResponse, err = extractJSONField(jsonObj, responseKey)
					if err != nil {
						return nil, fmt.Errorf("failed to extract field '%s' from JSON: %w", responseKey, err)
					}
				} else {
					valueFromResponse = jsonObj
				}

				valueFromResponse, err = parseResponseValue(valueFromResponse, responseKey, *valueType)
				if err != nil {
					return nil, fmt.Errorf("error parsing value: %v", err)
				}
			}
			if respAny.Value != nil {
				var resp interface{}
				err = k.cdc.UnpackAny(&respAny, &resp)
				if err != nil {
					return nil, fmt.Errorf("response comparison: error unpacking: %v", err)
				}
				valueFromResponse, err = parseResponseValue(queryCallback, responseKey, *valueType)
				if err != nil {
					return nil, fmt.Errorf("error parsing value: %v", err)
				}
			}
		}
	} else {
		var respAny cdctypes.Any
		if len(responses) == 0 {
			k.Logger(ctx).Debug("response comparison: could not parse response data")
		} else {
			respAny = *responses[responseIndex]
		}
		protoMsg, err := k.interfaceRegistry.Resolve(respAny.TypeUrl)
		if err != nil {
			return false, fmt.Errorf("failed to resolve type URL %s: %w", respAny.TypeUrl, err)
		}

		err = proto.Unmarshal(respAny.Value, protoMsg)
		if err != nil {
			return false, fmt.Errorf("error unmarshalling response: %v", err)
		}
		valueFromResponse, err = parseResponseValue(protoMsg, responseKey, *valueType)
		if err != nil {
			return false, fmt.Errorf("error parsing value: %v", err)
		}
	}

	return valueFromResponse, nil
}

// getPreviousResponseValue retrieves the previous response value for difference calculation
func (k Keeper) getPreviousResponseValue(ctx sdk.Context, flowID uint64, responseIndex uint32, responseKey string, valueType string, isICQ bool) (interface{}, error) {
	// Get the flow to access previous responses
	flowHistory, err := k.GetFlowHistory(ctx, flowID)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %d", flowID)
	}

	// Need at least 1 execution to have a previous value
	if len(flowHistory) == 0 {
		return nil, fmt.Errorf("no execution history available for flow %d", flowID)
	}

	k.Logger(ctx).Debug("getPreviousResponseValue", "history_length", len(flowHistory))

	// Look for the most recent successful execution
	for i := len(flowHistory) - 1; i >= 0; i-- {
		history := flowHistory[i]

		// Handle ICQ query responses
		if isICQ && len(history.QueryResponses) > 0 {
			if int(responseIndex) < len(history.QueryResponses) {
				k.Logger(ctx).Debug("Found ICQ query response", "index", i, "response", history.QueryResponses[responseIndex])
				return k.getValueFromResponse(
					ctx,
					responseIndex,
					responseKey,
					&valueType,
					[]byte(history.QueryResponses[responseIndex]),
					nil,
				)
			}
			continue
		}
	}
	for i := len(flowHistory) - 2; i >= 0; i-- {
		history := flowHistory[i]

		// Skip failed or empty executions
		if !history.Executed || len(history.Errors) > 0 {
			continue
		}
		// Handle regular message responses
		if len(history.MsgResponses) > 0 && int(responseIndex) < len(history.MsgResponses) {
			k.Logger(ctx).Debug("Found message response", "index", i)
			return k.getValueFromResponse(
				ctx,
				responseIndex,
				responseKey,
				&valueType,
				nil,
				history.MsgResponses,
			)
		}
	}

	return nil, fmt.Errorf("no successful execution found in history")
}

// parseResponseValue retrieves and parses the value of a response key to the specified response type
func parseResponseValue(response interface{}, responseKey, responseType string) (interface{}, error) {

	val, err := traverseFields(response, responseKey)
	if err != nil {
		return nil, err
	}

	switch responseType {
	case "string":
		if val.Kind() == reflect.String {
			return val.String(), nil
		}
	case "sdk.Coin":
		if val.Kind() == reflect.Slice && val.Type().Elem().Name() == "Coin" {
			coins := val.Interface().(sdk.Coins)
			return coins[0], nil
		}
		if val.Kind() == reflect.Struct {
			amountField := val.FieldByName("Amount")
			denomField := val.FieldByName("Denom")
			if amountField.IsValid() && denomField.IsValid() && amountField.Type() == reflect.TypeOf(math.Int{}) && denomField.Kind() == reflect.String {
				amount := amountField.Interface().(math.Int)
				return sdk.Coin{
					Amount: amount,
					Denom:  denomField.String(),
				}, nil
			}
		}
	case "sdk.Coins":
		if val.Kind() == reflect.Slice && val.Type().Elem().Name() == "Coin" {
			coins := val.Interface().(sdk.Coins)
			return coins, nil
		}
	case "math.Int":
		if val.Kind() == reflect.Struct {

			if val.Type().Name() == "Int" {
				return val.Interface().(math.Int), nil
			}

		}
		if val.Kind() == reflect.String {
			intVal, ok := math.NewIntFromString(val.String())
			if !ok {
				return nil, fmt.Errorf("failed to parse math.Int from string: %v", val.String())
			}
			return intVal, nil
		}
		if val.Kind() == reflect.Float64 {
			intVal := math.NewInt(int64(val.Float()))
			return intVal, nil
		}
		if val.Kind() == reflect.Int || val.Kind() == reflect.Int64 {
			intVal := math.NewInt(reflect.ValueOf(val.Interface()).Int())
			return intVal, nil
		} else {
			fmt.Println("parsing math.Int from response failed", val)
		}
	case "math.Dec":
		if val.Kind() == reflect.String {
			decVal, err := math.NewDecFromString(val.Interface().(string))
			if err != nil {
				return nil, err
			}
			fmt.Print("parsed math.Dec from string", decVal)
			return decVal, nil
		}
	case "[]string":
		if val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.String {
			return val.Interface().([]string), nil
		}
	case "[]math.Int":
		if val.Kind() == reflect.Slice && val.Type().Elem().Name() == "Int" {
			return val.Interface().([]math.Int), nil
		}

	}

	return nil, fmt.Errorf("field could not be parsed as %s in %v ", responseType, val)
}

// parseOperand parses the operand based on the response type
func parseOperand(operand string, responseType string) (interface{}, error) {
	switch responseType {
	case "string":
		return operand, nil
	case "sdk.Coin":
		coin, err := sdk.ParseCoinNormalized(operand)
		return coin, err
	case "sdk.Coins":
		coins, err := sdk.ParseCoinsNormalized(operand)
		return coins, err
	case "math.Int":
		var sdkInt math.Int
		sdkInt, ok := math.NewIntFromString(operand)
		if !ok {
			return nil, fmt.Errorf("unsupported int operand")
		}
		return sdkInt, nil
	case "math.Dec":
		dec, err := math.NewDecFromString(operand)
		fmt.Print("parsed math.Dec from string", dec)
		if err != nil {
			return nil, err
		}
		return dec, nil

	case "[]string":
		var strArr []string
		err := json.Unmarshal([]byte(operand), &strArr)
		return strArr, err
	case "[]math.Int":
		var intArr []math.Int
		err := json.Unmarshal([]byte(operand), &intArr)
		return intArr, err
	case "int":
		var intVal int64
		_, err := fmt.Sscanf(operand, "%d", &intVal)
		if err != nil {
			return nil, fmt.Errorf("failed to parse int operand: %w", err)
		}
		return intVal, nil
		// Add more cases for other response types as needed
		// For example, if you have a custom proto type, you can handle it here
		// case "customProtoType":
		// 	var customProto CustomProtoType
		// 	err := proto.Unmarshal([]byte(operand), &customProto)
		// 	return customProto, err
	}
	return nil, fmt.Errorf("unsupported operand type: %s", responseType)
}

// parseICQResponse parses the ICQ response
func parseICQResponse(response []byte, valueType *string) (interface{}, error) {
	switch *valueType {
	case "string":
		return string(response), nil
	case "sdk.Coin":
		var coin sdk.Coin
		err := coin.Unmarshal(response)
		return coin, err
	//not tested
	case "sdk.Coins":
		var coins sdk.Coins
		for len(response) > 0 {
			var coin sdk.Coin
			if err := coin.Unmarshal(response); err != nil {
				return nil, err
			}
			coins = append(coins, coin)

			// Move the response slice forward to the next coin
			remaining, err := remainingBytesAfterCoinUnmarshal(coin, response)
			if err != nil {
				return nil, err
			}
			response = remaining
		}
		return coins, nil
	case "math.Int":
		var sdkInt math.Int
		err := sdkInt.Unmarshal(response)
		return sdkInt, err
	case "math.Dec":
		var decStr string
		err := json.Unmarshal(response, &decStr)
		if err != nil {
			return nil, err
		}
		dec, err := math.NewDecFromString(decStr)
		if err != nil {
			return nil, err
		}

		return dec, nil
	case "[]string":
		var strings []string
		err := json.Unmarshal(response, &strings)
		return strings, err
	case "[]math.Int":
		var ints []math.Int
		for len(response) > 0 {
			var intVal math.Int
			if err := intVal.Unmarshal(response); err != nil {
				return nil, err
			}
			ints = append(ints, intVal)

			remaining, err := remainingBytesAfterIntUnmarshal(intVal, response)
			if err != nil {
				return nil, err
			}
			response = remaining
		}
		return ints, nil
	case "json":
		var jsonObj map[string]interface{}
		if err := json.Unmarshal(response, &jsonObj); err != nil {
			return nil, err
		}
		return jsonObj, nil
	}

	if strings.Contains(*valueType, "osmosistwapv1beta1.TwapRecord") {
		var twapRecord osmosistwapv1beta1.TwapRecord
		if err := twapRecord.Unmarshal(response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal osmosistwapv1beta1.TwapRecord: %w", err)
		}

		// Check if we have a specific field requested (format: "osmosistwapv1beta1.TwapRecord.FieldName")
		parts := strings.Split(*valueType, ".")
		if len(parts) > 2 {
			fieldName := parts[2]
			*valueType = "math.Dec"
			rv := reflect.ValueOf(twapRecord)
			field := rv.FieldByName(fieldName)
			if !field.IsValid() {
				return nil, fmt.Errorf("invalid field name: %s", fieldName)
			}

			legacyDec, ok := field.Interface().(math.LegacyDec)
			if !ok {
				return nil, fmt.Errorf("field %s is not a LegacyDec", fieldName)
			}

			dec, err := math.NewDecFromString(legacyDec.String())
			if err != nil {
				return nil, fmt.Errorf("failed to parse LegacyDec: %w", err)
			}

			return dec, nil
		}
		dec, err := math.NewDecFromString(twapRecord.P0LastSpotPrice.String())
		if err != nil {
			return nil, fmt.Errorf("failed to parse P0LastSpotPrice: %w", err)
		}

		*valueType = "math.Dec"
		return dec, nil

	}

	return nil, fmt.Errorf("unsupported operand type: %s", *valueType)
}

// contains checks if a value contains an operand
func contains(value, operand interface{}) bool {
	val := reflect.ValueOf(value)
	operandVal := reflect.ValueOf(operand)

	switch val.Kind() {
	case reflect.String:
		return strings.Contains(val.String(), operandVal.String())
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), operandVal.Interface()) {
				return true
			}
		}
	}
	// Custom case for sdk.Coins
	if coins, ok := value.(sdk.Coins); ok {
		if coin, ok := operand.(sdk.Coins); ok {
			_, notOk := coins.SafeSub(coin...)
			return !notOk
		}
	}
	return false
}

func compareNumbers(value, operand interface{}, compareFunc func(int64, int64) bool) (bool, error) {
	val := reflect.ValueOf(value)
	operandVal := reflect.ValueOf(operand)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareFunc(val.Int(), operandVal.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Convert unsigned integers to int64 safely
		valInt64 := int64(val.Uint())
		operandInt64 := int64(operandVal.Uint())
		return compareFunc(valInt64, operandInt64), nil
	case reflect.Struct:
		if val.Type().Name() == "Int" {
			// Use the `Int64` method for math.Int
			return compareFunc(
				val.MethodByName("Int64").Call(nil)[0].Int(),
				operandVal.MethodByName("Int64").Call(nil)[0].Int(),
			), nil
		}
		if val.Type() == reflect.TypeOf(sdk.Coin{}) {
			valDenom := val.FieldByName("Denom").String()
			operandDenom := operandVal.FieldByName("Denom").String()

			// Ensure the Denom fields are equal
			if valDenom != operandDenom {
				return false, fmt.Errorf("denom mismatch: %s vs %s", valDenom, operandDenom)
			}

			valAmount := val.FieldByName("Amount")
			operandAmount := operandVal.FieldByName("Amount")
			if valAmount.IsValid() && operandAmount.IsValid() {
				// Ensure Amount field has Int64 method
				valAmountMethod := valAmount.MethodByName("Int64")
				operandAmountMethod := operandAmount.MethodByName("Int64")

				if valAmountMethod.IsValid() && operandAmountMethod.IsValid() {
					valInt := valAmountMethod.Call(nil)[0].Int()
					operandInt := operandAmountMethod.Call(nil)[0].Int()
					return compareFunc(valInt, operandInt), nil
				}
				return false, fmt.Errorf("amount field missing Int64 method")
			}
			return false, fmt.Errorf("sdk.Coin struct missing Amount field")
		}
	}
	return false, fmt.Errorf("unsupported numeric type: %v", val.Kind())
}

// TraverseFields traverses the nested fields of a struct or slice based on the provided keys
func traverseFields(msgInterface interface{}, inputKey string) (reflect.Value, error) {
	keys := strings.Split(inputKey, ".")
	val := reflect.ValueOf(msgInterface)

	// If input is a primitive and key is empty or single, return directly
	if len(keys) == 1 && (val.Kind() != reflect.Struct && val.Kind() != reflect.Slice && val.Kind() != reflect.Ptr) {
		return val, nil
	}

	if inputKey != "" {
		for _, key := range keys {
			// Handle slices
			if val.Kind() == reflect.Slice {
				index, err := parseIndex(key, val.Len())
				if err != nil {
					return reflect.Value{}, err
				}
				if index >= val.Len() {
					return reflect.Value{}, fmt.Errorf("index %d out of bounds for slice", index)
				}
				val = val.Index(index)
			} else {
				// If the value is a pointer, get the element it points to
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}

				// Ensure we're dealing with a struct
				if val.Kind() != reflect.Struct {
					return reflect.Value{}, fmt.Errorf("expected a struct, got %v", val.Kind())
				}

				val = val.FieldByName(key)
				if !val.IsValid() {
					return reflect.Value{}, fmt.Errorf("field %s not found in interface %v", key, msgInterface)
				}
			}
		}
	}
	return val, nil
}

// parseIndex parses a string index for slices, supporting negative indices.
// 'length' is the length of the slice being indexed.
func parseIndex(key string, length int) (int, error) {
	var index int
	_, err := fmt.Sscanf(key, "[%d]", &index)
	if err != nil {
		return -1, fmt.Errorf("invalid or missing slice index in key: %s", key)
	}

	// Convert negative index to a positive one.
	if index < 0 {
		index = length + index
	}

	// Validate the index to ensure it's within the slice bounds.
	if index < 0 || index >= length {
		return -1, fmt.Errorf("index %d out of bounds for slice of length %d", index, length)
	}

	return index, nil
}

// remainingBytesAfterCoinUnmarshal calculates the remaining bytes after an sdk.Coin is unmarshalled
func remainingBytesAfterCoinUnmarshal(coin sdk.Coin, response []byte) ([]byte, error) {
	// This would typically involve calculating the size of the coin in bytes,
	// or you could also re-encode the coin to know its byte size for skipping.
	// Placeholder for your actual implementation:
	encodedCoin, err := coin.Marshal() // Encoding `coin` again to get its length
	if err != nil {
		return nil, err
	}
	if len(response) < len(encodedCoin) {
		return nil, fmt.Errorf("remaining bytes insufficient")
	}
	return response[len(encodedCoin):], nil
}

// remainingBytesAfterIntUnmarshal calculates the remaining bytes after an math.Int is unmarshalled
func remainingBytesAfterIntUnmarshal(intVal math.Int, response []byte) ([]byte, error) {
	encodedInt, err := intVal.Marshal()
	if err != nil {
		return nil, err
	}
	if len(response) < len(encodedInt) {
		return nil, fmt.Errorf("remaining bytes insufficient")
	}
	return response[len(encodedInt):], nil
}

// extractJSONField extracts a field from a JSON object using dot notation (e.g. "foo.bar.baz")
func extractJSONField(jsonObj map[string]interface{}, key string) (interface{}, error) {
	keys := strings.Split(key, ".")
	var val interface{} = jsonObj
	for _, k := range keys {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("intermediate value for key '%s' is not a JSON object", k)
		}
		val, ok = m[k]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found in JSON object", k)
		}
	}
	return val, nil
}

// normalizeJSONTypes attempts to convert float64 to int64 or string if the other operand is int/string, for JSON extracted values
func normalizeJSONTypes(a, b interface{}) (interface{}, interface{}) {
	// If both are nil or same type, nothing to do
	if a == nil || b == nil || reflect.TypeOf(a) == reflect.TypeOf(b) {
		return a, b
	}

	// Try to parse numeric strings to numbers for JSON comparisons
	aNum, aIsNum := tryParseNumeric(a)
	bNum, bIsNum := tryParseNumeric(b)
	if aIsNum && bIsNum {
		// If both are numbers, use int64 if both are ints, else float64
		if aInt, okA := aNum.(int64); okA {
			if bInt, okB := bNum.(int64); okB {
				return aInt, bInt
			}
		}
		// fallback to float64
		aF := toFloat64(aNum)
		bF := toFloat64(bNum)
		return aF, bF
	}
	// If only one is numeric, keep as is (will likely fail comparison, but that's correct)

	// float64 from JSON, int64 operand
	if af, ok := a.(float64); ok {
		switch bv := b.(type) {
		case int:
			return int(af), bv
		case int64:
			return int64(af), bv
		case string:
			return fmt.Sprintf("%v", af), bv
		}
	}
	if bf, ok := b.(float64); ok {
		switch av := a.(type) {
		case int:
			return av, int(bf)
		case int64:
			return av, int64(bf)
		case string:
			return av, fmt.Sprintf("%v", bf)
		}
	}
	// string/int cross-compare
	if as, ok := a.(string); ok {
		switch bv := b.(type) {
		case int:
			return as, fmt.Sprintf("%d", bv)
		case int64:
			return as, fmt.Sprintf("%d", bv)
		}
	}
	if bs, ok := b.(string); ok {
		switch av := a.(type) {
		case int:
			return fmt.Sprintf("%d", av), bs
		case int64:
			return fmt.Sprintf("%d", av), bs
		}
	}
	return a, b
}

// tryParseNumeric attempts to parse an interface{} as int64 or float64 if it's a string or number
func tryParseNumeric(v interface{}) (interface{}, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		// If it's an integer value, return as int64
		if val == float64(int64(val)) {
			return int64(val), true
		}
		return val, true
	case string:
		// Try int64 first
		var i int64
		_, err := fmt.Sscanf(val, "%d", &i)
		if err == nil {
			return i, true
		}
		// Try float64
		var f float64
		_, err = fmt.Sscanf(val, "%f", &f)
		if err == nil {
			return f, true
		}
		return nil, false
	default:
		return nil, false
	}
}

// calculateDifference calculates the difference between current and previous values for difference mode
func (k Keeper) calculateDifference(ctx sdk.Context, currentValue, previousValue interface{}) (interface{}, error) {
	switch current := currentValue.(type) {
	case math.Dec:
		prevDec, ok := previousValue.(math.Dec)
		if !ok {
			return nil, fmt.Errorf("previous value is not math.Dec, previous value type: %T, current value type: %T", previousValue, currentValue)
		}
		// Calculate difference with sign preserved
		diff, err := current.Sub(prevDec)
		if err != nil {
			return nil, fmt.Errorf("math.Dec subtraction failed: %w", err)
		}
		k.Logger(ctx).Debug("calculated difference", "current", current, "previous", prevDec, "difference", diff)
		return diff, nil

	case math.Int:
		prevInt, ok := previousValue.(math.Int)
		if !ok {
			return nil, fmt.Errorf("previous value is not math.Int")
		}
		// Calculate difference with sign preserved
		diff, err := current.SafeSub(prevInt)
		if err != nil {
			return nil, fmt.Errorf("math.Int subtraction failed: %w", err)
		}
		k.Logger(ctx).Debug("calculated difference", "current", current, "previous", prevInt, "difference", diff)
		return diff, nil

	case sdk.Coin:
		prevCoin, ok := previousValue.(sdk.Coin)
		if !ok {
			prevCoins, ok := previousValue.(sdk.Coins)
			if !ok {
				return nil, fmt.Errorf("previous value is not sdk.Coin or sdk.Coins")
			}
			prevCoin = prevCoins[0]
		}

		diff, err := current.SafeSub(prevCoin)
		if err != nil {
			diff, err = prevCoin.SafeSub(current)
			if err != nil {
				return nil, fmt.Errorf("sdk.Coin subtraction failed: %w", err)
			}
		}
		k.Logger(ctx).Debug("calculated difference", "current", current, "previous", prevCoin, "difference", diff)
		return diff, nil

	default:
		return nil, fmt.Errorf("unsupported type %T", current)
	}
}

// toFloat64 converts int64 or float64 to float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int64:
		return float64(val)
	case float64:
		return val
	default:
		return 0
	}
}
