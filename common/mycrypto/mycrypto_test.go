package mycrypto

import (
	"fmt"
	"github.com/tidusant/c3m/common/log"
	"os"
	"testing"
)

var (
	text  = `Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`
	text1 = `ontrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`
	text2 = `ntrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`
	text3 = `trary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`
	text4 = `rary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`
	text5 = `ary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`
	text6 = `ry to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.`
)

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setup() {

}

//test special char
func TestEncodeWNormkey(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeWNormkey ====")
	//check test data

	key := `www.test-long-test.com/abc`
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}
func TestEncodeWLongkey(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeWLongkey ====")
	//check test data

	key := `www.test-long-test.com/abc`
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}
func TestEncodeWSpecialChar(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeWSpecialChar ====")
	//check test data

	key := `www.test-long-test.com/abc` + `.*?/~!@#$%^&*(),.[];'{}<>:\"` + "`" + "\n\t "
	data := Base64Encode(text + key)
	log.Debugf("b64: %s", data)
	rs := EncodeW(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}
func TestEncodeWEmptyKey(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeWEmptyKey ====")
	//check test data
	key := ``
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}

func TestEncodeWSpace(t *testing.T) {
	fmt.Println("\n\n==== TestEncodeWSpace ====")
	//check test data
	key := ` `
	//data=mycrypto.Base64Encode(data)
	log.Debugf("b64: %s", text)
	rs := EncodeW(text, key)
	log.Debugf("VEncode: %s", rs)
	rs = DecodeW(rs, key)
	log.Debugf("VDecode: %s", rs)
	if rs != text {
		t.Fatalf("Test fail")
	}
}
func TestEncdoeWFromJS(t *testing.T) {
	fmt.Println("\n\n==== TestEncdoeWFromJS ====")
	rs := DecodeW("pDIgUZZ29dmQSqlwtpmP8plP9fmFgqlrKrIFgyIFHbRNgelvEcRNufl3EcRPuiRP5emLKijF1fkQtplvOdIP9cRQAunQZdRNu0RPqqlhKhk290lhKykrKqRQKyIFWuRP9vRPWbHGWijFWqkLKVHGAykrKbjGAulvO0mGSuRPIhk20pWMDpZtVbRP1qj2udIhKymLKemvEhRMRfVMJpnFEqlwVpk2gtUrKBjFWxHGStRN1sZ2gykwAeH2bbRPNpCPO0jF4plQSeIvEil29hRPO0RNqqkGKtIF4cD3utkvE5RNWekPguI2DpjF4pEvuhI2udjFNbRPgek2cuILK1lLKekvDpk2HpmPquRP1elvDpk2SiH3EhIBKVHGAykrK3k3StlhfpH29dl2EsmPE0mGRbRPIhk20pHBKVk3SukBKSlQW1kBKfHGWiHFmuULKqkvZpI29ykvlpmPqhk3EwjLK0jPDpH2u0IGVpk2HpmPquRQmelvZpjF4pH2gql3WyH2ObRPgyMwmPEhHGA1lvDbRPAyl2WemvEhIFZpmPquRQEdIP91HwAqHvguRQWemGSsIB4pCP9hIF0pBGKimF0pH29cIGVpIwSekBKiIFW0jF9dlhJgUsNfUsVhRPOdILJgUsNfUsViRP9vRLStIBKPjF5yHwEiRNSekv9hmF0pIGZpCFObk3S1kBRpTOAxIBKOnQAhIF1ulhKeIrKQk29tRPOdILKOmvubTBKrnBKMjFWulv8bRQmhjGA0IF4pjF4pWMDpZtVdROAxjGVpHv9ejhKylhKqRQAhIFO0jGWuRP9dRQAxIBK0jPEelwtpk2HpIGAxjFWiULK2IGS5RQKelQEbHGRpIQEhjF5wRQAxIBKBIF5qjGWiHF5sIB4pEPquRPIylwW0RPgykvDpk2HpCP9hIF0pBGKimF0bRLSVk3SukBKylQW1kBKtk2gelrKijGZpHF1umL4dRrfpH29cIGVpIwSekBKqRPgykvDpjF4pl2EsmPuekrJgUsNfUsVhUwAxjGVpmPEimJ", "this test")
	if rs != text {
		t.Fatalf("Test fail")
	}
}
