package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"ur/database"
	"ur/user"

	"goji.io/pat"

	goji "goji.io"

	"github.com/yourheropaul/inj"
)

type application struct {
	Environment            environmentMap          `inj:""`
	Config                 configurer              `inj:""`
	DatabaseSessionFactory database.SessionFactory `inj:""`
	ModelModelFactory      user.ModelFactory       `inj:""`
	Mux                    *goji.Mux               `inj:""`
}

var setupProcedures = []struct {
	description string
	fn          interface{}
}{
	{"Scanning environment", setupEnvironmentMap},
	{"Set up error reporter", setupErrorReporter},
	{"Set up application", setupApplication},
	{"Load config", loadConfig},
	{"Setup database", setupDatabase},
	{"Setup user model factory", setupUserModelFactory},
	{"Create HTTP multiplexer", setupMux},
	{"Setup user API", setupUserAPI},
	{"Assert dependency graph", assertDependencyGraph},
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

func loadConfig(grapher inj.Grapher, env environmentMap, reporter errorReporter) {

	cfg, err := newConfig("admin.scl", env)

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

func assertDependencyGraph(grapher inj.Grapher, reporter errorReporter) {

	if ok, errs := grapher.Assert(); !ok {
		for k, v := range errs {
			fmt.Printf("[%d] %s\n", k, v)
		}
		reporter.fatal("Exiting")
	}
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
