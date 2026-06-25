package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	pathpkg "path"
	"slices"
	"strconv"
	"strings"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

type apiFlags struct {
	method    string
	fields    []string
	rawFields []string
	queries   []string
	input     string
	headers   []string
	include   bool
	silent    bool
}

func newAPICmd(rt *app.Runtime) *cobra.Command {
	var flags apiFlags

	cmd := &cobra.Command{
		Use:   "api <endpoint>",
		Short: "Make a raw TAPD API request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			endpoint, err := normalizeAPIEndpoint(args[0])
			if err != nil {
				return err
			}

			method, err := normalizeAPIMethod(flags.method)
			if err != nil {
				return err
			}

			data, err := apiRequestData(cmd.InOrStdin(), method, flags)
			if err != nil {
				return err
			}

			opts, err := apiRequestOptions(flags.headers)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			req, err := client.NewRequest(cmd.Context(), method, endpoint, data, opts)
			if err != nil {
				return err
			}

			if err := appendAPIQuery(req.URL, nil, flags.queries); err != nil {
				return err
			}

			if !apiMethodHasBody(method) {
				if err := appendAPIQuery(req.URL, flags.fields, flags.rawFields); err != nil {
					return err
				}
			}

			var body json.RawMessage
			resp, err := client.Do(req, &body)
			if err != nil {
				return err
			}

			return writeAPIResponse(cmd.OutOrStdout(), resp, body, flags.include, flags.silent)
		},
	}

	cmd.Flags().StringVarP(&flags.method, "method", "X", http.MethodGet, "HTTP method: GET|POST|PUT|PATCH|DELETE")
	cmd.Flags().StringArrayVar(&flags.fields, "field", nil, "add a typed request parameter in key=value format")
	cmd.Flags().StringArrayVar(&flags.rawFields, "raw-field", nil, "add a string request parameter in key=value format")
	cmd.Flags().StringArrayVar(&flags.queries, "query", nil, "add a query parameter in key=value format")
	cmd.Flags().StringVar(&flags.input, "input", "", "JSON request body file, or - for stdin")
	cmd.Flags().StringArrayVarP(&flags.headers, "header", "H", nil, "add a request header in key:value format")
	cmd.Flags().BoolVarP(&flags.include, "include", "i", false, "include HTTP response status and headers")
	cmd.Flags().BoolVar(&flags.silent, "silent", false, "do not print response body")

	return cmd
}

func normalizeAPIEndpoint(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", errors.New("endpoint cannot be empty")
	}
	if strings.Contains(endpoint, "\\") || strings.ContainsAny(endpoint, "?#") {
		return "", fmt.Errorf("invalid endpoint %q, expected a TAPD API relative path", endpoint)
	}
	if strings.Contains(endpoint, "://") || strings.HasPrefix(endpoint, "//") {
		return "", fmt.Errorf("invalid endpoint %q, absolute URLs are not allowed", endpoint)
	}
	endpoint = strings.TrimPrefix(endpoint, "/")

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse endpoint: %w", err)
	}
	if parsed.IsAbs() || parsed.Host != "" {
		return "", fmt.Errorf("invalid endpoint %q, absolute URLs are not allowed", endpoint)
	}
	unescaped, err := url.PathUnescape(parsed.Path)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint %q, expected a TAPD API relative path", endpoint)
	}
	if strings.Contains(unescaped, "\\") || strings.ContainsAny(unescaped, "?#") {
		return "", fmt.Errorf("invalid endpoint %q, expected a TAPD API relative path", endpoint)
	}
	if strings.HasPrefix(unescaped, "//") {
		return "", fmt.Errorf("invalid endpoint %q, absolute URLs are not allowed", endpoint)
	}
	if hasParentPathSegment(unescaped) {
		return "", fmt.Errorf("invalid endpoint %q, parent path segments are not allowed", endpoint)
	}

	clean := pathpkg.Clean(unescaped)
	if clean == "." {
		return "", fmt.Errorf("invalid endpoint %q, expected a TAPD API relative path", endpoint)
	}
	if clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
		return "", fmt.Errorf("invalid endpoint %q, parent path segments are not allowed", endpoint)
	}

	return clean, nil
}

func hasParentPathSegment(path string) bool {
	for segment := range strings.SplitSeq(path, "/") {
		if segment == ".." {
			return true
		}
	}
	return false
}

func normalizeAPIMethod(method string) (string, error) {
	method = strings.ToUpper(strings.TrimSpace(method))
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return method, nil
	default:
		return "", fmt.Errorf("unsupported method %q, expected GET, POST, PUT, PATCH, or DELETE", method)
	}
}

func apiRequestData(stdin io.Reader, method string, flags apiFlags) (any, error) {
	if flags.input != "" {
		if len(flags.fields) > 0 || len(flags.rawFields) > 0 {
			return nil, errors.New("--input cannot be mixed with --field or --raw-field")
		}
		if !apiMethodHasBody(method) {
			return nil, fmt.Errorf("--input requires POST, PUT, or PATCH")
		}
		return readAPIInput(stdin, flags.input)
	}

	if !apiMethodHasBody(method) {
		return nil, nil
	}

	return apiFieldsData(flags.fields, flags.rawFields)
}

func readAPIInput(stdin io.Reader, input string) (json.RawMessage, error) {
	var (
		data []byte
		err  error
	)
	if input == "-" {
		data, err = io.ReadAll(stdin)
	} else {
		data, err = os.ReadFile(input)
	}
	if err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}
	data = bytes.TrimSpace(data)
	if !json.Valid(data) {
		return nil, fmt.Errorf("input must contain valid JSON")
	}
	return json.RawMessage(data), nil
}

func apiFieldsData(fields, rawFields []string) (map[string]any, error) {
	data := make(map[string]any, len(fields)+len(rawFields))
	for _, field := range rawFields {
		name, value, err := parseAPIField(field)
		if err != nil {
			return nil, err
		}
		data[name] = value
	}
	for _, field := range fields {
		name, value, err := parseTypedAPIField(field)
		if err != nil {
			return nil, err
		}
		data[name] = value
	}
	if len(data) == 0 {
		return nil, nil
	}
	return data, nil
}

func appendAPIQuery(u *url.URL, fields, rawFields []string) error {
	q := u.Query()
	for _, field := range rawFields {
		name, value, err := parseAPIField(field)
		if err != nil {
			return err
		}
		q.Add(name, value)
	}
	for _, field := range fields {
		name, value, err := parseTypedAPIField(field)
		if err != nil {
			return err
		}
		q.Add(name, apiQueryValue(value))
	}
	u.RawQuery = q.Encode()
	return nil
}

func parseAPIField(field string) (string, string, error) {
	name, value, ok := strings.Cut(field, "=")
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		return "", "", fmt.Errorf("invalid field %q, expected key=value", field)
	}
	return name, value, nil
}

func parseTypedAPIField(field string) (string, any, error) {
	name, value, err := parseAPIField(field)
	if err != nil {
		return "", nil, err
	}

	typed, err := apiTypedValue(strings.TrimSpace(value))
	if err != nil {
		return "", nil, err
	}
	return name, typed, nil
}

func apiTypedValue(value string) (any, error) {
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	}

	if file, ok := strings.CutPrefix(value, "@"); ok {
		if file == "" {
			return nil, errors.New("field file reference cannot be empty")
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read field file: %w", err)
		}
		return string(data), nil
	}

	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f, nil
	}

	return value, nil
}

func apiQueryValue(value any) string {
	switch value := value.(type) {
	case nil:
		return "null"
	case bool:
		return strconv.FormatBool(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case string:
		return value
	default:
		return fmt.Sprint(value)
	}
}

func apiRequestOptions(headers []string) ([]tapd.RequestOption, error) {
	opts := make([]tapd.RequestOption, 0, len(headers))
	for _, header := range headers {
		name, value, ok := strings.Cut(header, ":")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if !ok || name == "" {
			return nil, fmt.Errorf("invalid header %q, expected key:value", header)
		}
		if strings.EqualFold(name, "Authorization") {
			return nil, errors.New("overriding Authorization is not allowed")
		}
		opts = append(opts, tapd.WithRequestHeader(name, value))
	}
	return opts, nil
}

func apiMethodHasBody(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}

func writeAPIResponse(w io.Writer, resp *tapd.Response, body json.RawMessage, include, silent bool) error {
	if include && resp != nil && resp.Response != nil {
		if _, err := fmt.Fprintf(w, "%s %s\n", resp.Proto, resp.Status); err != nil {
			return err
		}

		keys := make([]string, 0, len(resp.Header))
		for key := range resp.Header {
			keys = append(keys, key)
		}
		slices.Sort(keys)
		for _, key := range keys {
			if _, err := fmt.Fprintf(w, "%s: %s\n", key, strings.Join(resp.Header.Values(key), ", ")); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	if silent {
		return nil
	}
	if len(body) == 0 {
		body = json.RawMessage("null")
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		if _, writeErr := fmt.Fprintln(w, string(body)); writeErr != nil {
			return writeErr
		}
		return nil
	}
	if err := pretty.WriteByte('\n'); err != nil {
		return err
	}
	_, err := w.Write(pretty.Bytes())
	return err
}
