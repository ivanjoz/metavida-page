package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"strings"

	"github.com/fxamacker/cbor/v2"
)

type UsuarioToken struct {
	CompanyID int32  `json:"c,omitempty" cbor:"1,keyasint"`
	ID        int32  `json:"i" cbor:"2,keyasint"`
	Created   int32  `json:"e" cbor:"3,keyasint"`
	Hash      uint64 `json:"h" cbor:"4,keyasint"`
	User      string `json:"u" cbor:"5,keyasint"`
	Role      string `json:"r,omitempty" cbor:"6,keyasint"`
	Error     string `json:"-" cbor:"-"`
}


var authSecret = "metavida-local-secret"

func SetAuthSecret(secret string) {
	secret = strings.TrimSpace(secret)
	if secret != "" {
		authSecret = secret
	}
}

func GetAuthSecret() string { return authSecret }

func ComputeUsuarioTokenHash(usuarioToken UsuarioToken) uint64 {
	hashMac := hmac.New(sha256.New, []byte(authSecret))
	tokenPayloadBuffer := make([]byte, 8)
	binary.BigEndian.PutUint32(tokenPayloadBuffer[0:4], uint32(usuarioToken.CompanyID))
	binary.BigEndian.PutUint32(tokenPayloadBuffer[4:8], uint32(usuarioToken.Created))
	hashMac.Write([]byte("usrToken:v1"))
	hashMac.Write(tokenPayloadBuffer)
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, uint32(usuarioToken.ID))
	hashMac.Write(idBytes)
	hashMac.Write([]byte(usuarioToken.User))
	hashMac.Write([]byte(usuarioToken.Role))
	return binary.BigEndian.Uint64(hashMac.Sum(nil)[:8])
}

func EncodeUsuarioToken(usuarioToken UsuarioToken) (string, error) {
	usuarioToken.Hash = ComputeUsuarioTokenHash(usuarioToken)
	tokenBytes, err := cbor.Marshal(usuarioToken)
	if err != nil {
		return "", err
	}
	return BytesToBase64(tokenBytes, true), nil
}

func DecodeUsuarioToken(tokenBase64 string) (*UsuarioToken, error) {
	tokenBytes, err := Base64ToBytes(MakeB64UrlDecode(strings.TrimSpace(tokenBase64)))
	if err != nil {
		return nil, err
	}
	user := UsuarioToken{}
	if err := cbor.Unmarshal(tokenBytes, &user); err != nil {
		return nil, err
	}
	if user.Hash != ComputeUsuarioTokenHash(user) {
		user.Error = "El token de inicio de sesión no es válido."
	}
	return &user, nil
}

func DecodeUserToken(tokenBase64 string) (*UsuarioToken, error) {
	return DecodeUsuarioToken(tokenBase64)
}

func CheckUser(req *HandlerArgs, _ int) *UsuarioToken {
	userToken := req.Headers["authorization"]
	if len(userToken) < 8 {
		userToken = req.Headers["Authorization"]
	}

	user := UsuarioToken{}
	if len(userToken) < 8 {
		user.Error = "No se suministró un Token de user"
		return &user
	}

	tokenBase64 := strings.TrimSpace(strings.TrimPrefix(userToken, "Bearer "))
	if len(tokenBase64) < 8 {
		user.Error = "No se encontró la información del user."
		return &user
	}

	decodedUser, err := DecodeUsuarioToken(tokenBase64)
	if err != nil {
		user.Error = "Error al recuperar la información del user."
		return &user
	}
	return decodedUser
}

func GetUserFromAuthorization(authorization string) *UsuarioToken {
	args := HandlerArgs{Headers: map[string]string{"Authorization": authorization}}
	return CheckUser(&args, 0)
}

