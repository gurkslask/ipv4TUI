package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Octet struct {
	num int
	s   [8]bool
}

func (o *Octet) Setnum(n int) error {
	// Set Octet with a int between 0-255
	// Returns error if not between those numbers
	o.num = n
	err := o.Validate()
	if err != nil {
		return err
	}

	o.s = o.GetboolSlice()
	return nil
}

func (o Octet) Printint() {
	// Prints the decimal number of the octet
	fmt.Println(o.num)
}
func (o Octet) GetString() string {
	// Prints the decimal number as a string
	return fmt.Sprint(o.num)
}

func (o Octet) Printbyte() {
	// Prints the binary numbers as a string
	fmt.Printf("%08b\n", byte(o.num))
}
func (o Octet) Getbyte() string {
	// Returns the binary numbers as a string
	return fmt.Sprintf("%08b", byte(o.num))
}
func (o Octet) GetbyteAsByte() byte {
	// Returns the binary numbers as a string
	return byte(o.num)
}

func (o Octet) Getbitfrombyte(n int) bool {
	// Gets the n bit from the byte
	// 1, 2, 4, 16, 32, 64, 128
	return byte(o.num)&byte(n) > 0
}
func (o *Octet) GetboolSlice() [8]bool {
	// Returns a slice of all the booleans
	var res [8]bool
	j := 1
	for i := 0; i < 7; i++ {
		res[i] = o.Getbitfrombyte(j)
		j = j * 2
	}
	return res
}
func (o Octet) Validate() error {
	if o.num < 0 || o.num > 255 {
		return fmt.Errorf("Invalid number, not between 0-255, number was: %v", o.num)
	} else {
		return nil
	}
}
func octetValidator(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if i > 255 || i < 0 {
		err = fmt.Errorf("Invalid number %v", s)
	}
	return err
}

type IPv4 struct {
	address string
	octets  [4]Octet
}

func NewIPv4FromBinary(b [4]byte) (IPv4, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v.", int(b[0])))
	sb.WriteString(fmt.Sprintf("%v.", int(b[1])))
	sb.WriteString(fmt.Sprintf("%v.", int(b[2])))
	sb.WriteString(fmt.Sprintf("%v", int(b[3])))
	return NewIPv4(sb.String())
}

func NewIPv4(n string) (IPv4, error) {
	var ip IPv4
	err := ip.SetAddress(n)
	if err != nil {
		return ip, err
	}
	return ip, nil
}
func NewIPv4SubnetMask(n string) (IPv4, error) {
	var ip IPv4
	err := ip.SetAddress(n)
	if err != nil {
		return ip, err
	}
	err = ip.CheckIfValidSubnetmask()
	if err != nil {
		return ip, err
	}
	return ip, nil
}

func (ip *IPv4) SetAddress(in string) error {
	//fmt.Println(in)
	err := ip.CheckIfValidIPAddress(in)
	if err != nil {
		return err
	}
	splittedoctets := strings.Split(in, ".")
	for i := 0; i < 4; i++ {
		err := ip.octets[i].Setnum(StrToInt(splittedoctets[i]))
		if err != nil {
			return err
		}
	}
	ip.address = in
	return nil
}
func (ip IPv4) PrintDecimal() string {
	return ip.address
}
func (ip IPv4) PrintBinary() string {
	var sb strings.Builder
	for i := 0; i < 4; i++ {
		sb.WriteString(ip.octets[i].Getbyte())
	}
	return sb.String()
}

// Returns [4]byte of the ip adress in 4 bytes
func (ip IPv4) Getbytes() [4]byte {
	var sb [4]byte
	for i := range 4 {
		sb[i] = ip.octets[i].GetbyteAsByte()
	}
	return sb
}

// Returns [32]bool of all the bits
func (ip IPv4) GetBits() [32]bool {
	var res [32]bool
	b := ip.Getbytes()
	for octetnum, octetbyte := range b {
		bb := GetBoolsFromByte(octetbyte)
		for bitnum, bit := range bb {
			res[bitnum+octetnum*8] = bit
		}
	}
	return res
}
func (ip IPv4) CheckIfValidIPAddress(in string) error {
	splittedoctets := strings.Split(in, ".")
	if len(splittedoctets) != 4 {
		return fmt.Errorf("Not a valid IP %v", in)
	}
	for i := 0; i < 4; i++ {
		err := ip.octets[i].Setnum(StrToInt(splittedoctets[i]))
		if err != nil {
			return err
		}
	}
	return nil
}
func (ip IPv4) CheckIfValidSubnetmask() error {
	bs := ip.PrintBinary()
	//fmt.Println("SUBNET BINARY", bs)
	firstzero := false
	for k, v := range bs {
		if !firstzero && v == '1' {
			// We havent seen a zero, and the actual rune is 1, do nothing
			continue
		} else if firstzero && v == '1' {
			//fmt.Println("SUBNET IP ADRESS", ip.address)
			return fmt.Errorf("Not a valid subnetmask %x %#U, KEY %v", v, v, k)
		} else if !firstzero && v == '0' {
			firstzero = true
		}
	}
	return nil
}

// Calculates a netaddress based on a ip and subnetmask, returns IPv4 struct
func CalcNetAddress(ip IPv4, subnetmask IPv4) IPv4 {
	ipb := ip.Getbytes()
	smb := subnetmask.Getbytes()
	var bb [4]byte
	for k, _ := range ipb {
		bb[k] = ipb[k] & smb[k]
	}
	resIP, _ := NewIPv4FromBinary(bb)
	return resIP
}

// Calculates a broadcastaddress based on a ip and subnetmask, returns IPv4 struct
func CalcBroadcastAddress(ip IPv4, subnetmask IPv4) IPv4 {
	cidr := CalcCIDR(subnetmask)
	ipb := ip.GetBits()
	smb := subnetmask.GetBits()
	var res [32]bool
	var bb [4]byte
	for i := range 32 {
		if i < cidr {
			res[i] = ipb[i] && smb[i]
		} else {
			res[i] = ipb[i] || !smb[i]
		}
	}
	for i := range 4 {
		bb[i] = GetByteFromBools(res[i*8 : i*8+8])
	}
	resIP, _ := NewIPv4FromBinary(bb)
	return resIP
}

// Calculates a CIDR address from a subnet mask. Example: 255.255.255.0 = 24
func CalcCIDR(subnetmask IPv4) int {
	res := 0
	for _, v := range subnetmask.Getbytes() {
		bb := GetBoolsFromByte(v)
		for _, bit := range bb {
			if bit {
				res = res + 1
			} else {
				return res
			}
		}
	}
	return res
}

// This function calculates the amount of host addresses and net addresses
// with a given cidr (subnetmask). The function excludes Net-address and Broadcast-address
// Returns netaddress, host address
func CalcCombinations(cidr int) (int, int) {
	netaddresses := 0
	for i := range cidr {
		netaddresses += int(math.Pow(2, float64(i)))
	}
	netaddresses -= 1
	hostaddresses := 0
	for i := range 32 - cidr {
		hostaddresses += int(math.Pow(2, float64(i)))
	}
	hostaddresses -= 1
	return netaddresses, hostaddresses

}
