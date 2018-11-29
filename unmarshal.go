package qs

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"github.com/kr/pretty"
)

var nameRegex = pcre.MustCompile(`\A[\[\]]*([^\[\]]+)\]*`, 0)
var objectRegex1 = pcre.MustCompile(`^\[\]\[([^\[\]]+)\]$`, 0)
var objectRegex2 = pcre.MustCompile(`^\[\](.+)$`, 0)

var enableConvertArrays = false

func ConvertArrays(enable bool) {
	enableConvertArrays = enable
}

func Unmarshal(qs string) (interface{}, error) {
	components := strings.Split(qs, "&")
	params := map[string]interface{}{}

	for _, c := range components {
		tuple := strings.Split(c, "=")

		for i, item := range tuple {
			if unesc, err := url.QueryUnescape(item); err == nil {
				tuple[i] = unesc
			}
		}

		key := ""

		if len(tuple) > 0 {
			key = tuple[0]
		}

		value := interface{}(nil)

		if len(tuple) > 1 {
			value = tuple[1]
		}

		if err := normalizeParams(params, key, value); err != nil {
			return nil, err
		}
	}

	if !enableConvertArrays {
		return params, nil
	}

	return convertArrays(params), nil
}

func getPos(component string) int {
	pos := strings.Index(component, "]=")
	if pos == -1 {
		pos = strings.Index(component, "=")
	} else {
		pos++
	}
	return pos
}

func splitKeyValue(component string) (key string, value string, err error) {
	pos := getPos(component)
	if pos == -1 {
		key = component
		value = ""
	} else {
		key = component[0:pos]
		value = component[pos+1:]
	}

	key, err = url.QueryUnescape(key)
	if err != nil {
		return "", "", err
	}

	value, err = url.QueryUnescape(value)
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

func Parse(qs string) (interface{}, error) {
	components := strings.Split(qs, "&")

	params := map[string]interface{}{}

	for _, c := range components {

		key, value, err := splitKeyValue(c)
		if err != nil {
			return nil, err
		}

		if err := normalizeParams2(params, key, value); err != nil {
			return nil, err
		}
		pretty.Println("params:", params)
	}

	if !enableConvertArrays {
		return params, nil
	}

	pretty.Println(params)

	return convertArrays(params), nil
}

func normalizeParams2(params map[string]interface{}, key string, value interface{}) error {

	fmt.Println(">>>>>>>----------------------------")
	pretty.Println("params:", params)
	pretty.Println("key:", key)
	pretty.Println("value:", value)
	fmt.Println("<<<<<<<----------------------------")

	nameMatcher := nameRegex.MatcherString(key, 0)
	k := nameMatcher.GroupString(1)
	after := ""

	if pos := nameRegex.FindIndex([]byte(key), 0); len(pos) == 2 {
		after = key[pos[1]:]
	}

	objectMatcher1 := objectRegex1.MatcherString(after, 0)
	objectMatcher2 := objectRegex2.MatcherString(after, 0)

	if k == "" {
		fmt.Println("porra nenhuma!")

		params = map[string]interface{}{}
		fmt.Println(params)
		return nil

	} else if after == "" {

		ival, ok := params[k]
		if !ok {
			fmt.Println("insert simples, key nao existe ainda")
			params[k] = value
			return nil
		}

		// somos arrayzes!

		fmt.Println("ggg.................")
		switch i := ival.(type) {
		case []interface{}:
			params[k] = append(i, value)
		case string:
			params[k] = []interface{}{i, value}
		case map[string]interface{}:
			pretty.Println("houston, we have a problem")

			array, ok := toArray(i)
			if !ok {
				return errors.New("houston, we have a problem")
			}
			array = append(array, value)
			params[k] = array

		default:
			fmt.Printf("\n1 - panic\n\n")
		}

		return nil

	} else if after == "[]" {

		ival, ok := params[k]
		if !ok {
			fmt.Println("insert simples, key nao existe ainda, array")
			params[k] = []interface{}{value}
			return nil
		}

		fmt.Println("fff.................")
		switch i := ival.(type) {
		case []interface{}:
			params[k] = append(i, value)
		case string:
			params[k] = []interface{}{i, value}
		case map[string]interface{}:
			pretty.Println("houston, temos outro problema")
		default:
			fmt.Printf("\n1 - panic\n\n")
		}

		return nil

	} else if objectMatcher1.Matches() || objectMatcher2.Matches() {

		fmt.Println("viiiixii")

		childKey := ""

		if objectMatcher1.Matches() {
			childKey = objectMatcher1.GroupString(1)
		} else if objectMatcher2.Matches() {
			childKey = objectMatcher2.GroupString(1)
		}

		if childKey != "" {
			ival, ok := params[k]

			if !ok {
				params[k] = []interface{}{}
				ival = params[k]
			}

			array, ok := ival.([]interface{})

			if !ok {
				return fmt.Errorf("Expected type '[]interface{}' for key '%s', but got '%T'", k, ival)
			}

			if length := len(array); length > 0 {
				if hash, ok := array[length-1].(map[string]interface{}); ok {
					if _, ok := hash[childKey]; !ok {
						normalizeParams(hash, childKey, value)
						return nil
					}
				}
			}

			newHash := map[string]interface{}{}
			normalizeParams(newHash, childKey, value)
			params[k] = append(array, newHash)

			return nil
		}
	}

	fmt.Println("aqui")

	ival, ok := params[k]
	if !ok {
		params[k] = map[string]interface{}{}
		ival = params[k]
	}

	switch i := ival.(type) {
	case map[string]interface{}:
		return normalizeParams(i, after, value)
	case string:
		// viramos array
		params[k] = []interface{}{i, value}
		return nil
	default:
		return fmt.Errorf("Expected type 'map[string]interface{}' for key '%s', but got '%T'", k, ival)
	}
	/*
		hash, ok := ival.(map[string]interface{})

		if !ok {
			return fmt.Errorf("Expected type 'map[string]interface{}' for key '%s', but got '%T'", k, ival)
		}

		if err := normalizeParams(hash, after, value); err != nil {
			return err
		}

		return nil
	*/
}

func indexArray(in map[string]interface{}) ([]int, bool) {

	if len(in) == 0 {
		return nil, false
	}

	arr := []int{}
	for key := range in {
		i, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			return nil, false
		}
		arr = append(arr, int(i))
	}
	sort.Ints(arr)
	return arr, true
}

func toArray(in map[string]interface{}) ([]interface{}, bool) {

	indexArr, ok := indexArray(in)
	if !ok {
		return nil, false
	}

	pretty.Println("indexarray:", indexArr)

	arr := []interface{}{}
	for _, index := range indexArr {
		key := strconv.Itoa(index)
		arr = append(arr, in[key])
	}

	return arr, true
}

func convertArrays(in map[string]interface{}) interface{} {

	arr, ok := toArray(in)
	if ok { // I am a array
		pretty.Println("array!", arr)
		for i, value := range arr {
			switch v := value.(type) {
			case map[string]interface{}:
				arr[i] = convertArrays(v)
			}
		}
		return arr
	}

	// I am a map
	for key, value := range in {
		switch v := value.(type) {
		case map[string]interface{}:
			in[key] = convertArrays(v)
		}
	}
	return in
}

func normalizeParams(params map[string]interface{}, key string, value interface{}) error {
	nameMatcher := nameRegex.MatcherString(key, 0)
	k := nameMatcher.GroupString(1)
	after := ""

	if pos := nameRegex.FindIndex([]byte(key), 0); len(pos) == 2 {
		after = key[pos[1]:]
	}

	objectMatcher1 := objectRegex1.MatcherString(after, 0)
	objectMatcher2 := objectRegex2.MatcherString(after, 0)

	if k == "" {
		return nil

	} else if after == "" {
		params[k] = value
		return nil

	} else if after == "[]" {
		ival, ok := params[k]

		if !ok {
			params[k] = []interface{}{value}
			return nil
		}

		array, ok := ival.([]interface{})

		if !ok {
			return fmt.Errorf("Expected type '[]interface{}' for key '%s', but got '%T'", k, ival)
		}

		params[k] = append(array, value)
		return nil

	} else if objectMatcher1.Matches() || objectMatcher2.Matches() {

		childKey := ""

		if objectMatcher1.Matches() {
			childKey = objectMatcher1.GroupString(1)
		} else if objectMatcher2.Matches() {
			childKey = objectMatcher2.GroupString(1)
		}

		if childKey != "" {
			ival, ok := params[k]

			if !ok {
				params[k] = []interface{}{}
				ival = params[k]
			}

			array, ok := ival.([]interface{})

			if !ok {
				return fmt.Errorf("Expected type '[]interface{}' for key '%s', but got '%T'", k, ival)
			}

			if length := len(array); length > 0 {
				if hash, ok := array[length-1].(map[string]interface{}); ok {
					if _, ok := hash[childKey]; !ok {
						normalizeParams(hash, childKey, value)
						return nil
					}
				}
			}

			newHash := map[string]interface{}{}
			normalizeParams(newHash, childKey, value)
			params[k] = append(array, newHash)

			return nil
		}
	}

	ival, ok := params[k]

	if !ok {
		params[k] = map[string]interface{}{}
		ival = params[k]
	}

	hash, ok := ival.(map[string]interface{})

	if !ok {
		return fmt.Errorf("Expected type 'map[string]interface{}' for key '%s', but got '%T'", k, ival)
	}

	if err := normalizeParams(hash, after, value); err != nil {
		return err
	}

	return nil
}
