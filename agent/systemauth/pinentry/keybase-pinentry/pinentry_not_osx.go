// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

//go:build !darwin
// +build !darwin

package pinentry

type pinentrySecretStoreInfo struct{}

func (pi *pinentryInstance) useSecretStore(useSecretStore bool) (pinentrySecretStoreInfo, error) {
	return pinentrySecretStoreInfo{}, nil
}

func (pi *pinentryInstance) shouldStoreSecret(info pinentrySecretStoreInfo) bool {
	return false
}
