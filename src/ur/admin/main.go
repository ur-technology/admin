package main

import (
	"fmt"
	"ur/database"

	"github.com/yourheropaul/inj"
)

type application struct {
	Environment            environmentMap          `inj:""`
	Config                 configurer              `inj:""`
	DatabaseSessionFactory database.SessionFactory `inj:""`
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
}
