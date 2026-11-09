package testtools

import (
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

func TestRandomForeignID(t *testing.T) {
	const count = 1000

	seen := make(map[string]struct{}, count)
	for range count {
		id := RandomForeignID()
		require.NotEmpty(t, id)

		_, duplicate := seen[id]
		require.False(t, duplicate, "duplicate foreign ID %q", id)
		seen[id] = struct{}{}
	}
}

var fixedForeignIDJSONProperty = regexp.MustCompile(`"(?:foreignID|externalID)"\s*:\s*"[^"]*"`)

type fixedForeignIDFinding struct {
	line int
	why  string
}

func TestLiveSuiteForeignIDs(t *testing.T) {
	t.Run("detects fixed foreign IDs", func(t *testing.T) {
		fixtures := []struct {
			name string
			src  string
			want int
		}{
			{
				name: "direct string value",
				src:  "package p\nvar t = Transfer{ForeignID: \"ref-123\"}\n",
				want: 1,
			},
			{
				name: "literal-backed alias",
				src:  "package p\nconst ref = \"ref-123\"\nvar t = Transfer{ForeignID: ref}\n",
				want: 1,
			},
			{
				name: "literal-backed alias address",
				src:  "package p\nvar ref = \"ref-123\"\nvar t = Transfer{ForeignID: &ref}\n",
				want: 1,
			},
			{
				name: "pointer-wrapped string",
				src:  "package p\nvar t = Transfer{ForeignID: moov.PtrOf(\"ref-123\")}\n",
				want: 1,
			},
			{
				name: "foreignID JSON fixture property",
				src:  "package p\nvar raw = []byte(`{\"foreignID\":\"ref-123\"}`)\n",
				want: 1,
			},
			{
				name: "externalID JSON fixture property",
				src:  "package p\nvar raw = []byte(`{\"externalID\":\"ref-123\"}`)\n",
				want: 1,
			},
			{
				name: "dynamic foreign ID",
				src:  "package p\nvar t = Transfer{ForeignID: testtools.RandomForeignID()}\n",
				want: 0,
			},
			{
				name: "pointer to dynamic foreign ID",
				src:  "package p\nvar t = Transfer{ForeignID: moov.PtrOf(testtools.RandomForeignID())}\n",
				want: 0,
			},
			{
				name: "fixed resource identifier",
				src:  "package p\nvar w = GetWallet{AccountID: \"ebbf46c6-122a-4367-bc45-7dd555e1d3b9\"}\n",
				want: 0,
			},
			{
				name: "prose mention",
				src:  "package p\n// foreignID comes from your system\nvar x = 1\n",
				want: 0,
			},
		}

		for _, fixture := range fixtures {
			t.Run(fixture.name, func(t *testing.T) {
				findings, err := findFixedForeignIDs(fixture.name+".go", []byte(fixture.src))
				require.NoError(t, err)
				require.Len(t, findings, fixture.want)
			})
		}
	})

	t.Run("live suite test files", func(t *testing.T) {
		for _, dir := range []string{"pkg", "examples"} {
			root := filepath.Join("..", "..", dir)
			err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if entry.IsDir() || !strings.HasSuffix(path, "_test.go") {
					return nil
				}

				src, err := os.ReadFile(path)
				require.NoError(t, err)

				findings, err := findFixedForeignIDs(path, src)
				require.NoError(t, err)
				require.Empty(t, findings, "%s uses a fixed foreign ID", path)
				return nil
			})
			require.NoError(t, err)
		}
	})
}

func findFixedForeignIDs(filename string, src []byte) ([]fixedForeignIDFinding, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	literalAliases := make(map[string]bool)
	collectLiteralAliases(file, literalAliases)

	var findings []fixedForeignIDFinding
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CompositeLit:
			for _, element := range n.Elts {
				kv, ok := element.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || !isForeignIDField(key.Name) {
					continue
				}
				if why := fixedForeignIDValue(kv.Value, literalAliases); why != "" {
					findings = append(findings, fixedForeignIDFinding{fset.Position(kv.Pos()).Line, why})
				}
			}
		case *ast.BasicLit:
			if n.Kind != token.STRING {
				return true
			}
			value, err := strconv.Unquote(n.Value)
			if err != nil {
				return true
			}
			if fixedForeignIDJSONProperty.MatchString(value) {
				findings = append(findings, fixedForeignIDFinding{fset.Position(n.Pos()).Line, "fixed foreignID/externalID JSON property"})
			}
		}
		return true
	})

	return findings, nil
}

func collectLiteralAliases(file *ast.File, aliases map[string]bool) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.ValueSpec:
			if allStringLiterals(n.Values) {
				for _, name := range n.Names {
					aliases[name.Name] = true
				}
			}
		case *ast.AssignStmt:
			if allStringLiterals(n.Rhs) {
				for _, lhs := range n.Lhs {
					if id, ok := lhs.(*ast.Ident); ok {
						aliases[id.Name] = true
					}
				}
			}
		}
		return true
	})
}

func allStringLiterals(exprs []ast.Expr) bool {
	if len(exprs) == 0 {
		return false
	}
	for _, expr := range exprs {
		lit, ok := expr.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return false
		}
	}
	return true
}

func isForeignIDField(name string) bool {
	return name == "ForeignID" || name == "ExternalID"
}

func fixedForeignIDValue(expr ast.Expr, aliases map[string]bool) string {
	switch value := expr.(type) {
	case *ast.BasicLit:
		if value.Kind == token.STRING {
			return "direct string foreign ID"
		}
	case *ast.Ident:
		if aliases[value.Name] {
			return "literal-backed foreign ID alias"
		}
	case *ast.UnaryExpr:
		if value.Op == token.AND {
			if id, ok := value.X.(*ast.Ident); ok && aliases[id.Name] {
				return "literal-backed foreign ID alias"
			}
		}
	case *ast.CallExpr:
		for _, arg := range value.Args {
			if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				return "pointer-wrapped foreign ID string"
			}
		}
	}
	return ""
}
