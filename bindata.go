// Code generated by go-bindata. DO NOT EDIT. @generated
// sources:
// res/init_s_ports.sql
// res/init_s_usage.sql
// res/init_s_user.sql
// res/init_s_user_password.sql
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _resInit_s_portsSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\x2e\x4a\x4d\x2c\x49\x55\x28\x49\x4c\xca\x49\x55\x28\x8e\x2f\xc8\x2f\x2a\x29\x56\xd0\xe0\xe5\x52\x50\x00\x31\x15\x14\x32\xf3\x4a\x14\x0a\x8a\x32\x73\x13\x8b\x2a\x15\xb2\x53\x2b\x75\x40\x32\x89\x39\x99\x89\xc5\x0a\x65\x89\x45\xc9\x19\x89\x45\x1a\x86\x06\x06\x9a\x0a\x29\xa9\x69\x89\xa5\x39\x25\x0a\xea\xcf\xe6\xac\x7a\xd9\xda\xfb\x7c\xef\x3a\x75\x5e\x2e\x4d\x6b\x5e\x2e\x40\x00\x00\x00\xff\xff\x55\x94\x8f\x1a\x60\x00\x00\x00")

func resInit_s_portsSqlBytes() ([]byte, error) {
	return bindataRead(
		_resInit_s_portsSql,
		"res/init_s_ports.sql",
	)
}

func resInit_s_portsSql() (*asset, error) {
	bytes, err := resInit_s_portsSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/init_s_ports.sql", size: 96, mode: os.FileMode(438), modTime: time.Unix(1536589515, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resInit_s_usageSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\x28\x8e\x2f\x2d\x4e\x4c\x4f\xe5\xe5\xd2\xe0\xe5\x52\x50\x28\xc8\x2f\x2a\x51\x00\x83\xcc\x3c\x28\xc3\xcf\x3f\x44\xc1\x2f\xd4\xc7\x47\x07\x24\x9f\x92\x58\x92\x0a\x11\x06\xb1\x4a\x32\x73\x53\xd1\xe4\xf3\xcb\xf3\x42\x41\xe6\x29\x24\x65\xa6\x83\x8c\x70\x71\x75\x73\x0c\xf5\x09\x51\x30\x40\x55\x57\x5a\x00\x51\xa5\x80\x5b\x1d\x2f\x97\xa6\x35\x2f\x17\x20\x00\x00\xff\xff\xf6\x41\x68\x3d\xac\x00\x00\x00")

func resInit_s_usageSqlBytes() ([]byte, error) {
	return bindataRead(
		_resInit_s_usageSql,
		"res/init_s_usage.sql",
	)
}

func resInit_s_usageSql() (*asset, error) {
	bytes, err := resInit_s_usageSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/init_s_usage.sql", size: 172, mode: os.FileMode(438), modTime: time.Unix(1539435921, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resInit_s_userSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x8e\x31\xcb\xc2\x30\x10\x86\xf7\x40\xfe\xc3\x3b\xb6\xf0\x0d\xfd\x66\xa7\xd8\xdc\x70\x34\xbd\x6a\xda\x80\x9d\x42\xd1\x82\x2e\x0e\x2d\x15\x7f\xbe\x08\x29\x52\x9d\x9f\xe7\xee\x7d\x4a\x4f\xa6\x23\x74\x66\xef\x08\x73\x5c\xe6\x71\xd2\x2a\xd3\x0a\x00\x42\x4b\x3e\xb2\xc5\x63\x98\xce\xd7\x61\xca\xfe\x8b\x22\xc7\xc1\x73\x6d\x7c\x8f\x8a\x7a\x48\xd3\x41\x82\x73\x7f\x1f\x5f\x4c\x4d\xdb\x83\x2f\x49\xb8\xac\x7e\xa4\xc4\xa8\x36\xec\xa2\xb1\xd6\x53\xdb\x6e\x04\xad\xf2\x9d\x56\x29\x36\x08\x1f\x03\x81\xc5\xd2\x29\x35\xc7\x75\x3b\x2e\xb7\xfb\x65\x7c\xa2\x91\x44\x90\xad\xe8\xfd\xe1\x15\x00\x00\xff\xff\x46\x92\x3c\x8e\xef\x00\x00\x00")

func resInit_s_userSqlBytes() ([]byte, error) {
	return bindataRead(
		_resInit_s_userSql,
		"res/init_s_user.sql",
	)
}

func resInit_s_userSql() (*asset, error) {
	bytes, err := resInit_s_userSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/init_s_user.sql", size: 239, mode: os.FileMode(438), modTime: time.Unix(1537010345, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resInit_s_user_passwordSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\x28\x8e\x2f\x2d\x4e\x2d\x8a\x2f\x48\x2c\x2e\x2e\xcf\x2f\x4a\xe1\xe5\xd2\xe0\xe5\x52\x50\x08\x0d\x76\x0d\x8a\xf7\x74\x51\x28\x4b\x2c\x4a\xce\x48\x2c\xd2\x30\x34\x30\xd0\x54\x08\x08\xf2\xf4\x75\x0c\x8a\x54\xf0\x76\x8d\xd4\x01\x29\x0a\x70\x0c\x0e\x0e\xf7\x0f\x42\x53\xe5\xe7\x1f\xa2\xe0\x17\xea\xe3\xc3\xcb\xa5\x69\xcd\xcb\x05\x08\x00\x00\xff\xff\x93\xa6\x64\x8a\x6c\x00\x00\x00")

func resInit_s_user_passwordSqlBytes() ([]byte, error) {
	return bindataRead(
		_resInit_s_user_passwordSql,
		"res/init_s_user_password.sql",
	)
}

func resInit_s_user_passwordSql() (*asset, error) {
	bytes, err := resInit_s_user_passwordSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/init_s_user_password.sql", size: 108, mode: os.FileMode(438), modTime: time.Unix(1537035523, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"res/init_s_ports.sql":         resInit_s_portsSql,
	"res/init_s_usage.sql":         resInit_s_usageSql,
	"res/init_s_user.sql":          resInit_s_userSql,
	"res/init_s_user_password.sql": resInit_s_user_passwordSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"res": &bintree{nil, map[string]*bintree{
		"init_s_ports.sql":         &bintree{resInit_s_portsSql, map[string]*bintree{}},
		"init_s_usage.sql":         &bintree{resInit_s_usageSql, map[string]*bintree{}},
		"init_s_user.sql":          &bintree{resInit_s_userSql, map[string]*bintree{}},
		"init_s_user_password.sql": &bintree{resInit_s_user_passwordSql, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
