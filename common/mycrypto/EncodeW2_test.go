package mycrypto

import (
	"fmt"
	"github.com/tidusant/c3m/common/log"

	"testing"
)

//test special char
func TestEncodeW2Normkey(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeW2Normkey ====")
	//check test data

	key := `abc`
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW2(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW2(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}
func TestEncodeW2Longkey(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeW2Longkey ====")
	//check test data

	key := `www.test-long-test.com/abc`
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW2(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW2(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}
func TestEncodeW2SpecialChar(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeW2SpecialChar ====")
	//check test data

	key := `www.test-long-test.com/abc` + `.*?/~!@#$%^&*(),.[];'{}<>:\"` + "`" + "\n\t "
	data := Base64Encode(text + key)
	log.Debugf("b64: %s", data)
	rs := EncodeW2(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW2(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}
func TestEncodeW2EmptyKey(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeW2EmptyKey ====")
	//check test data
	key := ``
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW2(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW2(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}

func TestEncodeW2Space(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeW2Space ====")
	//check test data
	key := ` `
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW2(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW2(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}

//func TestEncodeW2FromJS(t *testing.T){
//	fmt.Println("\n\n==== TestEncdoeWFromJS ====")
//	rs:=DecodeW2("pDIgUZZ29dmQSqlwtpmP8plP9fmFgqlrKrIFgyIFHbRNgelvEcRNufl3EcRPuiRP5emLKijF1fkQtplvOdIP9cRQAunQZdRNu0RPqqlhKhk290lhKykrKqRQKyIFWuRP9vRPWbHGWijFWqkLKVHGAykrKbjGAulvO0mGSuRPIhk20pWMDpZtVbRP1qj2udIhKymLKemvEhRMRfVMJpnFEqlwVpk2gtUrKBjFWxHGStRN1sZ2gykwAeH2bbRPNpCPO0jF4plQSeIvEil29hRPO0RNqqkGKtIF4cD3utkvE5RNWekPguI2DpjF4pEvuhI2udjFNbRPgek2cuILK1lLKekvDpk2HpmPquRP1elvDpk2SiH3EhIBKVHGAykrK3k3StlhfpH29dl2EsmPE0mGRbRPIhk20pHBKVk3SukBKSlQW1kBKfHGWiHFmuULKqkvZpI29ykvlpmPqhk3EwjLK0jPDpH2u0IGVpk2HpmPquRQmelvZpjF4pH2gql3WyH2ObRPgyMwmPEhHGA1lvDbRPAyl2WemvEhIFZpmPquRQEdIP91HwAqHvguRQWemGSsIB4pCP9hIF0pBGKimF0pH29cIGVpIwSekBKiIFW0jF9dlhJgUsNfUsVhRPOdILJgUsNfUsViRP9vRLStIBKPjF5yHwEiRNSekv9hmF0pIGZpCFObk3S1kBRpTOAxIBKOnQAhIF1ulhKeIrKQk29tRPOdILKOmvubTBKrnBKMjFWulv8bRQmhjGA0IF4pjF4pWMDpZtVdROAxjGVpHv9ejhKylhKqRQAhIFO0jGWuRP9dRQAxIBK0jPEelwtpk2HpIGAxjFWiULK2IGS5RQKelQEbHGRpIQEhjF5wRQAxIBKBIF5qjGWiHF5sIB4pEPquRPIylwW0RPgykvDpk2HpCP9hIF0pBGKimF0bRLSVk3SukBKylQW1kBKtk2gelrKijGZpHF1umL4dRrfpH29cIGVpIwSekBKqRPgykvDpjF4pl2EsmPuekrJgUsNfUsVhUwAxjGVpmPEimJ","this test")
//	if rs != text {
//		t.Fatalf("Test fail")
//	}
//}
