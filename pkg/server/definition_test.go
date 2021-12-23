package server

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getVM() (vm *jsonnet.VM) {
	vm = jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{JPaths: []string{"testdata"}})
	return
}

func TestDefinition(t *testing.T) {
	testCases := []struct {
		name     string
		params   protocol.DefinitionParams
		expected *protocol.DefinitionLink
	}{
		{
			name: "test goto definition for var myvar",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_goto_definition.jsonnet",
					},
					Position: protocol.Position{
						Line:      5,
						Character: 19,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/test_goto_definition.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      0,
						Character: 15,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      0,
						Character: 11,
					},
				},
			},
		},
		{
			name: "test goto definition on function helper",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_goto_definition.jsonnet",
					},
					Position: protocol.Position{
						Line:      7,
						Character: 8,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/test_goto_definition.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 6,
					},
					End: protocol.Position{
						Line:      1,
						Character: 23,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 6,
					},
					End: protocol.Position{
						Line:      1,
						Character: 12,
					},
				},
			},
		},
		{
			name: "test goto inner definition",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_goto_definition_multi_locals.jsonnet",
					},
					Position: protocol.Position{
						Line:      6,
						Character: 11,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/test_goto_definition_multi_locals.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      4,
						Character: 10,
					},
					End: protocol.Position{
						Line:      4,
						Character: 28,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      4,
						Character: 10,
					},
					End: protocol.Position{
						Line:      4,
						Character: 18,
					},
				},
			},
		},
		{
			name: "test goto super index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_combined_object.jsonnet",
					},
					Position: protocol.Position{
						Line:      5,
						Character: 13,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/test_combined_object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 5,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      1,
						Character: 5,
					},
				},
			},
		},
		{
			name: "test goto super nested",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_combined_object.jsonnet",
					},
					Position: protocol.Position{
						Line:      5,
						Character: 15,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/test_combined_object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 22,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 9,
					},
				},
			},
		},
		{
			name: "test goto self object field function",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/test_basic_lib.libsonnet",
					},
					Position: protocol.Position{
						Line:      4,
						Character: 19,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/test_basic_lib.libsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 20,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      1,
						Character: 9,
					},
				},
			},
		},
		{
			name: "test goto super object field local defined obj 'foo'",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/oo-contrived.jsonnet",
					},
					Position: protocol.Position{
						Line:      12,
						Character: 17,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/oo-contrived.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 2,
					},
					End: protocol.Position{
						Line:      1,
						Character: 8,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 2,
					},
					End: protocol.Position{
						Line:      1,
						Character: 5,
					},
				},
			},
		},
		{
			name: "test goto super object field local defined obj 'g'",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/oo-contrived.jsonnet",
					},
					Position: protocol.Position{
						Line:      13,
						Character: 17,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/oo-contrived.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 2,
					},
					End: protocol.Position{
						Line:      2,
						Character: 19,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 2,
					},
					End: protocol.Position{
						Line:      2,
						Character: 3,
					},
				},
			},
		},
		{
			name: "test goto local var from other local var",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/oo-contrived.jsonnet",
					},
					Position: protocol.Position{
						Line:      6,
						Character: 9,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/oo-contrived.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      3,
						Character: 1,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 6,
					},
					End: protocol.Position{
						Line:      0,
						Character: 10,
					},
				},
			},
		},
		{
			name: "test goto local obj field from 'self.attr' from other obj",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-indexes.jsonnet",
					},
					Position: protocol.Position{
						Line:      9,
						Character: 18,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/goto-indexes.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 23,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 8,
					},
					End: protocol.Position{
						Line:      2,
						Character: 11,
					},
				},
			},
		},
		{
			name: "test goto local object 'obj' via obj index 'obj.foo'",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-indexes.jsonnet",
					},
					Position: protocol.Position{
						Line:      8,
						Character: 15,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "./testdata/goto-indexes.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 5,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 4,
					},
					End: protocol.Position{
						Line:      1,
						Character: 7,
					},
				},
			},
		},
		{
			name: "test goto imported file",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 22,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-basic-object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 0,
					},
					End: protocol.Position{
						Line:      0,
						Character: 0,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 0,
					},
					End: protocol.Position{
						Line:      0,
						Character: 0,
					},
				},
			},
		},
		{
			name: "test goto imported file at lhs index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      3,
						Character: 18,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-basic-object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 7,
					},
				},
			},
		},
		{
			name: "test goto imported file at rhs index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "./testdata/goto-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      4,
						Character: 18,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-basic-object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto import index",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-import-attribute.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 48,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-basic-object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      5,
						Character: 4,
					},
					End: protocol.Position{
						Line:      5,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto attribute of nested import",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-nested-imported-file.jsonnet",
					},
					Position: protocol.Position{
						Line:      2,
						Character: 15,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-basic-object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto dollar attribute",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-dollar-simple.jsonnet",
					},
					Position: protocol.Position{
						Line:      7,
						Character: 17,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-dollar-simple.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 2,
					},
					End: protocol.Position{
						Line:      3,
						Character: 3,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      1,
						Character: 2,
					},
					End: protocol.Position{
						Line:      1,
						Character: 11,
					},
				},
			},
		},
		{
			name: "goto dollar sub attribute",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-dollar-simple.jsonnet",
					},
					Position: protocol.Position{
						Line:      8,
						Character: 28,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-dollar-simple.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 4,
					},
					End: protocol.Position{
						Line:      2,
						Character: 15,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 4,
					},
					End: protocol.Position{
						Line:      2,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto dollar doesn't follow to imports",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-dollar-no-follow.jsonnet",
					},
					Position: protocol.Position{
						Line:      7,
						Character: 13,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-dollar-no-follow.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 2,
					},
					End: protocol.Position{
						Line:      3,
						Character: 30,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 2,
					},
					End: protocol.Position{
						Line:      3,
						Character: 6,
					},
				},
			},
		},
		{
			name: "goto attribute of nested import no object intermediary",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-nested-import-file-no-inter-obj.jsonnet",
					},
					Position: protocol.Position{
						Line:      2,
						Character: 15,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-basic-object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto self in import in binary",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-self-within-binary.jsonnet",
					},
					Position: protocol.Position{
						Line:      4,
						Character: 15,
					},
				},
				WorkDoneProgressParams: protocol.WorkDoneProgressParams{},
				PartialResultParams:    protocol.PartialResultParams{},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-basic-object.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      3,
						Character: 4,
					},
					End: protocol.Position{
						Line:      3,
						Character: 7,
					},
				},
			},
		},
		{
			name: "goto self attribute from local",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-self-in-local.jsonnet",
					},
					Position: protocol.Position{
						Line:      3,
						Character: 23,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-self-in-local.jsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 2,
					},
					End: protocol.Position{
						Line:      2,
						Character: 21,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      2,
						Character: 2,
					},
					End: protocol.Position{
						Line:      2,
						Character: 12,
					},
				},
			},
		},
		{
			name: "goto function parameter from inside function",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-functions.libsonnet",
					},
					Position: protocol.Position{
						Line:      7,
						Character: 10,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-functions.libsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      6,
						Character: 10,
					},
					End: protocol.Position{
						Line:      6,
						Character: 14,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      6,
						Character: 10,
					},
					End: protocol.Position{
						Line:      6,
						Character: 14,
					},
				},
			},
		},
		{
			name: "goto local func param",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-local-function.libsonnet",
					},
					Position: protocol.Position{
						Line:      2,
						Character: 25,
					},
				},
			},
			expected: &protocol.DefinitionLink{
				TargetURI: absUri(t, "testdata/goto-local-function.libsonnet"),
				TargetRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 11,
					},
					End: protocol.Position{
						Line:      0,
						Character: 12,
					},
				},
				TargetSelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      0,
						Character: 11,
					},
					End: protocol.Position{
						Line:      0,
						Character: 12,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := string(tc.params.TextDocument.URI)
			var content, err = os.ReadFile(filename)
			require.NoError(t, err)
			ast, err := jsonnet.SnippetToAST(filename, string(content))
			require.NoError(t, err)
			got, err := Definition(ast, &tc.params, getVM())
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestDefinitionFail(t *testing.T) {
	testCases := []struct {
		name     string
		params   protocol.DefinitionParams
		expected error
	}{
		{
			name: "goto local keyword fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-basic-object.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 3,
					},
				},
			},
			expected: fmt.Errorf("cannot find definition"),
		},

		{
			name: "goto index of std fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-std.jsonnet",
					},
					Position: protocol.Position{
						Line:      1,
						Character: 20,
					},
				},
			},
			expected: fmt.Errorf("cannot get definition of std lib"),
		},
		{
			name: "goto comment fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-comment.jsonnet",
					},
					Position: protocol.Position{
						Line:      0,
						Character: 1,
					},
				},
			},
			expected: fmt.Errorf("cannot find definition"),
		},

		{
			name: "goto range index fails",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-local-function.libsonnet",
					},
					Position: protocol.Position{
						Line:      15,
						Character: 57,
					},
				},
			},
			expected: fmt.Errorf("unexpected node type when finding bind for 'ports'"),
		},
		{
			name: "goto super fails as no LHS object exists",
			params: protocol.DefinitionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "testdata/goto-local-function.libsonnet",
					},
					Position: protocol.Position{
						Line:      33,
						Character: 23,
					},
				},
			},
			expected: fmt.Errorf("could not find a lhs object"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := string(tc.params.TextDocument.URI)
			var content, err = os.ReadFile(filename)
			require.NoError(t, err)
			ast, err := jsonnet.SnippetToAST(filename, string(content))
			require.NoError(t, err)
			got, err := Definition(ast, &tc.params, getVM())
			require.Error(t, err)
			assert.Equal(t, tc.expected.Error(), err.Error())
			assert.Nil(t, got)
		})
	}
}

func absUri(t *testing.T, path string) protocol.DocumentURI {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)
	return protocol.URIFromPath(abs)
}
