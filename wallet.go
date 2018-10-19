package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"base58"
)

//创建钱包结构
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey []byte
}
//创建生成新的钱包的方法
func NewWallet()*Wallet  {
	//创建曲线
	curve := elliptic.P256()
	//生成私钥
	privatekey, err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil{
		log.Panic(err)
	}
	//生成公钥,公钥是一对坐标点，将坐标点进行编码拼接
	publickeyOrig := privatekey.PublicKey
	publikey := append(publickeyOrig.X.Bytes(),publickeyOrig.Y.Bytes()...)
	//返回新生成的钱包
	return &Wallet{PrivateKey:privatekey,PublicKey:publikey}
}
//为钱包绑定生成地址的方法
func (w *Wallet)NewAddress()string  {
	//先从钱包中取得公钥
	pubkey := w.PublicKey
	//对公钥进行sha256和160编码
	rip160hash := ShaRipHash(pubkey)
	//进行字节的拼接，将版本号和进行hash运算的hash值进行拼接
	var version = []byte("00")
	payload := append(version,rip160hash...)
	//将拼接好的payload复制一份，进行两次的sha256的hash运算
	checkCode := CheckSum(payload)
	//再将取得的前四字节的hash值和原来的21字节的hash值进行拼接
	payload = append(payload,checkCode...)
	//对新生成的hash值进行sha256编码,然后返回地址
	address := base58.Encode(payload)
	return address

}
func ShaRipHash(pubkey []byte)[]byte  {
	//先对公钥进行sha256编码
	shaHash := sha256.Sum256(pubkey)
	//在对sha256生成的hash进行rip160bianma
	//先创建rip160编码器
	rip160Shaer := ripemd160.New()
	_, err := rip160Shaer.Write(shaHash[:])
	if err != nil{
		log.Panic(err)
	}
	rip160Hash := rip160Shaer.Sum(nil)
	//返回经过两次编码生成的hash值
	return rip160Hash[:]

}
func CheckSum(payload []byte)[]byte  {
	//进行第一次的sha256运算
	hash1 := sha256.Sum256(payload)
	//进行第二次的sha256运算
	hash2 := sha256.Sum256(hash1[:])
	//取进行两次sha256的hash值的前四个字节
	checkCode := hash2[:4]
	return checkCode

}































