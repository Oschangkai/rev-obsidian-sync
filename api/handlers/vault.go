package handlers

import (
	"github.com/acheong08/obsidian-sync/database"
	"github.com/acheong08/obsidian-sync/utilities"
	"github.com/acheong08/obsidian-sync/vault"
	"github.com/gin-gonic/gin"
	password_generator "github.com/sethvargo/go-password/password"
)

func ListVaults(c *gin.Context) {
	type request struct {
		Token string `json:"token" binding:"required"`
	}
	type response struct {
		Shared []any          `json:"shared"`
		Vaults []*vault.Vault `json:"vaults"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	dbConnection := c.MustGet("db").(*database.Database)
	vaults, err := dbConnection.GetVaults(email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response{
		Shared: []any{},
		Vaults: vaults,
	})

}

func CreateVault(c *gin.Context) {
	type request struct {
		KeyHash *string `json:"keyhash"`
		Name    string  `json:"name" binding:"required"`
		Salt    *string `json:"salt"`
		Token   string  `json:"token" binding:"required"`
	}
	// Response is vault details
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	dbConnection := c.MustGet("db").(*database.Database)
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	var password string
	var salt string
	var keyHash string
	// Generate password if keyhash is not provided
	if req.Salt == nil {
		password, err = password_generator.Generate(20, 5, 5, false, true)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		salt, err = password_generator.Generate(20, 5, 5, false, true)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		keyHash = ""
	} else {
		if req.KeyHash != nil {
			keyHash = *req.KeyHash
		} else {
			c.JSON(400, gin.H{"error": "keyhash must be provided if salt is provided"})
		}
	}
	vault, err := dbConnection.NewVault(req.Name, email, password, salt, keyHash)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, vault)

}

func DeleteVault(c *gin.Context) {
	type request struct {
		Token    string `json:"token" binding:"required"`
		VaultUID string `json:"vault_uid" binding:"required"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	dbConnection := c.MustGet("db").(*database.Database)
	email, err := utilities.GetJwtEmail(req.Token)
	if err != nil {
		// Unauthorized
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	err = dbConnection.DeleteVault(req.VaultUID, email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{})
}