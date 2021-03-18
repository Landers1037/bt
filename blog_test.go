/*
landers Apps
Author: landers
Github: github.com/landers1037
*/

package main

import (
	"testing"
)

func TestParseMd(t *testing.T) {
	f, e := ReadMd("test.md")
	res := ParseMd(f)
	if e != nil {
		t.Error("error")
	}else {
		t.Logf("\n%s", res)
	}
}

func TestParseAbs(t *testing.T) {
	f, e := ReadMd("test.md")
	res := ParseMd(f)
	if e != nil {
		t.Error("error")
	}else {
		res = ParseAbs(res)
		t.Logf("\n%s", res)
	}
}

func TestParseMdAbs(t *testing.T) {
	f, e := ReadMd("test.md")
	res := ParseMd(f)
	if e != nil {
		t.Error("error")
	}else {
		res = ParseMdAbs(res)
		t.Logf("\n%s", res)
	}
}

func TestParseMeta(t *testing.T) {
	f, e := ReadMd("test.md")
	if e != nil {
		t.Error("error")
	}else {
		res := ParseMeta(f)
		t.Logf("\n%+v", res)
	}
}

func TestParseMetaYaml(t *testing.T) {
	f, e := ReadMd("test.md")
	if e != nil {
		t.Error("error")
	}else {
		res := ParseMetaYaml(f)
		t.Logf("\n%+v", res)
	}
}

func TestParseYmalFront(t *testing.T) {
	f, e := ReadMd("test.md")
	if e != nil {
		t.Error("error")
	}else {
		res := ParseYmalFront(f)
		t.Logf("\n%+v", res)
	}
}

func TestNewMd(t *testing.T) {
	e := NewMd("new")
	if e != nil {
		t.Errorf(e.Error())
	}else {
		t.Log("success create md")
	}
}

func TestProcessAllMd(t *testing.T) {
	var tmpList = [100]string{}
	ProcessAllMd(tmpList[:100], 0)
}
