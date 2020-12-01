package blockchain

import (
	caMsp "github.com/off-grid-block/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
)

// InvokeHello
func Register(s *SetupSDK, data caMsp.RegistrationRequest) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	// new User information

	caClient, err := caMsp.New(s.Fsdk.Context())
	if err != nil {
		return "", errors.WithMessage(err, "Failed to create new msp client")
	}
	enrollSecret, err := caClient.Register(&data)
	if err != nil {
		return "", errors.WithMessage(err, "Unable to register user with CA")
	}

	return string(enrollSecret), nil
}
