package main

import "testing"

func TestParseLustreMetadataOperations(t *testing.T) {

	var data string = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"jobid":"35044931","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35070653","target":"hebe-MDT0000"},"value":[1639743019.545,"43"]},{"metric":{"jobid":"35189820","target":"hebe-MDT0000"},"value":[1639743019.545,"4"]},{"metric":{"jobid":"35166602","target":"hebe-MDT0001"},"value":[1639743019.545,"31"]},{"metric":{"jobid":"35189845","target":"hebe-MDT0001"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35048662","target":"hebe-MDT0001"},"value":[1639743019.545,"27"]},{"metric":{"jobid":"35097923","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35156251","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35178783","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35170040","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"touch.6812","target":"hebe-MDT0002"},"value":[1639743019.545,"3"]},{"metric":{"jobid":"xrootd.6812","target":"hebe-MDT0002"},"value":[1639743019.545,"13"]},{"metric":{"jobid":"rsync.9334","target":"hebe-MDT0000"},"value":[1639743019.545,"15"]},{"metric":{"jobid":"34842227","target":"hebe-MDT0001"},"value":[1639743019.545,"6"]},{"metric":{"jobid":"35048664","target":"hebe-MDT0001"},"value":[1639743019.545,"28"]},{"metric":{"jobid":"35178762","target":"hebe-MDT0001"},"value":[1639743019.545,"12"]},{"metric":{"jobid":"35166599","target":"hebe-MDT0001"},"value":[1639743019.545,"72"]},{"metric":{"jobid":"35166603","target":"hebe-MDT0001"},"value":[1639743019.545,"109"]},{"metric":{"jobid":"35168120","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35186134","target":"hebe-MDT0002"},"value":[1639743019.545,"2"]},{"metric":{"jobid":"35129848","target":"hebe-MDT0002"},"value":[1639743019.545,"2"]},{"metric":{"jobid":"touch.6812","target":"hebe-OST020c"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35070627","target":"hebe-MDT0000"},"value":[1639743019.545,"101"]},{"metric":{"jobid":"35189850","target":"hebe-MDT0000"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"slurmstepd.10388","target":"hebe-MDT0000"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"ll_sa_11831.0","target":"hebe-MDT0001"},"value":[1639743019.545,"185"]},{"metric":{"jobid":"35156968","target":"hebe-MDT0000"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35188793","target":"hebe-MDT0000"},"value":[1639743019.545,"2"]},{"metric":{"jobid":"35189309","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35185982","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"Reloader.7377","target":"hebe-MDT0000"},"value":[1639743019.545,"2"]},{"metric":{"jobid":"sftp-server.5524","target":"hebe-MDT0001"},"value":[1639743019.545,"734"]},{"metric":{"jobid":"35039426","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35154964","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35177818","target":"hebe-MDT0001"},"value":[1639743019.545,"15"]},{"metric":{"jobid":"35032996","target":"hebe-MDT0002"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35112154","target":"hebe-MDT0002"},"value":[1639743019.545,"2"]},{"metric":{"jobid":"35188756","target":"hebe-MDT0000"},"value":[1639743019.545,"1"]},{"metric":{"jobid":"35048628","target":"hebe-MDT0001"},"value":[1639743019.545,"83"]},{"metric":{"jobid":"35051717","target":"hebe-MDT0001"},"value":[1639743019.545,"13"]},{"metric":{"jobid":"35056989","target":"hebe-MDT0001"},"value":[1639743019.545,"5"]}]}}`
	var content []byte = []byte(data)

	var lustreMetadataOperations *[]metadataInfo = parseLustreMetadataOperations(&content)

	var got_count int = len(*lustreMetadataOperations)
	var expected_count int = 41

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
}
