package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"psql/internal/configuration"
	psqldb "psql/internal/db"
	"strings"
)

var FSO configuration.OsFileSystem

func main() {
	var (
		domain string // = "SIGMA"
		envKey string // = "psi"
		scenarioName string // = "UVS_SCENARIO_ONE"
		ctx = context.Background()

	)
	logger := zap.NewExample()
	// Не может быть ошибки т.к. работаем с stdout
	//nolint: errcheck
	logger.Sync()

	// ReplaceGlobals replaces the global Logger and SugaredLogger, and returns a
	// function to restore the original values. It's safe for concurrent use.
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	flag.StringVar(&scenarioName, "scenarioName", "", "rlm scenario name")
	flag.StringVar(&envKey, "envKey", "", "envKey name: psi or master")
	flag.StringVar(&domain, "domain", "", "domain: alpha or sigma")
	flag.Parse()

	if len(scenarioName)==0 && len(envKey)==0 && len(domain)==0 {
		zap.S().Panicw("Empty flags", "scenarioName", scenarioName, "envKey", envKey, "domain", domain)
	}

	config, err := configuration.ReadConfig("/opt/yevgen/githome/publicgit/gb/gb_go_psql/psql/cmd/config.yaml", FSO)
	if err != nil {
		zap.S().Panicw("Read configuration file error.", "err", err)
	}

	srv := config["COMMON"].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["server"].(string)
	port := config["COMMON"].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["port"].(string)
	db := config["COMMON"].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["db"].(string)
	table := config["COMMON"].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["table"].(string)
	cred_prefix := config["COMMON"].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["cred_prefix"].(string)
	// j_id := config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_id"].(string)
	zap.S().Infow("Psql db parameters","srv", srv,"port",port, "db", db,"table",table)
	fmt.Println(config)


	registryJobInfo := psqldb.RegistryJobInfo{
		S_Protocol: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["s_protocol"].(string),
		S_Host: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["s_host"].(string),
		S_Domain: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["s_domain"].(string),
		S_Jenkins_path: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["s_jenkins_path"].(string),
		J_Build_type: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_build_type"].(string),
		J_Token: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_token"].(string),
		J_Notes: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_notes"].(string),
		J_Mdata: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_mdata"].(string),
		J_Published: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_published"].(string),
		J_Version: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_version"].(string),
		J_Version_Tag: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_version_tag"].(string),
		S_User_Name: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["s_user_name"].(string),
		S_User_Token: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["s_user_token"].(string),
	}

	//J_Id: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_id"].(string),
	//J_Path: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_path"].(string),
	//J_Job_name: config[strings.ToUpper(scenarioName)].(map[interface{}]interface{})[strings.ToUpper(domain)].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[envKey].(map[interface{}]interface{})["j_job_name"].(string),



	psql, err := psqldb.NewPool(ctx, srv, port, db, cred_prefix)
	if err != nil {
		zap.S().Panicw("Couldnot connect to th registry", "srv",srv,"port",port,"db",db,"err", err)
	}
	defer psql.Close()

	ids, err := psql.SearchJID(ctx, table,j_id)
	if err != nil {
		zap.S().Errorw("Search error", "err", err)
	}
	if len(ids) > 0 {
		zap.S().Infow("This jid was found in registry. Exit.", "id", ids)
		return
	} else {
		resId, err := psql.Insert(ctx, table, registryJobInfo)
		if err != nil {
			zap.S().Panicw("Couldnot insert into registry", "err", err)
		}
		zap.S().Infow("New registry entry.", "id", resId)
	}
}

