package iplib

import (
	"testing"

	common "github.com/gurkslask/common-go-libs/common"
)

func TestOctetType(t *testing.T) {
	var o Octet
	var ans error

	ans = o.Setnum(42)
	if ans != nil {
		t.Errorf("Should work")
	}
	ans = o.Setnum(442)
	if ans == nil {
		t.Errorf("Shouldnt work")
	}
	ans = o.Setnum(-442)
	if ans == nil {
		t.Errorf("Shouldnt work")
	}
	t.Run("Getstring", func(t *testing.T) {
		o := Octet{num: 42}
		ans := o.GetString()
		if ans != "42" {
			t.Errorf("Should be 42")
		}
	})
	t.Run("Getbyte", func(t *testing.T) {
		o := Octet{num: 1}
		ans := o.Getbyte()
		if ans != "00000001" {
			t.Errorf("Should be 42")
		}
	})
	t.Run("Getbit", func(t *testing.T) {
		o := Octet{num: 1}
		ans := o.Getbitfrombyte(1)
		if ans != true {
			t.Errorf("Should be 42")
		}
	})
}

func TestIPv4Addr(t *testing.T) {
	t.Run("Test create IP", func(t *testing.T) {
		addr := "192.168.1.1"
		_, err := NewIPv4(addr)
		if err != nil {
			t.Errorf("%v doesnt work", addr)
		}
	})
	t.Run("Test create IP too few", func(t *testing.T) {
		addr := "123.132.123"
		_, err := NewIPv4(addr)
		if err == nil {
			t.Errorf("%v doesnt work", addr)
		}
	})
	t.Run("Test create IP2 only text", func(t *testing.T) {
		addr := "hej"
		_, err := NewIPv4(addr)
		if err == nil {
			t.Errorf("%v doesnt work", addr)
		}
	})
	t.Run("Test create too large", func(t *testing.T) {
		addr := "99999.999999.999999.99999"
		_, err := NewIPv4(addr)
		if err == nil {
			t.Errorf("%v doesnt work", addr)
		}
	})
	t.Run("Test create byte", func(t *testing.T) {
		addr := [4]byte{11, 111, 123, 54}
		_, err := NewIPv4FromBinary(addr)
		if err != nil {
			t.Errorf("%v doesnt work because %v", addr, err)
		}
	})
	t.Run("Test subnetmask", func(t *testing.T) {
		addr := "255.255.255.0"
		_, err := NewIPv4SubnetMask(addr)
		if err != nil {
			t.Errorf("%v doesnt work because %v", addr, err)
		}
	})
	t.Run("Test subnetmask", func(t *testing.T) {
		addr := "128.255.255.0"
		_, err := NewIPv4SubnetMask(addr)
		if err == nil {
			t.Errorf("%v doesnt work because %v", addr, err)
		}
	})

}

func TestGetByteFromBools(t *testing.T) {
	t.Run("Test bools", func(t *testing.T) {
		bnum := byte(255)
		ans := []bool{true, true, true, true, true, true, true, true}
		res := common.GetByteFromBools(ans)
		if res != bnum {
			t.Errorf("This byte: %08b should be %v, got: %v\n", bnum, ans, res)
		}
	})
	t.Run("Test bools 192", func(t *testing.T) {
		bnum := byte(192)
		ans := []bool{true, true, false, false, false, false, false, false}
		res := common.GetByteFromBools(ans)
		if res != bnum {
			t.Errorf("The bits are: %v, we want number %v, we got %v", ans, bnum, res)
		}
	})
	t.Run("Test bools 192 too long", func(t *testing.T) {
		bnum := byte(192)
		ans := []bool{true, true, false, false, false, false, false, false, false, true}
		res := common.GetByteFromBools(ans)
		if res != bnum {
			t.Errorf("The bits are: %v, we want number %v, we got %v", ans, bnum, res)
		}
	})
}
func TestGetBoolFromBytes(t *testing.T) {
	t.Run("Test bools", func(t *testing.T) {
		bnum := byte(255)
		ans := [8]bool{true, true, true, true, true, true, true, true}
		res := common.GetBoolsFromByte(bnum)
		if res != ans {
			t.Errorf("This byte: %08b should be %v, got: %v\n", bnum, ans, res)
		}
	})
	t.Run("Test bools 192", func(t *testing.T) {
		bnum := byte(192)
		ans := [8]bool{true, true, false, false, false, false, false, false}
		res := common.GetBoolsFromByte(bnum)
		if res != ans {
			t.Errorf("This byte: %08b should be %v, got: %v\n", bnum, ans, res)
		}
	})

}
func TestNetAddress(t *testing.T) {
	t.Run("Test 192.168.1.1/24", func(t *testing.T) {
		ipaddr, _ := NewIPv4("192.168.1.1")
		subnetmask, _ := NewIPv4SubnetMask("255.255.255.0")
		ans, _ := NewIPv4("192.168.1.0")
		res := CalcNetAddress(ipaddr, subnetmask)
		if res.address != ans.address {
			t.Errorf("This ip addr %v and this subnetmask %v \n Netaddr should be %v got: %v \n", ipaddr.address, subnetmask.address, ans.address, res.address)
		}
	})
	t.Run("Test 10.1.3.4/16", func(t *testing.T) {
		ipaddr, _ := NewIPv4("10.1.3.4")
		subnetmask, _ := NewIPv4SubnetMask("255.255.0.0")
		ans, _ := NewIPv4("10.1.0.0")
		res := CalcNetAddress(ipaddr, subnetmask)
		if res.address != ans.address {
			t.Errorf("This ip addr %v and this subnetmask %v \n Netaddr should be %v got: %v \n", ipaddr.address, subnetmask.address, ans.address, res.address)
		}
	})
	t.Run("Test 10.1.3.4/8", func(t *testing.T) {
		ipaddr, _ := NewIPv4("10.1.3.4")
		subnetmask, _ := NewIPv4SubnetMask("255.0.0.0")
		ans, _ := NewIPv4("10.0.0.0")
		res := CalcNetAddress(ipaddr, subnetmask)
		if res.address != ans.address {
			t.Errorf("This ip addr %v and this subnetmask %v \n Netaddr should be %v got: %v \n", ipaddr.address, subnetmask.address, ans.address, res.address)
		}
	})
}
func TestBroadcastAddress(t *testing.T) {
	t.Run("Test 192.168.1.1/24", func(t *testing.T) {
		ipaddr, _ := NewIPv4("192.168.1.1")
		subnetmask, _ := NewIPv4SubnetMask("255.255.255.0")
		ans, _ := NewIPv4("192.168.1.255")
		res := CalcBroadcastAddress(ipaddr, subnetmask)
		if res.address != ans.address {
			t.Errorf("This ip addr %v and this subnetmask %v \n Broadcast should be %v got: %v \n", ipaddr.address, subnetmask.address, ans.address, res.address)
		}
	})
	t.Run("Test 10.1.2.3/8", func(t *testing.T) {
		ipaddr, _ := NewIPv4("10.1.2.3")
		subnetmask, _ := NewIPv4SubnetMask("255.0.0.0")
		ans, _ := NewIPv4("10.255.255.255")
		res := CalcBroadcastAddress(ipaddr, subnetmask)
		if res.address != ans.address {
			t.Errorf("This ip addr %v and this subnetmask %v \n Broadcast should be %v got: %v \n", ipaddr.address, subnetmask.address, ans.address, res.address)
		}
	})
	t.Run("Test 10.1.2.3/9", func(t *testing.T) {
		ipaddr, _ := NewIPv4("10.1.2.3")
		subnetmask, _ := NewIPv4SubnetMask("255.128.0.0")
		ans, _ := NewIPv4("10.127.255.255")
		res := CalcBroadcastAddress(ipaddr, subnetmask)
		if res.address != ans.address {
			t.Errorf("This ip addr %v and this subnetmask %v \n Broadcast should be %v got: %v \n", ipaddr.address, subnetmask.address, ans.address, res.address)
		}
	})
}
func TestCIDR(t *testing.T) {
	t.Run("Test 24", func(t *testing.T) {
		subnetmask, _ := NewIPv4SubnetMask("255.255.255.0")
		ans := 24
		res := CalcCIDR(subnetmask)
		if res != ans {
			t.Errorf("This subnetmask %v, should have %v, got %v", subnetmask.address, ans, res)
		}
	})
	t.Run("Test 8", func(t *testing.T) {
		subnetmask, _ := NewIPv4SubnetMask("255.0.0.0")
		ans := 8
		res := CalcCIDR(subnetmask)
		if res != ans {
			t.Errorf("This subnetmask %v, should have %v, got %v", subnetmask.address, ans, res)
		}
	})
}
func TestGetBitsFromIp(t *testing.T) {
	t.Run("Test 192.192.0.1", func(t *testing.T) {
		ip, _ := NewIPv4("192.192.0.1")
		ans := [32]bool{true, true, false, false, false, false, false, false,
			true, true, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, true}
		res := ip.GetBits()
		if res != ans {
			t.Errorf("This ip %v, should have %v, got %v", ip.address, ans, res)
		}
	})
}
func TestNetHostAddresses(t *testing.T) {
	t.Run("Test 24", func(t *testing.T) {
		subnetmask := 24
		ansNet := 16777214
		ansHost := 254
		resNet, resHost := CalcCombinations(subnetmask)
		if resNet != ansNet {
			t.Errorf("This subnetmask %v, should have %v Host addresses, got %v", subnetmask, ansNet, resNet)
		}
		if resHost != ansHost {
			t.Errorf("This subnetmask %v, should have %v Host addresses, got %v", subnetmask, ansHost, resHost)
		}
	})
	t.Run("Test 16", func(t *testing.T) {
		subnetmask := 16
		ansNet := 65534
		ansHost := 65534
		resNet, resHost := CalcCombinations(subnetmask)
		if resNet != ansNet {
			t.Errorf("This subnetmask %v, should have %v Host addresses, got %v", subnetmask, ansNet, resNet)
		}
		if resHost != ansHost {
			t.Errorf("This subnetmask %v, should have %v Host addresses, got %v", subnetmask, ansHost, resHost)
		}
	})
}
