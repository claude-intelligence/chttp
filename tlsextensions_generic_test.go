package http

import (
	"bytes"
	"testing"

	utls "github.com/refraction-networking/utls"
)

// Chrome 147+ 的 ja3 含 genMap 未实现的 51764, 需要 GenericExtensions 提供原始数据。
func TestStringToSpecGenericExtensions(t *testing.T) {
	const ja3 = "771,4865-4866-4867,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17613-65037-51764-4832,4588-29-23-24,0"
	const ua = "Mozilla/5.0 (Macintosh) Chrome/156.0.0.0 Safari/537.36"

	if _, err := (&TLSExtensions{}).StringToSpec(ja3, ua, false, false); err == nil {
		t.Fatal("未提供原始数据时应报扩展不支持")
	}

	data := []byte{0x00, 0x04, 0xd6, 0x79, 0x09, 0x01}
	ext := &TLSExtensions{GenericExtensions: map[uint16][]byte{51764: data, 4832: {0, 0}}}
	spec, err := ext.StringToSpec(ja3, ua, false, false)
	if err != nil {
		t.Fatal(err)
	}
	found := map[uint16][]byte{}
	for _, e := range spec.Extensions {
		if g, ok := e.(*utls.GenericExtension); ok {
			found[g.Id] = g.Data
		}
	}
	if !bytes.Equal(found[51764], data) || !bytes.Equal(found[4832], []byte{0, 0}) {
		t.Fatalf("GenericExtension 未按原始数据生成: %x", found)
	}
}
