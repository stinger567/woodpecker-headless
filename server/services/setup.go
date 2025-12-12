// Copyright 2024 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"go.woodpecker-ci.org/woodpecker/v3/server/services/config"
	"go.woodpecker-ci.org/woodpecker/v3/server/services/registry"
	"go.woodpecker-ci.org/woodpecker/v3/server/services/secret"
	"go.woodpecker-ci.org/woodpecker/v3/server/services/utils"
	"go.woodpecker-ci.org/woodpecker/v3/server/store"
	"go.woodpecker-ci.org/woodpecker/v3/server/store/types"
)

func setupRegistryService(store store.Store, dockerConfig string) registry.Service {
	if dockerConfig != "" {
		return registry.NewCombined(
			registry.NewDB(store),
			registry.NewFilesystem(dockerConfig),
		)
	}

	return registry.NewDB(store)
}

func setupSecretService(store store.Store) secret.Service {
	// TODO(1544): fix encrypted store
	// // encryption
	// encryptedSecretStore := encryptedStore.NewSecretStore(v)
	// err := encryption.Encryption(c, v).WithClient(encryptedSecretStore).Build()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("could not create encryption service")
	// }

	return secret.NewDB(store)
}

func setupConfigService(c *cli.Command, client *utils.Client) (config.Service, error) {
	timeout := c.Duration("forge-timeout")
	retries := c.Uint("forge-retry")
	if retries == 0 {
		return nil, fmt.Errorf("WOODPECKER_FORGE_RETRY can not be 0")
	}
	configFetcher := config.NewForge(timeout, retries)

	if endpoint := c.String("config-service-endpoint"); endpoint != "" {
		httpFetcher := config.NewHTTP(endpoint, client)
		return config.NewCombined(configFetcher, httpFetcher), nil
	}

	return configFetcher, nil
}

// setupSignatureKeys generate or load key pair to sign webhooks requests (i.e. used for service extensions).
func setupSignatureKeys(_store store.Store) (ed25519.PrivateKey, crypto.PublicKey, error) {
	privKeyID := "signature-private-key"

	privKey, err := _store.ServerConfigGet(privKeyID)
	if errors.Is(err, types.RecordNotExist) {
		_, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
		}
		err = _store.ServerConfigSet(privKeyID, hex.EncodeToString(privKey))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to store private key: %w", err)
		}
		log.Debug().Msg("created private key")
		return privKey, privKey.Public(), nil
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to load private key: %w", err)
	}
	privKeyStr, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode private key: %w", err)
	}
	privateKey := ed25519.PrivateKey(privKeyStr)
	return privateKey, privateKey.Public(), nil
}
