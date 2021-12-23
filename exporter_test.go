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
