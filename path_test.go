// Copyright 2013 Julien Schmidt. All rights reserved.
// Based on the path package, Copyright 2009 The Go Authors.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package httprouter

import (
	"bytes"
	"testing"
)

type cleanPathTest struct {
	path, result []byte
}

var cleanTests = []cleanPathTest{
	// Already clean
	{[]byte("/"), []byte("/")},
	{[]byte("/abc"), []byte("/abc")},
	{[]byte("/a/b/c"), []byte("/a/b/c")},
	{[]byte("/abc/"), []byte("/abc/")},
	{[]byte("/a/b/c/"), []byte("/a/b/c/")},

	// missing root
	{[]byte(""), []byte("/")},
	{[]byte("a/"), []byte("/a/")},
	{[]byte("abc"), []byte("/abc")},
	{[]byte("abc/def"), []byte("/abc/def")},
	{[]byte("a/b/c"), []byte("/a/b/c")},

	// Remove doubled slash
	{[]byte("//"), []byte("/")},
	{[]byte("/abc//"), []byte("/abc/")},
	{[]byte("/abc/def//"), []byte("/abc/def/")},
	{[]byte("/a/b/c//"), []byte("/a/b/c/")},
	{[]byte("/abc//def//ghi"), []byte("/abc/def/ghi")},
	{[]byte("//abc"), []byte("/abc")},
	{[]byte("///abc"), []byte("/abc")},
	{[]byte("//abc//"), []byte("/abc/")},

	// Remove . elements
	{[]byte("."), []byte("/")},
	{[]byte("./"), []byte("/")},
	{[]byte("/abc/./def"), []byte("/abc/def")},
	{[]byte("/./abc/def"), []byte("/abc/def")},
	{[]byte("/abc/."), []byte("/abc/")},

	// Remove .. elements
	{[]byte(".."), []byte("/")},
	{[]byte("../"), []byte("/")},
	{[]byte("../../"), []byte("/")},
	{[]byte("../.."), []byte("/")},
	{[]byte("../../abc"), []byte("/abc")},
	{[]byte("/abc/def/ghi/../jkl"), []byte("/abc/def/jkl")},
	{[]byte("/abc/def/../ghi/../jkl"), []byte("/abc/jkl")},
	{[]byte("/abc/def/.."), []byte("/abc")},
	{[]byte("/abc/def/../.."), []byte("/")},
	{[]byte("/abc/def/../../.."), []byte("/")},
	{[]byte("/abc/def/../../.."), []byte("/")},
	{[]byte("/abc/def/../../../ghi/jkl/../../../mno"), []byte("/mno")},

	// Combinations
	{[]byte("abc/./../def"), []byte("/def")},
	{[]byte("abc//./../def"), []byte("/def")},
	{[]byte("abc/../../././../def"), []byte("/def")},
}

func TestPathClean(t *testing.T) {
	for _, test := range cleanTests {
		if s := CleanPathB(test.path); !bytes.Equal(s, test.result) {
			t.Errorf("CleanPath(%q) = %q, want %q", test.path, s, test.result)
		}
		if s := CleanPathB(test.result); !bytes.Equal(s, test.result) {
			t.Errorf("CleanPath(%q) = %q, want %q", test.result, s, test.result)
		}
	}
}

func TestPathCleanMallocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}
	for _, test := range cleanTests {
		allocs := testing.AllocsPerRun(100, func() { CleanPathB(test.result) })
		if allocs > 0 {
			//t.Errorf("CleanPath(%q): %v allocs, want zero", test.result, allocs)
		}
	}
}

func BenchmarkPathClean(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, test := range cleanTests {
			CleanPathB(test.path)
		}
	}
}

func genLongPaths() (testPaths []cleanPathTest) {
	for i := 1; i <= 1234; i++ {
		ss := bytes.Repeat([]byte("a"), i)

		correctPath := append([]byte("/"), ss...)
		testPaths = append(testPaths, cleanPathTest{
			path:   correctPath,
			result: correctPath,
		}, cleanPathTest{
			path:   ss,
			result: correctPath,
		}, cleanPathTest{
			path:   append([]byte("//"), ss...),
			result: correctPath,
		}, cleanPathTest{
			path:   append(append([]byte("/"), ss...), []byte("/b/..")...),
			result: correctPath,
		})
	}
	return
}

func TestPathCleanLong(t *testing.T) {
	cleanTests := genLongPaths()

	for _, test := range cleanTests {
		if s := CleanPathB(test.path); !bytes.Equal(s, test.result) {
			t.Errorf("CleanPath(%q) = %q, want %q", test.path, s, test.result)
		}
		if s := CleanPathB(test.result); !bytes.Equal(s, test.result) {
			t.Errorf("CleanPath(%q) = %q, want %q", test.result, s, test.result)
		}
	}
}

func BenchmarkPathCleanLong(b *testing.B) {
	cleanTests := genLongPaths()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, test := range cleanTests {
			CleanPathB(test.path)
		}
	}
}
