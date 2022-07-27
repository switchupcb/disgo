package tools

import (
	"fmt"
	"go/format"
	"sort"
	"strings"
)

var (
	// constMap represents a map of constants referenced in functions (`key = val`)
	constMap map[string]string

	// EndpointBaseURL represents the EndpointBaseURL declaration variable name.
	EndpointBaseURL = "EndpointBaseURL"

	// CDNEndpointBaseURL represents the CDNEndpointBaseURL declaration variable name.
	CDNEndpointBaseURL = "CDNEndpointBaseURL"
)

// Endpoints reads the contents of a dasgo file and outputs disgo endpoints.
func Endpoints(data []byte) ([]byte, error) {
	constMap = make(map[string]string)
	constMap["slash"] = "/"
	content := parseEndpointDecl(string(data))
	contentdata := []byte(content)

	// gofmt
	fmtdata, err := format.Source(contentdata)
	if err != nil {
		fmt.Println(string(contentdata))

		return contentdata, fmt.Errorf("an error occurred while formatting the generated code.\n%w", err)
	}

	return fmtdata, nil
}

// parseEndpointDecl parses the endpoint declarations into a map.
func parseEndpointDecl(content string) string {
	var funcput strings.Builder

	base := EndpointBaseURL
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		decl := strings.Fields(line)

		if len(decl) > 0 {
			switch decl[0] {
			case EndpointBaseURL:
				constMap[EndpointBaseURL] = strings.Join(decl[2:], "")

			case CDNEndpointBaseURL:
				constMap[CDNEndpointBaseURL] = strings.Join(decl[2:], "")
				base = CDNEndpointBaseURL

			default:
				if len(decl) == 3 && decl[0] != "//" {
					url := decl[2][1 : len(decl[2])-1]
					funcput.WriteString(generateComment(decl[0]) + "\n")
					funcput.WriteString(generateFunc(decl[0], base, url) + "\n")
				}
			}
		}
	}

	var output strings.Builder
	output.WriteString("package wrapper\n\n")
	output.WriteString(generateConst(constMap))
	output.WriteString(funcput.String() + "\n")
	output.WriteString(manual())

	return output.String()
}

// generateConst generates the constant declarations for the endpoint file.
func generateConst(cm map[string]string) string {
	var decl strings.Builder
	decl.WriteString("// Discord API Endpoints\n")
	decl.WriteString("const (\n")
	decl.WriteString(EndpointBaseURL + " = " + cm[EndpointBaseURL] + "\n")
	decl.WriteString(CDNEndpointBaseURL + " = " + cm[CDNEndpointBaseURL] + "\n")
	delete(cm, EndpointBaseURL)
	delete(cm, CDNEndpointBaseURL)

	decls := make([]string, len(cm))
	i := 0
	for variable, value := range cm {
		decls[i] = variable + " = " + "\"" + value + "\""
		i++
	}

	sort.Slice(decls, func(i, j int) bool {
		return decls[i] < decls[j]
	})

	for _, line := range decls {
		decl.WriteString(line + "\n")
	}
	decl.WriteString(")\n")

	return decl.String()
}

// generateComment generates a comment for an endpoint function.
func generateComment(endpoint string) string {
	return "// " + endpoint + " builds a query for an HTTP request."
}

// generateFunc generates an endpoint function.
func generateFunc(endpoint, base, url string) string {
	urlparams := parameters(url)
	funcparams := make([]string, 0, len(urlparams))
	for _, param := range urlparams {
		if constMap[param] == "" {
			funcparams = append(funcparams, param)
		}
	}

	var p string
	if len(funcparams) > 0 {
		p = strings.Join(funcparams, ",") + " string"
	}

	var f strings.Builder
	f.WriteString("func " + endpoint + "(" + p + ")" + " string {\n")
	f.WriteString("return " + base + " +" + strings.Join(urlparams, "+ slash +") + "\n")
	f.WriteString("}\n")

	return f.String()
}

// parameters returns a list of parameters from a Discord API Endpoint URL.
func parameters(url string) []string {
	params := strings.Split(url, "/")
	for i, param := range params {
		// filter invalid characters (for a variable) from parameters
		params[i] = alphastring(param)

		// add constant variables to the constant map (`filtered = raw`)
		if param[0] != '{' {
			constMap[params[i]] = param
		}
	}

	return params
}

// alphastring returns a copy of s with only alphabetical characters.
func alphastring(s string) string {
	var alpha strings.Builder
	for _, c := range s {
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
			alpha.WriteString(string(c))
		}
	}

	return alpha.String()
}

// manual returns manual code added to the end of the file.
func manual() string {
	var m strings.Builder
	m.WriteString("var (")
	m.WriteString("EndpointModifyChannelGroupDM = EndpointModifyChannel\n")
	m.WriteString("EndpointModifyChannelGuild = EndpointModifyChannel\n")
	m.WriteString("EndpointModifyChannelThread = EndpointModifyChannel\n")
	m.WriteString(")")

	return m.String()
}
