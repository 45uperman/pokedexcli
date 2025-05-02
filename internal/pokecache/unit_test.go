package pokecache

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	testCache := NewCache(30 * time.Second)
	url := "https://pokeapi.co/api/v2/location-area"

	data, ok := testCache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			t.Errorf("error making request: %s", err)
			t.FailNow()
		}

		data, err = io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			t.Errorf("error reading response: %s", err)
			t.FailNow()
		}
	}

	testCache.Add(url, data)

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    url,
			expected: string(data),
		},
	}

	for _, c := range cases {
		testCache := NewCache(1 * time.Second)
		url := c.input

		data, ok := testCache.Get(url)
		if !ok {
			res, err := http.Get(url)
			if err != nil {
				t.Errorf("error making request: %s", err)
				t.FailNow()
			}

			data, err = io.ReadAll(res.Body)
			res.Body.Close()

			if err != nil {
				t.Errorf("error reading response: %s", err)
				t.FailNow()
			}
		}

		testCache.Add(url, data)

		data, ok = testCache.Get(url)
		if !ok {
			t.Errorf("pokecache failed to cache response")
			t.FailNow()
		}
		if c.expected != string(data) {
			t.Errorf("Bad cached data\nEXPECTED: %s\nACTUAL: %s", c.expected, data)
			t.FailNow()
		}

		time.Sleep(2 * time.Second)
		_, ok = testCache.Get(url)
		if ok {
			t.Errorf("pokecache failed to reap data")
			t.FailNow()
		}

	}
}
