// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package main

import (
	"testing"
)

func TestParseLustreMetadataOperations(t *testing.T) {

	var data string = `{"status":"success","data":{"resultType":"vector","result":[
		{"metric":{"jobid":"35044931","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},
		{"metric":{"jobid":"35070653","target":"hebe-MDT0000"},"value":[1639743019.545,"43"]},
		{"metric":{"jobid":"35189820","target":"hebe-MDT0000"},"value":[1639743019.545,"4"]},
		{"metric":{"jobid":"35166602","target":"hebe-MDT0001"},"value":[1639743019.545,"31"]},
		{"metric":{"jobid":"35189845","target":"hebe-MDT0001"},"value":[1639743019.545,"1"]},
		{"metric":{"jobid":"35048662","target":"hebe-MDT0001"},"value":[1639743019.545,"27"]},
		{"metric":{"jobid":"cp.5689","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},
		{"metric":{"jobid":"35056989","target":"hebe-OST022d"},"value":[1639743019.545,"5"]},
		{"metric":{"jobid":"touch.6812","target":"hebe-OST020c"},"value":[1639743019.545,"1"]}
		]}}`

	var content []byte = []byte(data)

	var lustreMetadataOperations *[]metadataInfo
	var err error

	lustreMetadataOperations, err = parseLustreMetadataOperations(&content)

	if err != nil {
		t.Error(err)
	}

	var got_count int = len(*lustreMetadataOperations)
	var expected_count int = 7

	if expected_count != got_count {
		t.Errorf("Expected count of metadata operations: %d - got: %d", expected_count, got_count)
	}

	var metadataInfo metadataInfo = (*lustreMetadataOperations)[0]
	var expected_jobid string = "35044931"
	var expected_target string = "hebe-MDT0002"

	if metadataInfo.jobid != expected_jobid {
		t.Errorf("Expected jobid: %s - got: %s", expected_jobid, metadataInfo.jobid)
	}

	if metadataInfo.target != expected_target {
		t.Errorf("Expected target: %s - got: %s", expected_target, metadataInfo.target)
	}

	for _, metadataInfo := range *lustreMetadataOperations {
		if !regexMetadataMDT.MatchString(metadataInfo.target) {
			t.Error("Only MDT as target is allowed:", metadataInfo.target)
		}
	}

}

func TestResolveProcInfo(t *testing.T) {

	users := userInfoMap{
		1001: userInfo{user: "alice", uid: 1001, gid: 100},
		1002: userInfo{user: "carol", uid: 1002, gid: 999}, // GID not in groups
	}

	groups := groupInfoMap{
		100: groupInfo{group: "staff", gid: 100},
	}

	// Simple procname.uid — verify all returned fields
	info, err := resolveProcInfo("cp.1001", users, groups)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("Expected non-nil procInfo for cp.1001")
	}
	if info.procName != "cp" {
		t.Errorf("Expected procName 'cp', got '%s'", info.procName)
	}
	if info.userName != "alice" {
		t.Errorf("Expected userName 'alice', got '%s'", info.userName)
	}
	if info.groupName != "staff" {
		t.Errorf("Expected groupName 'staff', got '%s'", info.groupName)
	}

	// Dotted procname — only procName parsing differs
	info, err = resolveProcInfo("my.app.1001", users, groups)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("Expected non-nil procInfo for my.app.1001")
	}
	if info.procName != "my.app" {
		t.Errorf("Expected procName 'my.app', got '%s'", info.procName)
	}

	// No dot separator → skip (nil, nil)
	info, err = resolveProcInfo("nodot", users, groups)
	if err != nil {
		t.Errorf("Unexpected error for 'nodot': %v", err)
	}
	if info != nil {
		t.Error("Expected nil procInfo for 'nodot'")
	}

	// Non-numeric UID → error
	info, err = resolveProcInfo("cp.notanumber", users, groups)
	if err == nil {
		t.Error("Expected error for non-numeric UID")
	}
	if info != nil {
		t.Error("Expected nil procInfo for non-numeric UID")
	}

	// Unknown UID → skip (nil, nil)
	info, err = resolveProcInfo("cp.9999", users, groups)
	if err != nil {
		t.Errorf("Unexpected error for unknown UID: %v", err)
	}
	if info != nil {
		t.Error("Expected nil procInfo for unknown UID")
	}

	// Known UID but unknown GID → skip (nil, nil)
	info, err = resolveProcInfo("cp.1002", users, groups)
	if err != nil {
		t.Errorf("Unexpected error for unknown GID: %v", err)
	}
	if info != nil {
		t.Error("Expected nil procInfo for unknown GID")
	}
}

func TestParseLustreTotalBytes(t *testing.T) {

	var data string = `{"status":"success","data":{"resultType":"vector","result":[
		{"metric":{"jobid":"35652133"},"value":[1640181380.814,"319215.8800539506"]},
		{"metric":{"jobid":"35239994"},"value":[1640181380.814,"125747.2"]},
		{"metric":{"jobid":"35651038"},"value":[1640181380.814,"379697.46153350436"]},
		{"metric":{"jobid":"35683050"},"value":[1640181380.814,"955.7333333333333"]},
		{"metric":{"jobid":"35676304"},"value":[1640181380.814,"893883.7333333333"]},
		{"metric":{"jobid":"35682305"},"value":[1640181380.814,"819.2"]},
		{"metric":{"jobid":"35676288"},"value":[1640181380.814,"689493.3333333334"]},
		{"metric":{"jobid":"35676299"},"value":[1640181380.814,"248627.2"]}
		]}}`

	var content []byte = []byte(data)

	var lustreThroughputInfo *[]throughputInfo
	var err error

	lustreThroughputInfo, err = parseLustreTotalBytes(&content)

	if err != nil {
		t.Error(err)
	}

	var got_count int = len(*lustreThroughputInfo)
	var expected_count int = 8

	if expected_count != got_count {
		t.Errorf("Expected count of metadata operations: %d - got: %d", expected_count, got_count)
	}

	var throughputInfo throughputInfo = (*lustreThroughputInfo)[0]
	var expected_jobid string = "35652133"

	if throughputInfo.jobid != expected_jobid {
		t.Errorf("Expected jobid: %s - got: %s", expected_jobid, throughputInfo.jobid)
	}
}
