package cluster

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/stefanuzzo/internal/configuration"
	"github.com/stefanuzzo/utilities"
)

const nodeIdFilename = "nodeid"

func GetOrSetNodeId(configuration *configuration.Configuration) (uuid.UUID, bool, error) {
	nodeIdFilepath := utilities.BuildPath(configuration.RunDirectory, nodeIdFilename)

	nodeId, err := readNodeId(nodeIdFilepath)
	if err != nil {
		return uuid.Nil, false, err
	}

	generated := false

	if nodeId == uuid.Nil {
		_, err = utilities.DeleteFile(nodeIdFilepath)
		if err != nil {
			return uuid.Nil, false, err
		}

		nodeId, err = generateUuid(configuration)
		if err != nil {
			return uuid.Nil, false, err
		}

		generated = true

		err = os.WriteFile(nodeIdFilepath, []byte(nodeId.String()), 0)
		if err != nil {
			return uuid.Nil, false, err
		}
	}

	return nodeId, generated, nil
}

func readNodeId(nodeIdFilepath string) (uuid.UUID, error) {
	exists, err := utilities.FileExists(nodeIdFilepath)
	if err != nil {
		return uuid.Nil, err
	}

	if !exists {
		return uuid.Nil, nil
	}

	nodeIdBytes, err := os.ReadFile(nodeIdFilepath)
	if err != nil {
		return uuid.Nil, err
	}

	nodeIdString := string(nodeIdBytes[:])

	nodeId, err := uuid.Parse(nodeIdString)
	if err != nil {
		return uuid.Nil, err
	}

	return nodeId, nil
}

func generateUuid(configuration *configuration.Configuration) (uuid.UUID, error) {
	switch configuration.NodeIdUuidVersion {
	case 1:
		return uuid.NewUUID()

	case 3:
		return uuid.NewMD5(uuid.NameSpaceOID, []byte(configuration.NodeIdUuidString)), nil

	case 4:
		return uuid.NewRandom()

	case 5:
		return uuid.NewSHA1(uuid.NameSpaceOID, []byte(configuration.NodeIdUuidString)), nil

	default:
		return uuid.Nil, fmt.Errorf("unhandled uuid version: %d", configuration.NodeIdUuidVersion)
	}
}
