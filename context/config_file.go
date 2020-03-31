package context

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const defaultHostname = "github.com"

type configEntry struct {
	User  string
	Token string `yaml:"oauth_token"`
}

/*

Problem: the config file is treated as exclusively an auth config right now. We hardcode our parsing to the current structure, precluding handling of other keys.

I can take a few approaches:

 1. keep this code the same but target auth.yml and make the variables auth specific
 2. update this code to do the same stuff but operate on a top level auth key

And in the background I'm nervous that we won't be able to preserve comments when writing the config
back out. I want to test that hypothesis; but I think the higher priority is getting the auth config
corded off.

*/

func parseOrSetupConfigFile(fn string) (*configEntry, error) {
	entry, err := parseConfigFile(fn)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return setupConfigFile(fn)
	}
	return entry, err
}

func parseConfigFile(fn string) (*configEntry, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseConfig(f)
}

// ParseDefaultConfig reads the configuration file
func ParseDefaultConfig() (*configEntry, error) {
	return parseConfigFile(configFile())
}

type config struct {
	Root *yaml.Node
	Auth struct {
		Host  string
		Token string
	}
	Editor   string
	Protocol string
}

func defaultConfig() config {
	return config{
		Protocol: "https",
		// we leave editor as empty string to signal that we should use environment variables
	}
}

func parseConfig(r io.Reader) (*configEntry, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var root yaml.Node
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}
	if len(root.Content) < 1 {
		return nil, fmt.Errorf("malformed config")
	}
	if root.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected a top level map")
	}

	config := defaultConfig()

	fmt.Println("!!!!!!!!!!!!!!!!!!!")
	for i, v := range root.Content[0].Content {
		switch v.Value {
		case "auth":
			fmt.Printf("found auth config at position %d\n", i)
			// TODO adapt existing code for auth
			//for i := 0; i < len(config.Content[0].Content[0].Content)-1; i = i + 2 {
			//	fmt.Println(config.Content[0].Content[i].Value)

			//	if config.Content[0].Content[i].Value == defaultHostname {
			//		var entries []configEntry
			//		err = config.Content[0].Content[i+1].Decode(&entries)
			//		if err != nil {
			//			return nil, err
			//		}
			//		return &entries[0], nil
			//	}
			//}
		case "protocol":
			protocolValue := root.Content[0].Content[i+1].Value
			if protocolValue != "ssh" && protocolValue != "https" {
				return nil, fmt.Errorf("got unexpected value for protocol: %s", protocolValue)
			}
			config.Protocol = protocolValue
		case "editor":
			editorValue := root.Content[0].Content[i+1].Value
			if !filepath.IsAbs(editorValue) {
				return nil, fmt.Errorf("editor should be an absolute path; got: %s", editorValue)
			}
			config.Editor = editorValue
		case "aliases":
			fmt.Printf("found alias config at position %d\n", i)
			fmt.Println("but alias support is not implemented yet sorry")
		}

	}
	fmt.Printf("%#v\n", config)
	fmt.Println("!!!!!!!!!!!!!!!!!!!")
	return nil, fmt.Errorf("could not find config entry for %q", defaultHostname)
}
