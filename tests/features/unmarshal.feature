Feature: Unmarshal using the `env` tag
    Scenario: All primitive types should be unmarshalled with zero values when no environment variables are set
        Given a "base" struct test struct is instantiated
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {}
        """
        And there should be no errors

    Scenario: String types should be unmarshalled correctly when valid values are provided
        Given a "base" struct test struct is instantiated
        And the environment variable "TEST_ENV_STR" is set to "passed"
        And the environment variable "TEST_ENV_STR_PTR" is set to "passed_ptr"
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {
            "String": "passed",
            "StringPtr": "passed_ptr"
        }
        """
        And there should be no errors


    Scenario: All integer types should be unmarshalled correctly when valid values are provided

        Given a "base" struct test struct is instantiated
        And the environment variable "TEST_ENV_INT_PTR" is set to "1"
        And the environment variable "TEST_ENV_INT_MAX" is set to "2147483647"
        And the environment variable "TEST_ENV_INT8_MAX" is set to "127"
        And the environment variable "TEST_ENV_INT16_MAX" is set to "32767"
        And the environment variable "TEST_ENV_INT32_MAX" is set to "2147483647"
        And the environment variable "TEST_ENV_INT64_MAX" is set to "9223372036854775807"
        And the environment variable "TEST_ENV_INT_MIN" is set to "-2147483648"
        And the environment variable "TEST_ENV_INT8_MIN" is set to "-128"
        And the environment variable "TEST_ENV_INT16_MIN" is set to "-32768"
        And the environment variable "TEST_ENV_INT32_MIN" is set to "-2147483648"
        And the environment variable "TEST_ENV_INT64_MIN" is set to "-9223372036854775808"
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {

            "IntPtr": 1,
            "IntMax": 2147483647,
            "Int8Max": 127,
            "Int16Max": 32767,
            "Int32Max": 2147483647,
            "Int64Max": 9223372036854775807,
            "IntMin": -2147483648,
            "Int8Min": -128,
            "Int16Min": -32768,
            "Int32Min": -2147483648,
            "Int64Min": -9223372036854775808
        }
        """
        And there should be no errors
    Scenario: All unexported types should be ignored 
        Given a "base" struct test struct is instantiated
        And the environment variable "TEST_ENV_UNEXPORTED" is set to "ingores_unexported_fields"
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {
            "Unexported": ""
        }
        """
        And there should be no errors


    Scenario: All Unsigned Integer types should be unmarshalled correctly when valid values are provided
        Given a "base" struct test struct is instantiated
        And the environment variable "TEST_ENV_UINT_PTR" is set to "1"
        And the environment variable "TEST_ENV_UINT_MAX" is set to "4294967295"
        And the environment variable "TEST_ENV_UINT8_MAX" is set to "255"
        And the environment variable "TEST_ENV_UINT16_MAX" is set to "65535"
        And the environment variable "TEST_ENV_UINT32_MAX" is set to "4294967295"
        And the environment variable "TEST_ENV_UINT64_MAX" is set to "18446744073709551615"
        And the environment variable "TEST_ENV_UINT_MIN" is set to "0"
        And the environment variable "TEST_ENV_UINT8_MIN" is set to "0"
        And the environment variable "TEST_ENV_UINT16_MIN" is set to "0"
        And the environment variable "TEST_ENV_UINT32_MIN" is set to "0"
        And the environment variable "TEST_ENV_UINT64_MIN" is set to "0"
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {
            "UintPtr": 1,
            "UintMax": 4294967295,
            "Uint8Max": 255,
            "Uint16Max": 65535,
            "Uint32Max": 4294967295,
            "Uint64Max": 18446744073709551615,
            "UintMin": 0,
            "Uint8Min": 0,
            "Uint16Min": 0,
            "Uint32Min": 0,
            "Uint64Min": 0
        }
        """
        And there should be no errors


    Scenario: All Float types should be unmarshalled correctly when valid values are provided
        Given a "base" struct test struct is instantiated
        And the environment variable "TEST_ENV_FLOAT32_PTR" is set to "1.0"
        And the environment variable "TEST_ENV_FLOAT64_PTR" is set to "1.0"
        And the environment variable "TEST_ENV_FLOAT32_MAX" is set to "3.40282346638528859811704183484516925440e+38"
        And the environment variable "TEST_ENV_FLOAT64_MAX" is set to "1.79769313486231570814527423731704356798070e+308"
        And the environment variable "TEST_ENV_FLOAT32_MIN" is set to "1.401298464324817070923729583289916131280e-45"
        And the environment variable "TEST_ENV_FLOAT64_MIN" is set to "4.9406564584124654417656879286822137236505980e-324"
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {
            "Float32Ptr": 0,
            "Float64Ptr": 0,
            "Float32Max": 3.40282346638528859811704183484516925440e+38,
            "Float64Max": 1.79769313486231570814527423731704356798070e+308,
            "Float32Min": 1.401298464324817070923729583289916131280e-45,
            "Float64Min": 4.9406564584124654417656879286822137236505980e-324
        }
        """
        And there should be no errors
    Scenario: All Boolean types should be unmarshalled correctly when valid values are provided     
        Given a "base" struct test struct is instantiated
        And the environment variable "TEST_ENV_BOOL_TRUE" is set to "true"
        And the environment variable "TEST_ENV_BOOL_FALSE" is set to "false"
        And the environment variable "TEST_ENV_BOOL_PTR_TRUE" is set to "true"
        And the environment variable "TEST_ENV_BOOL_PTR_FALSE" is set to "false"
        And the environment variable "TEST_ENV_BOOL_YES" is set to "yes"
        And the environment variable "TEST_ENV_BOOL_NO" is set to "no"
        And the environment variable "TEST_ENV_BOOL_ON" is set to "on"
        And the environment variable "TEST_ENV_BOOL_OFF" is set to "off"
        And the environment variable "TEST_ENV_BOOL_1" is set to "1"
        And the environment variable "TEST_ENV_BOOL_0" is set to "0"
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {
            "BoolTrue": true,
            "BoolFalse": false,
            "BoolPtrTrue": true,
            "BoolPtrFalse": false,
            "BoolYes": true,
            "BoolNo": false,
            "BoolOn": true,
            "BoolOff": false,
            "Bool1": true,
            "Bool0": false
        }
        """
        And there should be no errors

     Scenario: All nested types should be unmarshalled correctly when valid values are provided     
        Given a "base" struct test struct is instantiated
        And the environment variable "TEST_STRUCT_FIELD" is set to "passed"
        And the environment variable "TEST_STRUCT_FIELD_PTR" is set to "passed_ptr"
        When the struct is passed as an argument to Unmarshal
        Then the struct should have the following values:
        """
        {
            "NestedStruct": {
                "Field": "passed"
            },
            "NestedStructPointer": {
                "Field": "passed_ptr"
            }
         
        }
        """
        And there should be no errors