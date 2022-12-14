package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/goutil/strutil"
	"go/ast"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Generator struct {
	spinner spinner.Model
	cfg     *Config
	pg      *Program
	maker   *Maker
	done    bool

	cmdPath      string
	internalPath string
	bizPath      string
	confPath     string
	dataPath     string
	entPath      string
	servicePath  string

	err error
}

func (g *Generator) solvePath() {
	g.cmdPath = filepath.Join(g.cfg.modPath, "cmd")
	g.internalPath = filepath.Join(g.cfg.modPath, "internal")
	g.bizPath = filepath.Join(g.internalPath, "biz")
	g.confPath = filepath.Join(g.internalPath, "conf")
	g.dataPath = filepath.Join(g.internalPath, "data")
	g.entPath = filepath.Join(g.internalPath, "ent")
	g.servicePath = filepath.Join(g.internalPath, "service")

	if g.cfg.typ > 0 {
		return
	}

	g.maker.MKDir(g.cfg.modPath)
	g.maker.MKDir(g.cmdPath)
	g.maker.MKDir(g.internalPath)
	g.maker.MKDir(g.bizPath)
	g.maker.MKDir(g.confPath)
	g.maker.MKDir(g.dataPath)
	g.maker.MKDir(g.servicePath)
}

func (g *Generator) genBizs() {
	if g.err != nil {
		return
	}

	g.maker.Template(filepath.Join(g.cmdPath, "main.go"), "main.gohtml", g.cfg)
	g.maker.Template(filepath.Join(g.cmdPath, "wire.go"), "wire.gohtml", g.cfg)
	g.maker.Template(filepath.Join(g.bizPath, "biz.go"), "biz.gohtml", g.cfg)
	g.maker.Template(filepath.Join(g.confPath, "conf.go"), "conf.gohtml", nil)
	g.maker.Template(filepath.Join(g.confPath, "config.yaml"), "config.yaml", nil)
	g.maker.Template(filepath.Join(g.dataPath, "data.go"), "data.gohtml", g.cfg)
	g.maker.Template(filepath.Join(g.servicePath, "service.go"), "service.gohtml", g.cfg)
}

func (g *Generator) genServices() {
	if g.err != nil {
		return
	}

	for _, name := range g.cfg.Services {
		lowerName := strings.ToLower(name)
		fileName := lowerName + ".go"
		g.cfg.CurrService = name
		g.cfg.CurrServiceLower = lowerName
		g.maker.Template(filepath.Join(g.bizPath, fileName), "bgreeter.gohtml", g.cfg)
		g.maker.Template(filepath.Join(g.dataPath, fileName), "dgreeter.gohtml", g.cfg)
		g.maker.Template(filepath.Join(g.servicePath, fileName), "sgreeter.gohtml", g.cfg)
	}
}

func (g *Generator) createModule() {
	if g.err != nil {
		return
	}

	g.maker.Chdir(g.cfg.modPath)
	g.maker.Cmd("go mod init", g.cfg.ModURL)
	if g.usingWork() {
		g.maker.Chdir(filepath.Dir(g.cfg.modPath))
		g.maker.Cmd("go work use", g.cfg.ModName)
		g.maker.Cmd("go work sync")
	}
	g.genEnt()
	g.genWire()
}

func (g *Generator) genWire() {
	if g.err != nil {
		return
	}

	g.maker.Chdir(g.cmdPath)
	g.maker.Cmd("go get -u github.com/google/wire")
	g.maker.Cmd("go run github.com/google/wire/cmd/wire ./...")
}

func (g *Generator) genEnt() {
	if g.err != nil || len(g.cfg.Services) == 0 {
		return
	}

	g.maker.Chdir(g.internalPath)
	g.maker.Cmd("go get -u entgo.io/ent/cmd/ent")
	g.maker.Cmd("go run entgo.io/ent/cmd/ent init", g.cfg.Services...)
	g.maker.Chdir(g.entPath)
	g.maker.Cmd("go run entgo.io/ent/cmd/ent generate ./schema")
}

func (g *Generator) usingWork() bool {
	_, err := os.Stat(filepath.Join(filepath.Dir(g.cfg.modPath), "go.work"))
	return err == nil
}

func (g *Generator) addService(suffix string) func(int) ([]string, []string, error) {
	return func(int) ([]string, []string, error) {
		var services []string
		for _, s := range g.cfg.Services {
			services = append(services, "New"+s+suffix)
		}
		return services, nil, nil
	}
}

func (g *Generator) margeService() {
	if g.err != nil {
		return
	}

	spath := filepath.Join(g.servicePath, "service.go")
	g.maker.MergeService(filepath.Join(g.bizPath, "biz.go"), "ProviderSet", ast.Var, g.addService("Usecase"))
	g.maker.MergeService(filepath.Join(g.dataPath, "data.go"), "ProviderSet", ast.Var, g.addService("Repo"))
	g.maker.MergeService(spath, "ProviderSet", ast.Var, g.addService("Service"))
	g.maker.MergeService(spath, "MakeServices", ast.Fun, func(count int) ([]string, []string, error) {
		var ps, rs []string
		for _, s := range g.cfg.Services {
			ps = append(ps, "s"+strconv.Itoa(count)+" *"+s+"Service")
			rs = append(rs, "s"+strconv.Itoa(count))
			count++
		}
		return ps, rs, nil
	})
}

func (g *Generator) removeService() {
	if g.err != nil {
		return
	}

	for _, s := range g.cfg.Services {
		fileName := strings.ToLower(s) + ".go"
		g.maker.RemoveAll(filepath.Join(g.bizPath, fileName))
		g.maker.RemoveAll(filepath.Join(g.dataPath, fileName))
		g.maker.RemoveAll(filepath.Join(g.servicePath, fileName))
		g.maker.RemoveAll(filepath.Join(g.entPath, "schema", fileName))
	}
	g.maker.Chdir(g.entPath)
	g.maker.Cmd("go run entgo.io/ent/cmd/ent generate ./schema")
}

func (g *Generator) revertService() {
	if g.err != nil {
		return
	}

	for _, s := range g.cfg.Services {
		g.maker.RevertService(filepath.Join(g.bizPath, "biz.go"), "ProviderSet", "New"+s+"Usecase", ast.Var)
		g.maker.RevertService(filepath.Join(g.dataPath, "data.go"), "ProviderSet", "New"+s+"Repo", ast.Var)
		g.maker.RevertService(filepath.Join(g.servicePath, "service.go"), "ProviderSet", "New"+s+"Service", ast.Var)
		g.maker.RevertService(filepath.Join(g.servicePath, "service.go"), "MakeServices", s+"Service", ast.Fun)
	}
}

func (g *Generator) revertWork() {
	if g.err != nil || !g.usingWork() {
		return
	}

	workPath := filepath.Join(filepath.Dir(g.cfg.modPath), "go.work")
	g.pg.Output("\x1b[1;33m[\x1b[1;5;33m%s\x1b[0m\x1b[1;33m]\x1b[0m remove module from \x1b[1;4;34m[%s]\x1b[0m", placeholder, workPath)
	before := time.Now()
	var f []byte
	f, g.err = os.ReadFile(workPath)
	if g.err != nil {
		return
	}
	reader := bufio.NewReader(bytes.NewReader(f))
	var buffer bytes.Buffer
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		line := string(buf)
		if !strings.HasSuffix(line, filepath.Base(g.cfg.modPath)) {
			buffer.Write(buf)
			buffer.WriteString("\n")
		}
	}

	if buffer.Len() > 0 {
		g.err = os.WriteFile(workPath, buffer.Bytes(), 0744)
		if g.err == nil {
			g.pg.Output("\x1b[1A\x1b[1;33m[\x1b[1;33m%10s\x1b[0m\x1b[1;33m]\x1b[0m remove module from \x1b[1;34m[%s]\x1b[0m",
				dur(time.Since(before)), workPath)
		}
	}
}

func (g *Generator) restoreModURL() {
	modPath := filepath.Join(g.cfg.modPath, "go.mod")
	var f *os.File
	f, g.err = os.Open(modPath)
	if g.err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		g.err = fmt.Errorf("go.mod mayby empty")
		return
	}

	txt := scanner.Text()
	var ts []string
	if strutil.IsNotBlank(txt) {
		ts = strings.Split(txt, " ")
	}
	if len(ts) < 2 {
		g.err = fmt.Errorf("go.mod is invalid")
		return
	}
	g.cfg.ModURL = ts[1]
}

func (g *Generator) gen() {
	g.pg.NewLine()
	g.solvePath()
	switch g.cfg.typ {
	case 0: // Create Module
		g.genBizs()
		g.genServices()
		g.createModule()
	case 1: // Add Service
		g.restoreModURL()
		g.genEnt()
		g.genServices()
		g.margeService()
		g.genWire()
	case 2: // Remove Module
		g.maker.RemoveAll(g.cfg.modPath)
		g.revertWork()
	case 3: // Remove Service
		g.removeService()
		g.revertService()
		g.genWire()
	}

	if g.err == nil {
		g.err = g.maker.err
	}

	if g.err != nil {
		msg := ExitErrStyle.Render(g.err.Error() + "\nexit status 1")
		if !g.maker.out {
			msg = "\x1b[1A" + msg
		}
		g.pg.Output(msg)
	} else if g.cfg.typ != 2 {
		g.pg.Output("\x1b[1A\x1b[1A")
	}

	g.done = true
}

func NewGenerator(cfg *Config, pg *Program) *Generator {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = SpinnerStyle
	g := &Generator{spinner: s, cfg: cfg, pg: pg, maker: NewMaker(pg)}
	go g.gen()
	return g
}

func (g *Generator) Update(msg tea.Msg) tea.Cmd {
	if g.done {
		return tea.Quit
	}

	var cmd tea.Cmd
	g.spinner, cmd = g.spinner.Update(msg)
	return cmd
}

func (g *Generator) View() string {
	if g.done {
		return "\n"
	}
	return TLBmarginStyle.Render(g.spinner.View()) + BlurredStyle.Render(" waiting for generate...\n")
}

func (g *Generator) Callback(*Program) (string, bool) {
	return "", false
}
