package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"ur/database"
	"ur/file"
	"ur/template"
	"ur/user"

	"goji.io/pat"

	goji "goji.io"

	"github.com/NYTimes/gziphandler"
	"github.com/yourheropaul/inj"
)

const (
	devFilePath      = "src/ur/admin/frontend"
	aceTemplatesPath = "ace"
	configFilePath   = "admin.scl"
)

var devModeFlag = flag.Bool("dev", false, "Run server in development mode")

type application struct {
	Environment            environmentMap          `inj:""`
	Config                 configurer              `inj:""`
	DatabaseSessionFactory database.SessionFactory `inj:""`
	ModelModelFactory      user.ModelFactory       `inj:""`
	TemplateFactory        template.Factory        `inj:""`
	FileSystem             file.System             `inj:""`
	Mux                    *goji.Mux               `inj:""`
}

func (a *application) httpError(w http.ResponseWriter, code int, err string, args ...interface{}) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(err, args...)))
}

func (a *application) frontendHTTPHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := a.TemplateFactory.Template("index")

	if err != nil {
		a.httpError(w, 500, "Can't load index.ace: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.Execute(
		w,
		map[string]interface{}{
			"DevMode": *devModeFlag,
			"Config":  a.Config,
		},
	); err != nil {
		a.httpError(w, 500, "Can't load index.ace: %s", err)
		return
	}
}

var setupProcedures = []struct {
	description string
	fn          interface{}
}{
	{"Parsing command line flags", parseFlags},
	{"Scanning environment", setupEnvironmentMap},
	{"Set up error reporter", setupErrorReporter},
	{"Set up application", setupApplication},
	{"Set up file system", setupFileSystem},
	{"Load config", loadConfig},
	{"Setup database", setupDatabase},
	{"Setup user model factory", setupUserModelFactory},
	{"Create HTTP multiplexer", setupMux},
	{"Setup user API", setupUserAPI},
	{"Setup template factory", setupTemplateFactory},
	{"Assert dependency graph", assertDependencyGraph},
	{"Start frontend handlers", setupFrontendAppHandlers},
}

func parseFlags() {
	flag.Parse()
}

func setupEnvironmentMap(grapher inj.Grapher) {
	grapher.Provide(parseEnvironment())
}

func setupErrorReporter(grapher inj.Grapher) {
	grapher.Provide(newOsErrorReporter())
}

func setupApplication(grapher inj.Grapher) {
	grapher.Provide(&application{})
}

func setupFileSystem(grapher inj.Grapher) {
	// FIXME! Add binary data
	grapher.Provide(file.NewDiskSystem(devFilePath))
}

func loadConfig(grapher inj.Grapher, env environmentMap, fs file.System, reporter errorReporter) {

	cfg, err := newConfig(fs, configFilePath, env)

	if err != nil {
		reporter.fatal(err.Error())
	}

	grapher.Provide(cfg)
}

func setupDatabase(grapher inj.Grapher, cfg configurer, reporter errorReporter) {

	sf, err := database.NewFirebaseSessionFactory(
		cfg.database().Endpoint,
		cfg.database().ServiceCredentials,
	)

	if err != nil {
		reporter.fatal("Setup database: %s", err)
	}

	grapher.Provide(sf)
}

func setupUserModelFactory(grapher inj.Grapher, sf database.SessionFactory, reporter errorReporter) {
	grapher.Provide(user.NewDatabaseModelFactory(sf))
}

func setupMux(grapher inj.Grapher) {
	grapher.Provide(goji.NewMux())
}

func setupUserAPI(grapher inj.Grapher, uf user.ModelFactory, mux *goji.Mux, reporter errorReporter) {

	userMux := goji.SubMux()

	api := user.NewAPIV1(userMux, uf)

	if err := api.Register(); err != nil {
		reporter.fatal("Couldn't register user API: %s", err)
	}

	mux.Handle(pat.New("/api/v1/user/*"), userMux)
}

func setupTemplateFactory(grapher inj.Grapher, fs file.System) {
	grapher.Provide(template.NewAceFactory(fs, aceTemplatesPath, *devModeFlag))
}

func assertDependencyGraph(grapher inj.Grapher, reporter errorReporter) {

	if ok, errs := grapher.Assert(); !ok {
		for k, v := range errs {
			fmt.Printf("[%d] %s\n", k, v)
		}
		reporter.fatal("Exiting")
	}
}

func setupFrontendAppHandlers(app *application, mux *goji.Mux, fs file.System) {

	appMux := goji.SubMux()

	appMux.HandleFunc(pat.New(""), app.frontendHTTPHandler)

	mux.Handle(
		pat.New("/assets/*"),
		http.StripPrefix("/assets/",
			gziphandler.GzipHandler(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					if r.URL.Path == "" {
						w.WriteHeader(http.StatusNotFound)
						return
					}

					path := "dist/" + r.URL.Path

					if f, tm, err := fs.ReadCloser(path); err == nil {
						defer f.Close()
						http.ServeContent(w, r, "/assets/"+r.URL.Path, tm, newDelayedReadSeeker(f))
						return
					}

					w.WriteHeader(http.StatusNotFound)
				}),
			),
		),
	)

	mux.Handle(pat.New("/"), appMux)
}

func main() {

	grapher := inj.NewGraph()

	for _, procedure := range setupProcedures {
		fmt.Printf("%s... ", procedure.description)
		grapher.Inject(procedure.fn, grapher)
		fmt.Println("OK")
	}

	grapher.Inject(serve)

	select {}
}

func serve(mux *goji.Mux) {

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	fmt.Printf("Running server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
