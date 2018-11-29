package qs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPos(t *testing.T) {
	require.Equal(t, 1, getPos("0=foo"))
	require.Equal(t, 3, getPos("foo=c++"))
	require.Equal(t, 5, getPos("a[>=]=23"))
	require.Equal(t, 5, getPos("a[==]=23"))
	require.Equal(t, -1, getPos("foo"))
	require.Equal(t, 3, getPos("foo="))
	require.Equal(t, 3, getPos("foo=bar"))
	require.Equal(t, 5, getPos(" foo = bar = baz "))
	require.Equal(t, 3, getPos("foo=bar=baz"))
	require.Equal(t, 4, getPos("a[b]=c"))
	require.Equal(t, 7, getPos("a[b][c]=d"))
}

func TestSplitKeyValue(t *testing.T) {
	var key, value string
	var err error

	key, value, err = splitKeyValue("0=foo")
	require.NoError(t, err)
	require.Equal(t, "0", key)
	require.Equal(t, "foo", value)

	key, value, err = splitKeyValue("foo=c++")
	require.NoError(t, err)
	require.Equal(t, "foo", key)
	require.Equal(t, "c  ", value)

	key, value, err = splitKeyValue("a[>=]=23")
	require.NoError(t, err)
	require.Equal(t, "a[>=]", key)
	require.Equal(t, "23", value)

	key, value, err = splitKeyValue("a[==]=23")
	require.NoError(t, err)
	require.Equal(t, "a[==]", key)
	require.Equal(t, "23", value)

	key, value, err = splitKeyValue("foo")
	require.NoError(t, err)
	require.Equal(t, "foo", key)
	require.Equal(t, "", value)

	key, value, err = splitKeyValue("foo=")
	require.NoError(t, err)
	require.Equal(t, "foo", key)
	require.Equal(t, "", value)

	key, value, err = splitKeyValue("foo=bar")
	require.NoError(t, err)
	require.Equal(t, "foo", key)
	require.Equal(t, "bar", value)

	key, value, err = splitKeyValue(" foo = bar = baz ")
	require.NoError(t, err)
	require.Equal(t, " foo ", key)
	require.Equal(t, " bar = baz ", value)

	key, value, err = splitKeyValue("foo=bar=baz")
	require.NoError(t, err)
	require.Equal(t, "foo", key)
	require.Equal(t, "bar=baz", value)

	key, value, err = splitKeyValue("a[b]=c")
	require.NoError(t, err)
	require.Equal(t, "a[b]", key)
	require.Equal(t, "c", value)

	key, value, err = splitKeyValue("a[b][c]=d")
	require.NoError(t, err)
	require.Equal(t, "a[b][c]", key)
	require.Equal(t, "d", value)
}

func TestUnmarshal(t *testing.T) {

	ConvertArrays(true)
	hash, err := Unmarshal("foo")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": nil})
	}

	hash, err = Unmarshal("foo=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": ""})
	}

	hash, err = Unmarshal("foo=bar")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "bar"})
	}

	hash, err = Unmarshal("foo=\"bar\"")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "\"bar\""})
	}

	hash, err = Unmarshal("foo=bar&foo=quux")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "quux"})
	}

	hash, err = Unmarshal("foo&foo=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": ""})
	}

	hash, err = Unmarshal("foo=1&bar=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "1", "bar": "2"})
	}

	hash, err = Unmarshal("&foo=1&&bar=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "1", "bar": "2"})
	}

	hash, err = Unmarshal("foo&bar=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": nil, "bar": ""})
	}

	hash, err = Unmarshal("foo=bar&baz=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "bar", "baz": ""})
	}

	hash, err = Unmarshal("my+weird+field=q1%212%22%27w%245%267%2Fz8%29%3F")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"my weird field": `q1!2"'w$5&7/z8)?`})
	}

	hash, err = Unmarshal("a=b&pid%3D1234=1023")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"pid=1234": "1023", "a": "b"})
	}

	hash, err = Unmarshal("foo[]")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{nil}})
	}

	hash, err = Unmarshal("foo[]=")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{""}})
	}

	hash, err = Unmarshal("foo[]=bar")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{"bar"}})
	}

	hash, err = Unmarshal("foo[]=1&foo[]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{"1", "2"}})
	}

	hash, err = Unmarshal("foo=bar&baz[]=1&baz[]=2&baz[]=3")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": "bar", "baz": []interface{}{"1", "2", "3"}})
	}

	hash, err = Unmarshal("foo[]=bar&baz[]=1&baz[]=2&baz[]=3")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"foo": []interface{}{"bar"}, "baz": []interface{}{"1", "2", "3"}})
	}

	hash, err = Unmarshal("x[y][z]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": "1"}}})
	}

	hash, err = Unmarshal("x[y][z][]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": []interface{}{"1"}}}})
	}

	hash, err = Unmarshal("x[y][z]=1&x[y][z]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": "2"}}})
	}

	hash, err = Unmarshal("x[y][z][]=1&x[y][z][]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": map[string]interface{}{"z": []interface{}{"1", "2"}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1"}}}})
	}

	hash, err = Unmarshal("x[y][][z][]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": []interface{}{"1"}}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][w]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "w": "2"}}}})
	}

	hash, err = Unmarshal("x[y][][v][w]=1")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"v": map[string]interface{}{"w": "1"}}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][v][w]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "v": map[string]interface{}{"w": "2"}}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][z]=2")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1"}, map[string]interface{}{"z": "2"}}}})
	}

	hash, err = Unmarshal("x[y][][z]=1&x[y][][w]=a&x[y][][z]=2&x[y][][w]=3")
	if assert.NoError(t, err) {
		assert.Equal(t, hash, map[string]interface{}{"x": map[string]interface{}{"y": []interface{}{map[string]interface{}{"z": "1", "w": "a"}, map[string]interface{}{"z": "2", "w": "3"}}}})
	}

	hash, err = Unmarshal("x[y]=1&x[y]z=2")
	assert.Error(t, err)

	hash, err = Unmarshal("x[y]=1&x[]=1")
	assert.Error(t, err)

	hash, err = Unmarshal("x[y]=1&x[y][][w]=2")
	assert.Error(t, err)

	ConvertArrays(false)
}

func TestUnmarshalAndConvert(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	expected := map[string]interface{}{
		"split_rules": []interface{}{
			map[string]interface{}{
				"amount":                "29386",
				"charge_processing_fee": "false",
				"liable":                "false",
				"recipient_id":          "re_cj16dgfc1046qbq60q5x8sslx",
			},
			map[string]interface{}{
				"amount":                "5597",
				"charge_processing_fee": "true",
				"liable":                "true",
				"recipient_id":          "re_ciqcoztri002q4a603yfhqvxq",
			},
		},
	}

	hash, err := Unmarshal("split_rules[0][amount]=29386&split_rules[0][recipient_id]=re_cj16dgfc1046qbq60q5x8sslx&split_rules[0][charge_processing_fee]=false&split_rules[0][liable]=false&split_rules[1][amount]=5597&split_rules[1][recipient_id]=re_ciqcoztri002q4a603yfhqvxq&split_rules[1][charge_processing_fee]=true&split_rules[1][liable]=true")
	require.NoError(t, err, "could not unmarshal")
	require.NotNil(t, hash)
	require.Equal(t, expected, hash)

}

func TestSimpleString(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var expected, actual interface{}
	var err error

	t.Run("foo=c++", func(t *testing.T) {
		actual, err = Parse("foo=c++")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo": "c  ",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[>=]=23", func(t *testing.T) {
		actual, err = Parse("a[>=]=23")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				">=": "23",
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[<=>]==23", func(t *testing.T) {
		actual, err = Parse("a[<=>]==23")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"<=>": "=23",
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[==]=23", func(t *testing.T) {
		actual, err = Parse("a[==]=23")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"==": "23",
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("foo", func(t *testing.T) {
		actual, err = Parse("foo")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo": "",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("foo=", func(t *testing.T) {
		actual, err = Parse("foo=")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo": "",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("foo=bar", func(t *testing.T) {
		actual, err = Parse("foo=bar")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo": "bar",
		}
		require.Equal(t, expected, actual)
	})

	t.Run(" foo = bar = baz ", func(t *testing.T) {
		actual, err = Parse(" foo = bar = baz ")
		require.NoError(t, err)
		expected = map[string]interface{}{
			" foo ": " bar = baz ",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("foo=bar=baz", func(t *testing.T) {
		actual, err = Parse("foo=bar=baz")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo": "bar=baz",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("foo=bar&bar=baz", func(t *testing.T) {
		actual, err = Parse("foo=bar&bar=baz")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo": "bar",
			"bar": "baz",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("foo2=bar2&baz2=", func(t *testing.T) {
		actual, err = Parse("foo2=bar2&baz2=")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo2": "bar2",
			"baz2": "",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("foo=bar&baz", func(t *testing.T) {
		actual, err = Parse("foo=bar&baz")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"foo": "bar",
			"baz": "",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("cht=p3&chd=t:60,40&chs=250x100&chl=Hello|World", func(t *testing.T) {
		actual, err = Parse("cht=p3&chd=t:60,40&chs=250x100&chl=Hello|World")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"cht": "p3",
			"chd": "t:60,40",
			"chs": "250x100",
			"chl": "Hello|World",
		}
		require.Equal(t, expected, actual)
	})
}

func TestSimpleArray(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("a=b&a=c", func(t *testing.T) {
		actual, err = Parse("a=b&a=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})
}

func TestExplicitArray(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("a[]=b", func(t *testing.T) {
		actual, err = Parse("a[]=b")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[]=b&a[]=c", func(t *testing.T) {
		actual, err = Parse("a[]=b&a[]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[]=b&a[]=c&a[]=d", func(t *testing.T) {
		actual, err = Parse("a[]=b&a[]=c&a[]=d")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c", "d"},
		}
		require.Equal(t, expected, actual)
	})
}

func TestNestedString(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("a[b]=c", func(t *testing.T) {
		actual, err = Parse("a[b]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"b": "c",
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[b][c]=d", func(t *testing.T) {
		actual, err = Parse("a[b][c]=d")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "d",
				},
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[b][c][d][e][f][g][h]=i", func(t *testing.T) {
		actual, err = Parse("a[b][c][d][e][f][g][h]=i")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": map[string]interface{}{
						"d": map[string]interface{}{
							"e": map[string]interface{}{
								"f": map[string]interface{}{
									"g": map[string]interface{}{
										"h": "i",
									},
								},
							},
						},
					},
				},
			},
		}
		require.Equal(t, expected, actual)
	})
}

func TestSimpleAndExplicitArray(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("a=b&a[]=c", func(t *testing.T) {
		actual, err = Parse("a=b&a[]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[]=b&a=c", func(t *testing.T) {
		actual, err = Parse("a[]=b&a=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[0]=b&a=c", func(t *testing.T) {
		actual, err = Parse("a[0]=b&a=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a=b&a[0]=c", func(t *testing.T) {
		actual, err = Parse("a=b&a[0]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[1]=b&a=c", func(t *testing.T) {
		actual, err = Parse("a[1]=b&a=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[]=b&a=c", func(t *testing.T) {
		actual, err = Parse("a[]=b&a=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a=b&a[1]=c", func(t *testing.T) {
		actual, err = Parse("a=b&a[1]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a=b&a[]=c", func(t *testing.T) {
		actual, err = Parse("a=b&a[]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})
}

func TestNestedArrays(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("a[b][]=c&a[b][]=d", func(t *testing.T) {
		actual, err = Parse("a[b][]=c&a[b][]=d")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"b": []interface{}{"c", "d"},
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[>=]=25", func(t *testing.T) {
		actual, err = Parse("a[>=]=25")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				">=": "25",
			},
		}
		require.Equal(t, expected, actual)
	})
}

func TestSpecifyArrayIndices(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("a[1]=c&a[0]=b&a[2]=d", func(t *testing.T) {
		actual, err = Parse("a[1]=c&a[0]=b&a[2]=d")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c", "d"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[1]=c&a[0]=b", func(t *testing.T) {
		actual, err = Parse("a[1]=c&a[0]=b")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"b", "c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[1]=c", func(t *testing.T) {
		actual, err = Parse("a[1]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"c"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[20]=a", func(t *testing.T) {
		actual, err = Parse("a[20]=a")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": []interface{}{"a"},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[12b]=c", func(t *testing.T) {
		actual, err = Parse("a[12b]=c")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"12b": "c",
			},
		}
		require.Equal(t, expected, actual)
	})
}

func TestEnconding(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("he%3Dllo=th%3Dere", func(t *testing.T) {
		actual, err = Parse("he%3Dllo=th%3Dere")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"he=llo": "th=ere",
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[b%20c]=d", func(t *testing.T) {
		actual, err = Parse("a[b%20c]=d")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"b c": "d",
			},
		}
		require.Equal(t, expected, actual)
	})

	t.Run("a[b]=c%20d", func(t *testing.T) {
		actual, err = Parse("a[b]=c%20d")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"a": map[string]interface{}{
				"b": "c d",
			},
		}
		require.Equal(t, expected, actual)
	})
}

func TestBrackets(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run(`pets=["tobi"]`, func(t *testing.T) {
		actual, err = Parse(`pets=["tobi"]`)
		require.NoError(t, err)
		expected = map[string]interface{}{
			"pets": `["tobi"]`,
		}
		require.Equal(t, expected, actual)
	})

	t.Run(`operators=[">=", "<="]`, func(t *testing.T) {
		actual, err = Parse(`operators=[">=", "<="]`)
		require.NoError(t, err)
		expected = map[string]interface{}{
			"operators": `[">=", "<="]`,
		}
		require.Equal(t, expected, actual)
	})
}

func TestEmptyValues(t *testing.T) {
	ConvertArrays(true)
	defer ConvertArrays(false)

	var actual, expected interface{}
	var err error

	t.Run("empty", func(t *testing.T) {
		actual, err = Parse("")
		require.NoError(t, err)
		expected = map[string]interface{}{}
		require.Equal(t, expected, actual)
	})

	t.Run("spaces", func(t *testing.T) {
		actual, err = Parse("  ")
		require.NoError(t, err)
		expected = map[string]interface{}{
			"  ": "",
		}
		require.Equal(t, expected, actual)
	})

}
