package discovery

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"log/slog"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/tools/go/packages"
)


type packageExplorer struct {
	cfg *Config
	
}

func NewPkgExplorer(optsFn...func(*Config) error) (px *packageExplorer, err error) {

	px = &packageExplorer{
		cfg: &Config{
		TagName: "env",
		MetaTags: []string{
			"envy",
			"default",
			"match",
			"required",
			"options",
			"help",
			"usage",
		},
	
		DirPattern: "./...",
		Context: context.Background(),
		Mode:  packages.LoadSyntax,
		Dir:   "",
		Env: []string{},
		Tests: false,
		
		},
	}
	if px.cfg.Cwd, err = os.Getwd(); err != nil {
		return
	}

	for _, fn := range optsFn {
		if err = fn(px.cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error applying options: %v\n", err)
			return
		}
	}
	fmt.Println(px.cfg)
	return
}


func Link(pkg *packages.Package) func( ast.Node) string {
	return func( node ast.Node) string {
		var buf = bytes.NewBuffer(nil)
		file := pkg.Fset.File(node.Pos())
		buf.WriteString(fmt.Sprintf("./%s", path.Base(file.Name())))

		return buf.String()
	}
}

func LinkWithLineNums(pkg *packages.Package) func(string, ast.Node) string {
	return func(title string, node ast.Node) string {
		var buf = bytes.NewBuffer(nil)
		file := pkg.Fset.File(node.Pos())
		
		start_ln := file.Line(node.Pos())
		end_ln := file.Line(node.End())
		buf.WriteString(fmt.Sprintf("%s:(%d,%d)",  file.Name(), start_ln, end_ln))

		return buf.String()
	}
}

type occurence struct {
	Package string
	PackagePath string
	File string
	Location string
	Pos int
	Type string

	Field string
	Tag string
	EnvVar string
	MetaTags []string
}

func (p *packageExplorer) Walk()  {
	slog.Info("Loading packages",  "tag", p.cfg.TagName, "and meta tags", p.cfg.MetaTags)
	var pkgs []*packages.Package
	var err error
	if pkgs, err = packages.Load(&packages.Config{
		Mode:  p.cfg.Mode,
		Dir:   p.cfg.Dir,
		Env:  append(os.Environ(), p.cfg.Env...),
		Tests: p.cfg.Tests,
	}, p.cfg.DirPattern); err != nil {
		return
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}
	var report = map[string] []occurence{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if st, ok := s.Type.(*ast.StructType); ok {
								for _, field := range st.Fields.List {
									if field.Tag != nil {
										if tag, ok := reflect.StructTag(field.Tag.Value[1:len(field.Tag.Value)-1]).Lookup(p.cfg.TagName); ok {
											sanitized_tag := strings.Split(tag,";")
											if len(sanitized_tag) == 0 {
												continue
											}
											for _, env_var := range strings.Split(sanitized_tag[0], "|") {
												var next_occurence = occurence{
														PackagePath: pkg.PkgPath,
														Package: pkg.Types.Name(),
														File: file.Name.String(),
														Pos: int(field.Tag.Pos()),
														Type: s.Name.Name,
														Field: field.Names[0].Name,
														Tag: field.Tag.Value[1:len(field.Tag.Value)-1],
														EnvVar: env_var,
														Location: LinkWithLineNums(pkg)("", field.Tag),
												}
												var occurences []occurence
												var exists bool
												if occurences, exists = report[env_var]; !exists {
													occurences = []occurence{}
													report[env_var] = occurences
												} 
												report[env_var] = append(occurences, next_occurence)
											}
										}
									}
								}
							}
						}
					}
				}
			}
			
		}
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
    t.AppendHeader(table.Row{ "Name",  "Tag", "Location",})
    t.SetStyle(table.Style{
        Name: "myNewStyle",
        Box: table.BoxStyle{
            BottomLeft:       "",
            BottomRight:      "",
            BottomSeparator:  "v",
            Left:             "",
            LeftSeparator:    "",
            MiddleHorizontal: "-",
            MiddleSeparator:  "+",
            MiddleVertical:   "|",
            PaddingLeft:      "  ",
            PaddingRight:     "  ",
            Right:            "",
            RightSeparator:   "",
            TopLeft:          "",
            TopRight:         "",
            TopSeparator:     "v",
            UnfinishedRow:    " ~~~",
        },
        Color: table.ColorOptions{
            IndexColumn:     text.Colors{text.BgBlack, text.FgWhite},
            Footer:          text.Colors{text.BgBlack, text.FgWhite},
            Header:          text.Colors{text.BgBlack, text.FgHiCyan},
            Row:             text.Colors{text.BgBlack, text.FgHiGreen},
			Separator: text.Colors{text.BgBlack, text.FgWhite},
			Border: text.Colors{text.BgBlack, text.FgWhite},
            // RowAlternate:    text.Colors{text.Bg, text.FgBlack},
        },
        Format: table.FormatOptions{
            Footer: text.FormatUpper,
            Header: text.FormatUpper,
            Row:    text.FormatDefault,
			RowAlign: text.AlignLeft,
        },
        Options: table.Options{
            DrawBorder:      true,
            SeparateColumns: true,
            SeparateFooter:  true,
            SeparateHeader:  true,
            SeparateRows:    false,
        },
    })
	for env_var, occrs := range report {
		for _, occ := range occrs {
			  t.AppendRows([]table.Row{{ env_var, occ.Tag,occ.Location, }})
		}
	}
	t.Render()
}
func DarkGreen(text string) string {
	return rgb(38, 0, 160, 0, text)
}
func Green(text string) string {
	return rgb( 38, 0, 255, 0, text)
}
func Purple(text string) string {
	return rgb( 38, 110, 0, 110, text)
}
func Red(text string) string {
	return color(31, text)
}
func Yellow(text string) string {
	return color(33, text)
}
func Blue(text string) string {
	return color(34, text)
}
func Magenta(text string) string {
	return color(35, text)
}
func Cyan(text string) string {
	return color(36, text)
}

func Info(text string) string {
	return rgb(38, 0, 255, 255, text) // Cyan color for info
}


func rgb(applies_to, r, g, b int, text string) string {
	if applies_to != 38 && applies_to != 48 {
		return text // return original text if applies_to is not 38 or 48
	}
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return text // return original text if color values are out of range
	}
	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm%s\x1b[0m", applies_to,r, g, b, text)
}


func color(color_code int, text string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_code, text)
}


type Config struct {
	Cwd string 
	DirPattern string `env:"ENVY_PACKAGE_LOAD_DIR_PATTERN" help:"Directory pattern to match." default:"./..."`
	// Mode controls the level of information returned for each package.
	Mode packages.LoadMode `env:"ENVY_PACKAGES_LOAD_MODE" default:"991" help:"Mode controls the level of information returned for each package."`

	// Context specifies the context for the load operation.
	// Cancelling the context may cause [Load] to abort and
	// return an error.
	Context context.Context `ctx:"timeout=30s"`
	// Nower func() time.Time 
	// Field string `default:"{{call .Nower}}"`
	// Logf is the logger for the config.
	// If the user provides a logger, debug logging is enabled.
	// If the GOPACKAGESDEBUG environment variable is set to true,
	// but the logger is nil, default to log.Printf.
	// Logf func(format string, args ...any)

	// Dir is the directory in which to run the build system's query tool
	// that provides information about the packages.
	// If Dir is empty, the tool is run in the current directory.
	Dir string `env:"ENVY_PACKAGE_LOAD_DIR" help:"Directory in which to run the build system's query tool. Defaults to the current working directory."`

	// Env is the environment to use when invoking the build system's query tool.
	// If Env is nil, the current environment is used.
	// As in os/exec's Cmd, only the last value in the slice for
	// each environment key is used. To specify the setting of only
	// a few variables, append to the current environment, as in:
	//
	//	opt.Env = append(os.Environ(), "GOOS=plan9", "GOARCH=386")
	//
	Env []string `env:"ENVY_PACKAGE_LOAD_ENV" help:"Environment variables to set when invoking the build system's query tool."`

	// BuildFlags is a list of command-line flags to be passed through to
	// the build system's query tool.
	// BuildFlags []string

	// Fset provides source position information for syntax trees and types.
	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
	// Fset *token.FileSet

	// ParseFile is called to read and parse each file
	// when preparing a package's type-checked syntax tree.
	// It must be safe to call ParseFile simultaneously from multiple goroutines.
	// If ParseFile is nil, the loader will uses parser.ParseFile.
	//
	// ParseFile should parse the source from src and use filename only for
	// recording position information.
	//
	// An application may supply a custom implementation of ParseFile
	// to change the effective file contents or the behavior of the parser,
	// or to modify the syntax tree. For example, selectively eliminating
	// unwanted function bodies can significantly accelerate type checking.
	// ParseFile func(fset *token.FileSet, filename string, src []byte) (*ast.File, error)

	// If Tests is set, the loader includes not just the packages
	// matching a particular pattern but also any related test packages,
	// including test-only variants of the package and the test executable.
	//
	// For example, when using the go command, loading "fmt" with Tests=true
	// returns four packages, with IDs "fmt" (the standard package),
	// "fmt [fmt.test]" (the package as compiled for the test),
	// "fmt_test" (the test functions from source files in package fmt_test),
	// and "fmt.test" (the test binary).
	//
	// In build systems with explicit names for tests,
	// setting Tests may have no effect.
	Tests bool `env:"ENVY_PACKAGE_LOAD_TESTS" default:"true" help:"Include test packages."`

	// Overlay is a mapping from absolute file paths to file contents.
	//
	// For each map entry, [Load] uses the alternative file
	// contents provided by the overlay mapping instead of reading
	// from the file system. This mechanism can be used to enable
	// editor-integrated tools to correctly analyze the contents
	// of modified but unsaved buffers, for example.
	//
	// The overlay mapping is passed to the build system's driver
	// (see "The driver protocol") so that it too can report
	// consistent package metadata about unsaved files. However,
	// drivers may vary in their level of support for overlays.
	// Overlay map[string][]byte
	TagName string `env:"ENVY_LS_ROOT_TAG_NAME" default:"env"`
	MetaTags []string `env:"ENVY_LS_META_TAGS" default:"envy|default|match|required|options|help|usage"` 
}