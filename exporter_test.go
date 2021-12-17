package main

import "testing"

func TestParseLustreMetadataOperations(t *testing.T) {

	var data string = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"jobid":"35056989"},"value":[1639742610.832,"5"]},{"metric":{"jobid":"35178762"},"value":[1639742610.832,"12"]},{"metric":{"jobid":"35112130"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35149256"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35154427"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35165999"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35186135"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35156968"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35177818"},"value":[1639742610.832,"12"]},{"metric":{"jobid":"35157376"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35170692"},"value":[1639742610.832,"15"]},{"metric":{"jobid":"35189617"},"value":[1639742610.832,"5"]},{"metric":{"jobid":"35051717"},"value":[1639742610.832,"9"]},{"metric":{"jobid":"35154904"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35170162"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35181417"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35186164"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35189034"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"touch.6812"},"value":[1639742610.832,"3"]},{"metric":{"jobid":"35174086"},"value":[1639742610.832,"30"]},{"metric":{"jobid":"sftp-server.5524"},"value":[1639742610.832,"803"]},{"metric":{"jobid":"35112232"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35189103"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35171060"},"value":[1639742610.832,"153"]},{"metric":{"jobid":"34842227"},"value":[1639742610.832,"6"]},{"metric":{"jobid":"35166362"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35166601"},"value":[1639742610.832,"67"]},{"metric":{"jobid":"35111494"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35178285"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35185311"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35070620"},"value":[1639742610.832,"202"]},{"metric":{"jobid":"wc.0"},"value":[1639742610.832,"212"]},{"metric":{"jobid":"35111615"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35165763"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35167514"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35189142"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"xrootd.6812"},"value":[1639742610.832,"16"]},{"metric":{"jobid":"Reloader.7377"},"value":[1639742610.832,"2"]},{"metric":{"jobid":"35048662"},"value":[1639742610.832,"60"]},{"metric":{"jobid":"ll_sa_11831.0"},"value":[1639742610.832,"305"]},{"metric":{"jobid":"35033493"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35189161"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35189618"},"value":[1639742610.832,"3"]},{"metric":{"jobid":"35048628"},"value":[1639742610.832,"85"]},{"metric":{"jobid":"35166599"},"value":[1639742610.832,"319"]},{"metric":{"jobid":"35114865"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35156188"},"value":[1639742610.832,"1"]},{"metric":{"jobid":"35189250"},"value":[1639742610.832,"1"]}]}}`
	var content []byte = []byte(data)

	var lustreMetadataOperations *[]metadataInfo = parseLustreMetadataOperations(&content)

	var got int = len(*lustreMetadataOperations)
	var expected int = 48

	if expected != got {
		t.Errorf("Expected count of metadata operations: %d - got %d", expected, got)
	}
}
