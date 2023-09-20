//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Development target
func Dev() error {
	return sh.Run("go", "build", "-trimpath", "-ldflags="+"-s -w", "-tags", "dev")
}

// Production target for local directory
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-trimpath", "-ldflags="+"-s -w")
}

// Production target install
func Install() error {
	return sh.Run("go", "install", "-trimpath", "-ldflags="+"-s -w", "./...", "-tags"+"prod")
}

func Clean() error {
	return nil
}
