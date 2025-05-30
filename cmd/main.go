package main

import (
	"blockchain_demo/pkg/script_vm"
	"blockchain_demo/pkg/sign/sign_ed25519"
	"blockchain_demo/pkg/wallet"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type Script struct {
	ScriptSig    string `json:"script_sig" xml:"script_sig"`
	ScriptPubKey string `json:"script_pub_key" xml:"script_pub_key"`
	SignedData   string `json:"signed_data" xml:"signed_data"`
}

func CreateWallet(c *gin.Context) {
	signer := sign_ed25519.Ed25519Signer{}

	keys, err := signer.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to create keys: %v", err),
		})
	}

	prefix := []byte{0x00} // Example prefix for mainnet
	wallet, err := wallet.CreateWallet(keys, prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to create wallet: %v", err),
		})
	}
	publicHash, err := wallet.GetPublicKeyHash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to create hash: %v", err),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "wallet created successfully",
		"data": gin.H{
			"address":         wallet.Address,
			"public_key":      hex.EncodeToString(wallet.Keys.PublicKey),
			"public_key_hash": hex.EncodeToString(publicHash),
			"private_key":     hex.EncodeToString(wallet.Keys.PrivateKey),
		},
	})
}

func runScript (vm *script_vm.VM,  script []byte, signedData []byte) (string, error) {
	err := vm.ParseScript(script)
	if err != nil {
		return "", fmt.Errorf("failed to parse script: %v", err)
	}
	scriptCode := vm.String()
	_, err = vm.Execute(signedData)
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %v", err)
	}

	return scriptCode, nil
}

func ScriptRun(c *gin.Context) {
	script := Script{}
	if err := c.ShouldBind(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to bind script: %v", err),
		})
		return
	}
	scriptSig, err := hex.DecodeString(script.ScriptSig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to decode script_sig: %v", err),
		})
		return
	}
	scriptPubKey, err := hex.DecodeString(script.ScriptPubKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to decode script_pub_key: %v", err),
		})
		return
	}
	signedData, err := hex.DecodeString(script.SignedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to decode signed_data: %v", err),
		})
		return
	}
	fullScript := append(scriptSig, scriptPubKey...)
	signer := sign_ed25519.Ed25519Signer{}
	vm := script_vm.New(&signer)
	scriptCode, err := runScript(vm, fullScript, signedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code": scriptCode,
	})
}

func ScriptCompile(c *gin.Context) {
	script := Script{}
	if err := c.ShouldBind(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to bind script: %v", err),
		})
		return
	}

	signer := sign_ed25519.Ed25519Signer{}
	vm := script_vm.New(&signer)
	scriptSigByteCode, err := vm.ParseString(script.ScriptSig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to parse script: %v", err),
		})
		return
	}

	vm = script_vm.New(&signer)
	scriptPubKeyByteCode, err := vm.ParseString(script.ScriptPubKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to parse script: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"scriptSig": hex.EncodeToString(scriptSigByteCode),
		"scriptPubKey": hex.EncodeToString(scriptPubKeyByteCode),
	})
}

func ScriptParse(c *gin.Context) {
	script := Script{}
	if err := c.ShouldBind(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("failed to bind script: %v", err),
		})
		return
	}

	signer := sign_ed25519.Ed25519Signer{}
	scriptSigCode := ""
	if script.ScriptSig != "" {
		vm := script_vm.New(&signer)

		scriptSig, err := hex.DecodeString(script.ScriptSig)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Errorf("failed to decode script_sig: %v", err),
			})
			return
		}
		err = vm.ParseScript(scriptSig)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Errorf("failed to parse script_sig: %v", err),
			})
			return
		}
		scriptSigCode = vm.String()
	}

	scriptPubKeyCode := ""
	if script.ScriptPubKey != "" {	
		vm := script_vm.New(&signer)
		scriptPubKey, err := hex.DecodeString(script.ScriptPubKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Errorf("failed to decode script_pub_key: %v", err),
			})
			return
		}		
		err = vm.ParseScript(scriptPubKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message":fmt.Errorf("failed to parse PubKey script: %v", err),
			})
			return
		}
		scriptPubKeyCode = vm.String()
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"scriptSig": scriptSigCode,
		"scriptPubKey": scriptPubKeyCode,
	})
}

func main() {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("../public/dist", true)))
	
	api := router.Group("/api")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	api.POST("/wallet", CreateWallet)
	api.POST("/sript/run", ScriptRun)
	api.POST("/sript/compile", ScriptCompile)
	api.POST("/sript/parse", ScriptParse)
	router.Run()
}

