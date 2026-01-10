package env

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/dubbikins/envy/v2/tag"
)
type Expander func(s string, mapping func(string) string) string
type Reader func(string) (string)

func ScanFunc(node *tag.Node, scanner *bufio.Scanner) (err error) {
			var token []byte
			var environment_variables = make([]string, 0,1)
			for scanner.Scan() {
				token = scanner.Bytes()
				if len(token)== 1 && token[0] == ';'{
					break 
				}else if len(token)== 1 && token[0] == '|'{
					continue
				} 
				environment_variables = append(environment_variables, string(bytes.TrimSpace(token)))
				
			}
			var key, value string
			for scanner.Scan() {
				if len(scanner.Bytes())== 1 && scanner.Bytes()[0] == ',' {
					return errors.Join(err, fmt.Errorf("env tag parsing error, unexpected ',' when scanning for key in field %s", node.Field().Name))
				}
				key = scanner.Text()
				if !scanner.Scan() || len(scanner.Bytes())== 1 && scanner.Bytes()[0] == ','  {
					return errors.Join(err, fmt.Errorf("env tag parsing error, expected '=' after key '%s' in field %s: got '%s'", key, node.Field().Name, scanner.Text()))
				}
				scanner.Scan() //consume the "="
				value = scanner.Text()
				if !scanner.Scan() {
					node.SetOption(key, value)
					break
				}else if len(scanner.Bytes())!= 1 || scanner.Bytes()[0] != ','{
					return errors.Join(err, fmt.Errorf("env tag parsing error, expected ',' after \"%s=%s\" in field %s: got \"%s\"", key, value, node.Field().Name, scanner.Text()))
				}
				node.SetOption(key, value)
			}
			if scanner.Err() != nil {
				return scanner.Err()
			}
			var env_reader Reader = os.Getenv
			var expand bool
			if _expand, found := node.Options()["expand"]; found && _expand == "true" {
			
					expand = true
			}
			// if _reader, found := node.Options()[".env"]; found && !strings.HasPrefix(_reader, "os") {
			// 	env_reader = func(s string) string {
			// 		var data, err = os.ReadFile(_reader)
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 		var vars = bytes.Split(data, []byte("\n"))
			// 		for _, kv := range vars {
			// 			if split := strings.IndexRune(string(kv), '='); split >= 0 {
			// 				if string(kv[:split]) == s {
			// 					return string(kv[split+1:])
			// 				}
			// 			}
			// 		}
			// 		return os.Getenv(s)
			// 	}
			// }

			for _, envVar := range environment_variables {
				
				var envVarValue string
				if expand {
					
					envVarValue = os.Expand(envVar, ExpandFn(node, env_reader))
				}else {
					envVarValue = env_reader(envVar)
				}
				if envVarValue != "" {
					_, err = node.Write([]byte(envVarValue))
					return
				}
				
				
				
			}
			if value := env_reader(node.Field().Name); value != "" {
				if _, err = node.Write([]byte(value)); err != nil {
						return
					}
			}
	return
}

func SplitFn (data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for advance < len(data){
		var char = data[advance]
		//If the first character is 
		if (char == ';' || char == '|' || char == '=' || char == ',') {
			//If it's the first character, then return the character
			if advance == 0  {
				advance = 1
				token = data[:advance]
				return
			}
			token = data[:advance]
			return
		}
		if char == '\'' && advance == 0  {
			data = data[1:]
			for advance < len(data) && data[advance] != '\''{
				advance+=1
			}
			token = data[:advance]
			advance += 2
			return
		} 
		advance+=1
	}

	if atEOF {
		advance = len(data)
		return 
	}

	token = data[:advance]
	return 
}


