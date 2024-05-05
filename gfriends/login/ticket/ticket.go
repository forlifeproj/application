package ticket

import (
	"encoding/base64"
	"fmt"

	conf "github.com/forlifeproj/application/gfriends/login/config"
	"github.com/forlifeproj/application/gfriends/login/util"
	fllog "github.com/forlifeproj/msf/log"
)

/*
 *  ticket 格式为:
 * --------------------------------------------------------------------------------------------------------------------------------
 * | version | login type | create_time | appid | sercret len | (crc32(appid_openid) | create_time | random) |
 * --------------------------------------------------------------------------------------------------------------------------------
 * |  4字节  |    4字节    |    4字节    | 4字节 |    4字节    |          (4字节   		 |    4字节    |    4字节)    |
 * --------------------------------------------------------------------------------------------------------------------------------
 */

 const (
	TICKET_MAX_LEN = 128
	PLAIN_TEST_LEN  20
	ENCRYPT_FIELD_NUM = 3
 )

type Ticket struct {
	Version      int
	LoginType    int
	CreateTime   int
	Appid        int
	Openid       string
	Random       int
	SecrectLen   int
	SecrectBytes []byte
	CreateTicket string
}

func (t *Ticket) SetVersion(version int) {
	t.Version = version
}

func (t *Ticket) SetLoginType(loginType int) {
	t.LoginType = loginType
}

func (t *Ticket) SetCreateTime(ct int) {
	t.CreateTime = ct
}

func (t *Ticket) SetAppid(appid int) {
	t.Appid = appid
}

func (t *Ticket) SetOpenid(openid string) {
	t.Openid = openid
}

func (t *Ticket) SetSecretLen(len int) {
	t.SecrectLen = len
}

func (t *Ticket) SetRandom(random int) {
	t.Random = random
}

func (t *Ticket) SetTicket(ticket string) {
	t.CreateTicket = ticket
}

func (t *Ticket) CreateTicket() string {
	if err := t.EncryptSecret(); err != nil {
		fllog.Log().Error(fmt.Sprintf("encrypt secret failed."))
		return ""
	}
	intArray := []int{
		t.Version,
		t.LoginType,
		t.CreateTime,
		t.Appid,
		len(t.SecrectBytes),
	}
	ticketBytes := util.IntArray2Bytes(intArray)
	ticketBytes = append(ticketBytes, t.SecrectBytes)
	t.CreateTicket := base64.StdEncoding.EncodeToString(bytticketByteseArray)
	return t.CreateTicket
}

func (t *Ticket) IsValidTicket() bool {
	if !t.ParseTicket() {
		fllog.Log().Error(fmt.Sprintf("parse ticket failed. ticket=[%s]", t.CreateTicket))
		return false
	}

	if !t.DecryptSecret() {
		fllog.Log().Error(fmt.Sprintf("decrypt secret failed. ticket=[%s]", t.CreateTicket))
		return false
	}
	
	return true
}

func (t *Ticket) ParseTicket() bool {
	if len(t.CreateTicket) == 0 || len(t.CreateTicket) > TICKET_MAX_LEN {
		fllog.Log().Error(fmt.Sprintf("invalid ticket len:%d ticke:%s", len(t.CreateTicket), t.CreateTicket))
		return false
	}

	byteArray, err := base64.StdEncoding.DecodeString(t.CreateTicket)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("Base64 Decode failed. err:%+v ticket:%s", err, t.CreateTicket))
		return false
	}

	plaintextBytes := make([]byte, 4 * 5)
	copy(plaintextBytes, byteArray[:PLAIN_TEST_LEN])
	
	plaintextInts := util.Bytes2IntArray(plaintextBytes)
	if len(plaintextInts) != 5 {
		fllog.Log().Error(fmt.Sprintf("plaintext bytes->intarray failed. plaintextBytes:%+v, plaintextInts:%+v", 
		plaintextBytes, plaintextInts))
		return false
	}
	t.SetVersion(plaintextInts[0])
	t.SetAppid(plaintextInts[1])
	t.SetLoginType(plaintextInts[2])
	t.SetCreateTime(plaintextInts[3])
	t.SetSecretLen(plaintextInts[4])	

	t.SecrectBytes = make([]byte, t.SecrectLen)
	copy(t.SecrectBytes, byteArray[PLAIN_TEST_LEN:])
	
	return true
}

// 加密secret部分
func (t *Ticket) EncryptSecret() error {
	intArray := []int{
		util.GetCrc32(fmt.Sprintf("%d_%s", t.Appid, t.Openid)),
		t.CreateTime,
		t.Random,
	}
	intBytes := util.IntArray2Bytes(intArray)

	secrectBytes, err := util.Encrypt(intBytes, conf.GConf.TicketSecret.SecretKey)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("encrypt intBytes:%s failed err:%+v", string(intBytes), err))
		return err
	}
	t.SecrectBytes = secrectBytes
	return nil
}

// 解密secret部分
func (t *Ticket) DecryptSecret() bool {
	decryptBytes , err := util.Decrypt(t.SecrectBytes, conf.GConf.TicketSecret.SecretKey)
	if err != nil {
		fllog.Log().Error(fmt.Sprintf("decrypt secretbytes failed. err:%+v",err))
		return false
	}
	intArray := util.Bytes2IntArray(decryptBytes)
	if len(intArray) != ENCRYPT_FIELD_NUM {
		fllog.Log().Error(fmt.Sprintf("invalid intArray:%+v",intArray))
		return false
	}

	if len(t.Openid) > 0 {
		curCrc32 := util.GetCrc32(fmt.Sprintf("%d_%s", t.Appid, t.Openid))
		if intArray[0] != curCrc32 {
			fllog.Log().Error("appid=%d_openid=%s curCrc32=%d not equal intArray[0]=%d", 
				t.Appid, t.Openid, curCrc32, intArray[0])
			return false
		}
	}
	if t.CreateTime != intArray[1] {
		fllog.Log().Error(fmt.Sprintf("secret createTime=%d not equal intArray[1]=%d", 
			t.CreateTime, intArray[1]))
		return false
	}
	return false
}
