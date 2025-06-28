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
	var valueFromResponse interface{}
	var err error

	if queryCallback != nil {
		valueFromResponse, err = parseICQResponse(queryCallback, comparison.ValueType)
		if err != nil {
			var respAny cdctypes.Any
			err := k.cdc.Unmarshal(queryCallback, &respAny)
			if err != nil {
				k.Logger(ctx).Debug("response comparison: could not parse query callback data", "CallbackData", queryCallback)
				return false, fmt.Errorf("response comparison: error parsing query callback data: %v", queryCallback)
			}
			if respAny.Value != nil {

				var resp interface{}
				err = k.cdc.UnpackAny(&respAny, &resp)
				if err != nil {
					return false, fmt.Errorf("response comparison: error unpacking: %v", err)
				}
				valueFromResponse, err = parseResponseValue(queryCallback, comparison.ResponseKey, comparison.ValueType)
				if err != nil {
					return false, err
				}
			}

		}
	} else {
		var respAny cdctypes.Any
		if len(responses) == 0 {
			k.Logger(ctx).Debug("response comparison: could not parse response data")
		} else {
			respAny = *responses[comparison.ResponseIndex]
		}
		protoMsg, err := k.interfaceRegistry.Resolve(respAny.TypeUrl)
		if err != nil {
			return false, fmt.Errorf("failed to resolve type URL %s: %w", respAny.TypeUrl, err)
		}

		err = proto.Unmarshal(respAny.Value, protoMsg)
		if err != nil {
			return false, err
		}
		valueFromResponse, err = parseResponseValue(protoMsg, comparison.ResponseKey, comparison.ValueType)
		if err != nil {
			return false, fmt.Errorf("error parsing value: %v", err)
		}
	}

	operand, err := parseOperand(comparison.Operand, comparison.ValueType)
	if err != nil {
		return false, fmt.Errorf("error parsing operand: %v", err)
	}
	switch comparison.Operator {
	case types.ComparisonOperator_EQUAL:
		return reflect.DeepEqual(valueFromResponse, operand), nil
	case types.ComparisonOperator_NOT_EQUAL:
		return !reflect.DeepEqual(valueFromResponse, operand), nil
	case types.ComparisonOperator_CONTAINS:
		return contains(valueFromResponse, operand), nil
	case types.ComparisonOperator_NOT_CONTAINS:
		return !contains(valueFromResponse, operand), nil
	case types.ComparisonOperator_SMALLER_THAN:
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a < b })
	case types.ComparisonOperator_LARGER_THAN:
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a > b })
	case types.ComparisonOperator_GREATER_EQUAL:
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a >= b })
	case types.ComparisonOperator_LESS_EQUAL:
		return compareNumbers(valueFromResponse, operand, func(a, b int64) bool { return a <= b })
	case types.ComparisonOperator_STARTS_WITH:
		return strings.HasPrefix(valueFromResponse.(string), operand.(string)), nil
	case types.ComparisonOperator_ENDS_WITH:
		return strings.HasSuffix(valueFromResponse.(string), operand.(string)), nil
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

		var valueFromResponse interface{}
		var err error
		if queryCallback != nil {
			valueFromResponse, err = parseICQResponse(queryCallback, feedbackLoop.ValueType)
			if err != nil {
				var respAny cdctypes.Any
				err := k.cdc.Unmarshal(queryCallback, &respAny)
				if err != nil {
					k.Logger(ctx).Debug("use value: could not parse query callback data", "CallbackData", queryCallback)
					return fmt.Errorf("use value: error parsing query callback data: %v", queryCallback)
				}
				if respAny.Value != nil {

					var resp interface{}
					err = k.cdc.UnpackAny(&respAny, &resp)
					if err != nil {
						return fmt.Errorf("use value: error unpacking: %v", err)
					}
					valueFromResponse, err = parseResponseValue(queryCallback, feedbackLoop.ResponseKey, feedbackLoop.ValueType)
					if err != nil {
						return err
					}
				}

			}

		} else {
			history, err := k.GetFlowHistory(ctx, flowID)
			if err != nil {
				return err
			}
			if len(history) == 0 {
				return nil
			}
			responsesAnys := history[len(history)-1].MsgResponses
			if len(responsesAnys) == 0 {
				return nil
			}
			if int(feedbackLoop.ResponseIndex) >= len(responsesAnys) || int(feedbackLoop.ResponseIndex) < 0 {
				continue
			}

			protoMsg, err := k.interfaceRegistry.Resolve(responsesAnys[feedbackLoop.ResponseIndex].TypeUrl)
			if err != nil {
				return fmt.Errorf("failed to resolve type URL %s: %w", responsesAnys[feedbackLoop.ResponseIndex].TypeUrl, err)
			}

			err = proto.Unmarshal(responsesAnys[feedbackLoop.ResponseIndex].Value, protoMsg)
			if err != nil {
				return err
			}
			valueFromResponse, err = parseResponseValue(protoMsg, feedbackLoop.ResponseKey, feedbackLoop.ValueType)
			if err != nil {
				return err
			}
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

// parseResponseValue retrieves and parses the value of a response key to the specified response type
func parseResponseValue(response interface{}, responseKey, responseType string) (interface{}, error) {

	val, err := traverseFields(response, responseKey)
	if err != nil {
		// if responseKey == ""{
		// 	val =
		// }
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
	case "sdk.Int":
		if val.Kind() == reflect.Struct && val.Type().Name() == "Int" {
			return val.Interface().(math.Int), nil
		}
	case "[]string":
		if val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.String {
			return val.Interface().([]string), nil
		}
	case "[]sdk.Int":
		if val.Kind() == reflect.Slice && val.Type().Elem().Name() == "Int" {
			return val.Interface().([]math.Int), nil
		}
		// case "[]sdk.Coin":
		// 	if val.Kind() == reflect.Slice && val.Type().Elem().Name() == "Coin" {
		// 		return val.Interface().([]sdk.Coin), nil
		// 	}
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
	case "sdk.Int":
		var sdkInt math.Int
		sdkInt, ok := math.NewIntFromString(operand)
		if !ok {
			return nil, fmt.Errorf("unsupported int operand")
		}
		return sdkInt, nil
	case "[]string":
		var strArr []string
		err := json.Unmarshal([]byte(operand), &strArr)
		return strArr, err
	case "[]sdk.Int":
		var intArr []math.Int
		err := json.Unmarshal([]byte(operand), &intArr)
		return intArr, err
		// case "[]sdk.Coin":
		// 	coinArr, err := sdk.ParseCoinsNormalized(operand)
		// 	return coinArr, err
	}
	return nil, fmt.Errorf("unsupported operand type: %s", responseType)
}

// parseICQResponse parses the ICQ response
func parseICQResponse(response []byte, valueType string) (interface{}, error) {
	switch valueType {
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

			// Move the response slice forward by the length of the unmarshaled coin.
			// The length of each unmarshalled `coin` varies, so adjust accordingly.
			remaining, err := remainingBytesAfterCoinUnmarshal(coin, response)
			if err != nil {
				return nil, err
			}
			response = remaining
		}
		return coins, nil
	case "sdk.Int":
		var sdkInt math.Int
		err := sdkInt.Unmarshal(response)
		return sdkInt, err
	case "[]string":
		var strings []string
		err := json.Unmarshal(response, &strings)
		return strings, err
	case "[]sdk.Int":
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
	}
	//idea: add more protos here
	return nil, fmt.Errorf("unsupported operand type: %s", valueType)
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
			// Use the `Int64` method for sdk.Int
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
	if inputKey != "" {
		for _, key := range keys {

			// Handle slices
			if val.Kind() == reflect.Slice {
				index, err := parseIndex(key)
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

// parseIndex parses a string index for slices
func parseIndex(key string) (int, error) {
	var index int
	_, err := fmt.Sscanf(key, "[%d]", &index)
	if err != nil {
		return -1, fmt.Errorf("invalid or missing slice index in key: %s", key)
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
