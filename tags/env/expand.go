package env

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"strings"

	"github.com/dubbikins/envy/v2/tag"
)

func ExpandFn(n *tag.Node, reader Reader) func(string) string {
	slog.Info("expandfn")
	return func(s string) string {
		slog.Info("********EXPANDING******", "var", s, "node", n.Value().Type())
		if strings.HasPrefix(s,"$") {
				
			var err error
			t := template.New("$")
			s = fmt.Sprintf(`{{%s}}`, s)
			if t, err = t.Parse(s); err != nil {
				panic(err)
			}
			var root = n.Root()
			
			w := bytes.NewBuffer(nil)
			if root == nil {
				panic("can't perform root expansion on nil root")
			}
			if err = t.Execute(w, root.Value().Interface()); err != nil {
				panic(err)
			}
			//slog.Info("expanded root", "root", root.Value().Type(), "rootValue", root.Value().Interface(), "expandedValue", w.String(), "template", s)
			return w.String()
			
			
		}else if strings.HasPrefix(s,".") {
			
			//slog.Info("expanding ref", "node", n.Parent().Value().Type())
			var err error
			t := template.New("$")
			s = fmt.Sprintf(`{{%s}}`, s)
			if t, err = t.Parse(s); err != nil {
				panic(err)
			}
		
			w := bytes.NewBuffer(nil)
			if err = t.Execute(w, n.Parent().Value().Interface()); err != nil {
				panic(err)
			}
			//slog.Info("expanded root", "root", n.Parent().Value().Type(), "rootValue", n.Parent().Value().Interface(), "expandedValue", w.String(), "template", s)
			return w.String()
			
		}else {
			return reader(s)
		}
		
	}
}
