package main

import (
	"blockchain_demo/pkg/script_vm"
	"blockchain_demo/pkg/sign/sign_ed25519"
	"blockchain_demo/pkg/wallet"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)
	
func CreateWallet(c *gin.Context) {
	signer := sign_ed25519.Ed25519Signer{}

	keys, err := signer.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to create keys: %v", err),
    })
	}

	prefix := []byte{0x00} // Example prefix for mainnet
	wallet, err := wallet.CreateWallet(keys, prefix)
	if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to create wallet: %v", err),
    })
	}
  publicHash, err := wallet.GetPublicKeyHash()
	if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to create hash: %v", err),
    })
	}
  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "wallet created successfully",
    "data": gin.H{
      "address": wallet.Address,
      "public_key": hex.EncodeToString(wallet.Keys.PublicKey),
      "public_key_hash": hex.EncodeToString(publicHash),
      "private_key": hex.EncodeToString(wallet.Keys.PrivateKey),
    }, 
  })
}
type Script struct {
  ScriptSig string `json:"script_sig" xml:"script_sig" binding:"required"`
  ScriptPubKey string `json:"script_pub_key" xml:"script_pub_key" binding:"required"`
  SignedData string `json:"signed_data" xml:"signed_data" binding:"required"`
}

func RunScript(c *gin.Context) {
  script := Script{}
  if err := c.ShouldBind(&script); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to bind script: %v", err),
    })
  }
  scriptSig, err := hex.DecodeString(script.ScriptSig)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to decode script_sig: %v", err),
    })
  }
  scriptPubKey, err := hex.DecodeString(script.ScriptPubKey) 
  if err != nil { 
    c.JSON(http.StatusBadRequest, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to decode script_pub_key: %v", err),
    })
  }
  signedData, err := hex.DecodeString(script.SignedData)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to decode signed_data: %v", err),
    })
  }
  fullScript := append(scriptSig, scriptPubKey...)
	signer := sign_ed25519.Ed25519Signer{}
	vm := script_vm.New(&signer)
	err = vm.ParseScript(fullScript)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to parse script: %v", err),
    })
  }
  scriptCode := vm.String()
  err = vm.Execute(signedData, nil)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "success":  false,
      "message": fmt.Sprintf("failed to parse script: %v", err),
    })
  }
  c.String(http.StatusOK, scriptCode)
}

func main() {
  router := gin.Default()
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  api := router.Group("/api")  
  api.POST("/wallet", CreateWallet)
  api.POST("/run_sript", RunScript)
  router.Run()
}
