package testtools

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var foreignIDName = regexp.MustCompile(`(?i)(foreign|external).*id`)

func TestNoFixedForeignIDs(t *testing.T) {
	t.Run("repository", func(t *testing.T) {
		for _, root := range []string{
			filepath.Join("..", "..", "pkg", "moov"),
			filepath.Join("..", "..", "examples"),
		} {
			err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() || !strings.HasSuffix(path, "_test.go") {
					return nil
				}
				src, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				for _, finding := range fixedForeignIDs(path, src) {
					t.Errorf("%s: fixed foreign identifier: %s", path, finding)
				}
				return nil
			})
			require.NoError(t, err)
		}
	})

	t.Run("fixtures", func(t *testing.T) {
		for _, tc := range foreignIDFixtures {
			t.Run(tc.name, func(t *testing.T) {
				got := fixedForeignIDs("fixture_test.go", []byte(tc.src))
				if tc.rejected {
					require.NotEmpty(t, got)
				} else {
					require.Empty(t, got)
				}
			})
		}
	})
}

type foreignIDFixture struct {
	name     string
	src      string
	rejected bool
}

var foreignIDFixtures = []foreignIDFixture{
	{
		name:     "named declaration",
		src:      "package fixture\n\nvar foreignID = \"fixed-value\"\n",
		rejected: true,
	},
	{
		name:     "named assignment",
		src:      "package fixture\n\nfunc f() {\n\texternalID := \"fixed-value\"\n\t_ = externalID\n}\n",
		rejected: true,
	},
	{
		name:     "foreign ID field",
		src:      "package fixture\n\nvar r = request{ForeignID: \"fixed-value\"}\n",
		rejected: true,
	},
	{
		name:     "pointer wrapped value",
		src:      "package fixture\n\nvar r = request{ForeignID: moov.PtrOf(\"fixed-value\")}\n",
		rejected: true,
	},
	{
		name:     "literal bound identifier",
		src:      "package fixture\n\nfunc f() {\n\tid := \"fixed-value\"\n\t_ = request{ForeignID: id}\n}\n",
		rejected: true,
	},
	{
		name:     "JSON object field",
		src:      "package fixture\n\nvar body = []byte(`{\"foreignID\":\"fixed-value\"}`)\n",
		rejected: true,
	},
	{
		name:     "computed identifier",
		src:      "package fixture\n\nfunc f() {\n\tforeignID := newID()\n\t_ = request{ForeignID: foreignID}\n}\n",
		rejected: false,
	},
	{
		name:     "seeded resource identifier",
		src:      "package fixture\n\nvar accountID = \"ebbf46c6-122a-4367-bc45-7dd555e1d3b9\"\n",
		rejected: false,
	},
	{
		name:     "plain text",
		src:      "package fixture\n\nvar message = \"the foreignID is not set\"\n",
		rejected: false,
	},
	{
		name:     "JSON without foreign ID",
		src:      "package fixture\n\nvar body = []byte(`{\"walletID\":\"ec7e1848\"}`)\n",
		rejected: false,
	},
	{
		name:     "null field",
		src:      "package fixture\n\nvar r = request{ForeignID: moov.SetNull[string]()}\n",
		rejected: false,
	},
}

func fixedForeignIDs(filename string, src []byte) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return []string{"parse error: " + err.Error()}
	}

	bindings := literalStringBindings(file)

	var findings []string
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec:
			for i, name := range node.Names {
				if i >= len(node.Values) || !foreignIDName.MatchString(name.Name) {
					continue
				}
				if value, ok := resolveStringLiteral(node.Values[i], bindings); ok {
					findings = append(findings, name.Name+" = "+value)
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || i >= len(node.Rhs) || !foreignIDName.MatchString(ident.Name) {
					continue
				}
				if value, ok := resolveStringLiteral(node.Rhs[i], bindings); ok {
					findings = append(findings, ident.Name+" = "+value)
				}
			}
		case *ast.CompositeLit:
			for _, elt := range node.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || !isForeignIDField(key.Name) {
					continue
				}
				if value, ok := resolveStringLiteral(kv.Value, bindings); ok {
					findings = append(findings, key.Name+": "+value)
				}
			}
		case *ast.BasicLit:
			if node.Kind == token.STRING {
				if value, ok := stringValue(node); ok && jsonHasForeignID(value) {
					findings = append(findings, "JSON foreignID field")
				}
			}
		}
		return true
	})

	return findings
}

func literalStringBindings(file *ast.File) map[string]string {
	bindings := make(map[string]string)
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec:
			for i, name := range node.Names {
				if i >= len(node.Values) {
					continue
				}
				if value, ok := directString(node.Values[i]); ok {
					bindings[name.Name] = value
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || i >= len(node.Rhs) {
					continue
				}
				if value, ok := directString(node.Rhs[i]); ok {
					bindings[ident.Name] = value
				}
			}
		}
		return true
	})
	return bindings
}

func resolveStringLiteral(expr ast.Expr, bindings map[string]string) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return stringValue(e)
	case *ast.Ident:
		value, ok := bindings[e.Name]
		return value, ok
	case *ast.ParenExpr:
		return resolveStringLiteral(e.X, bindings)
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			return resolveStringLiteral(e.X, bindings)
		}
	case *ast.CallExpr:
		if arg := wrapperArg(e); arg != nil {
			return resolveStringLiteral(arg, bindings)
		}
	}
	return "", false
}

func directString(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	return stringValue(lit)
}

func stringValue(lit *ast.BasicLit) (string, bool) {
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return value, true
}

func wrapperArg(call *ast.CallExpr) ast.Expr {
	if len(call.Args) != 1 {
		return nil
	}
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		if fun.Name == "PtrOf" || fun.Name == "Set" {
			return call.Args[0]
		}
	case *ast.SelectorExpr:
		if fun.Sel.Name == "PtrOf" || fun.Sel.Name == "Set" {
			return call.Args[0]
		}
	}
	return nil
}

func isForeignIDField(name string) bool {
	return strings.EqualFold(name, "ForeignID") || strings.EqualFold(name, "ExternalID")
}

func jsonHasForeignID(value string) bool {
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return false
	}
	return jsonValueHasForeignID(decoded)
}

func jsonValueHasForeignID(value any) bool {
	switch decoded := value.(type) {
	case map[string]any:
		for key, child := range decoded {
			if isForeignIDField(key) || jsonValueHasForeignID(child) {
				return true
			}
		}
	case []any:
		for _, child := range decoded {
			if jsonValueHasForeignID(child) {
				return true
			}
		}
	}
	return false
}
