#!/usr/bin/env bash

# Check if the "com.quexten.Goldwarden" Flatpak is installed
if flatpak list | grep -q "com.quexten.Goldwarden"; then
  flatpak run --command=goldwarden com.quexten.Goldwarden "$@"
else
  # If not installed, attempt to run the local version
  goldwarden "$@"
fi
