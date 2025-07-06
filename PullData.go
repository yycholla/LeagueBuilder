package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetRemoteVersion() (string, error) {
	resp, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		return "", err // Handle error appropriately, e.g., log it
	}
	defer resp.Body.Close()

	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", err // Handle error appropriately, e.g., log it
	}
	if len(versions) == 0 {
		return "", io.EOF // No versions found
	}
	return versions[0], nil // Return the latest version
}

func GetLocalVersion() (string, error) {
	data, err := os.ReadFile("data/version.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // Version file does not exist, return empty string
		}
		return "", err // Handle other errors
	}
	return strings.TrimSpace(string(data)), nil // Return the version from the file
}

func SetLocalVersion(version string) error {
	return ioutil.WriteFile("data/version.txt", []byte(version), 0644) // Write the version to the file
}

func UpdateAvailable() (bool, error) {
	remoteVersion, err := GetRemoteVersion()
	if err != nil {
		return false, err // Handle error appropriately, e.g., log it
	}
	localVersion, err := GetLocalVersion()
	if err != nil {
		return false, err // Handle error appropriately, e.g., log it
	}
	if localVersion == "" {
		return true, SetLocalVersion(remoteVersion) // If local version is empty, set it to remote version
	}
	if remoteVersion != localVersion {
		if err := SetLocalVersion(remoteVersion); err != nil {
			return false, err // Handle error appropriately, e.g., log it
		}
		return true, nil // Return true indicating an update was made
	}
	return false, nil // No update needed, return false
}

func FetchUpdate() error {
	updated, err := UpdateAvailable()
	if err != nil {
		return err // Handle error appropriately, e.g., log it
	}
	localVersion, err := GetLocalVersion()
	if err != nil {
		return err
	}
	if updated {
		url := "https://ddragon.leagueoflegends.com/cdn/dragontail-" + localVersion + ".tgz"
		resp, err := http.Get(url)
		if err != nil {
			return err // Handle error appropriately, e.g., log it
		}
		defer resp.Body.Close()

		if err := os.MkdirAll("data", 0755); err != nil {
			return err // Handle error appropriately, e.g., log it
		}

		outPath := filepath.Join("data", "dragontail-"+localVersion+".tgz")
		out, err := os.Create(outPath)
		if err != nil {
			return err // Handle error appropriately, e.g., log it
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		extractDir := filepath.Join("data", "dragontail-"+localVersion)
		if err := ExtractTar(outPath, extractDir); err != nil {
			return err // Handle error appropriately, e.g., log it
		}
	}
	return nil
}

func ExtractTar(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err // Handle error appropriately, e.g., log it
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err // Handle error appropriately, e.g., log it
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			return err // Handle error appropriately, e.g., log it
		}
		targetPath := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(targetPath, 0755) // Create directory
		case tar.TypeReg:
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err // Handle error appropriately, e.g., log it
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err // Handle error appropriately, e.g., log it
			}
		default:
			continue // Skip other types
		}
	}
	return nil
}
