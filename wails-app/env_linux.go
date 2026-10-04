//go:build linux

package main

import "os"

// WebKitGTK's DMA-BUF renderer fails on the proprietary NVIDIA driver ("Failed to create GBM
// buffer") and can leave the window blank. Disable it unless the user already chose a value.
func init() {
	if _, set := os.LookupEnv("WEBKIT_DISABLE_DMABUF_RENDERER"); !set {
		os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
	}
}
