package acquia

import (
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type AcquiaServerState struct {
	*http.Server  `json:"-"`
	Subscriptions []*Subscription `json:"subscriptions"`
	Tasks         TaskList        `json:"tasks"`
	Users         []*User         `json:"users"`
}

func (ss *AcquiaServerState) Version() string {
	return "1.0"
}

// Starts up and serves an httpd of this API server, binding to the provided address.
//
// This function blocks; it should typically be called in a goroutine.
func (ss *AcquiaServerState) Serve(l net.Listener) {
	// The Addr prop shouldn't actually be used, but set it to avoid triggering defaults
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(ss.Router())

	srv := &http.Server{Addr: l.Addr().String(), Handler: n}
	srv.Serve(l) // blocks inside here
}

type Subscription struct {
	Name         string         `json:"name"`
	Environments []*Environment `json:"environments"`
	Databases    []*Database    `json:"databases"`
	Users        []*User        `json:"users"`
}

type Database struct {
	Name string
}

type Environment struct {
	Name        string   `json:"name"`
	Domains     []string `json:"domains"`
	CodeVersion string   `json:"vcs_path"`
}

type User struct {
	Name     string
	Email    string
	Password string
	Key      string
}

type Task struct {
	Id          int    `json:"name"`
	Created     int    `json:"created"`
	Started     int    `json:"started"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Logs        string `json:"logs"`
	Queue       string `json:"queue"`
	Result      string `json:"result"`
	Sender      string `json:"sender"`
	State       int    `json:"state"`
}

type TaskList []*Task

func (tl TaskList) AddTask() *Task {
	t := &Task{Id: len(tl)}
	tl = append(tl, t)
	return t
}

// Creates a new AcquiaServerState, which acts as the basis for a mock API endpoint.
func NewServerInstance(subname string) *AcquiaServerState {
	aqs := &AcquiaServerState{
		Subscriptions: make([]*Subscription, 0),
		Tasks:         make([]*Task, 0),
	}

	// Make a task to take up the zero index
	aqs.Tasks.AddTask()

	aqs.Subscriptions = append(aqs.Subscriptions, NewSubscription(subname))
	return aqs
}

func NewSubscription(name string) *Subscription {
	s := &Subscription{
		Name: name,
		Environments: []*Environment{
			&Environment{Name: "dev", Domains: make([]string, 0), CodeVersion: "tags/WELCOME"},
			&Environment{Name: "test", Domains: make([]string, 0), CodeVersion: "tags/WELCOME"},
			&Environment{Name: "prod", Domains: make([]string, 0), CodeVersion: "tags/WELCOME"},
		},
	}

	return s
}

func (ss *AcquiaServerState) Router() *mux.Router {
	router := mux.NewRouter()

	// Drush
	router.HandleFunc("/me/drushrc", ss.drushRcHandler).Methods("GET")
	// Tasks
	router.HandleFunc("/sites/{site}/tasks", ss.taskListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/tasks/{task}", ss.taskRecordHandler).Methods("GET")
	// Domains
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", ss.domainDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}/cache", ss.domainVarnishPurgeHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/envs/{env}/domains", ss.domainListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", ss.domainRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", ss.domainAddHandler).Methods("POST")
	// Servers
	router.HandleFunc("/sites/{site}/envs/{env}/servers", ss.serverEnvListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/servers/{server}", ss.serverRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/servers/{server}/php-procs", ss.serverPhpProcsHandler).Methods("GET")
	// SSH Keys
	router.HandleFunc("/sites/{site}/sshkeys/{sshkeyid}", ss.sshkeyDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/sshkeys", ss.sshkeyListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/sshkeys/{sshkeyid}", ss.sshkeyRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/sshkeys", ss.sshkeyAddHandler).Methods("POST")
	// Databases
	router.HandleFunc("/sites/{site}/dbs/{db}", ss.dbDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/dbs", ss.dbListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/dbs", ss.dbCreateHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/dbs/{db}", ss.dbRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs", ss.dbEnvInstanceListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}", ss.dbEnvInstanceRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups", ss.dbEnvInstanceBackupListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups", ss.dbEnvInstanceBackupCreateHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}", ss.dbEnvInstanceBackupRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}", ss.dbEnvInstanceBackupDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}/download", ss.dbEnvInstanceBackupDownloadHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}/restore", ss.dbEnvInstanceBackupRestoreHandler).Methods("POST")
	// Workflow
	router.HandleFunc("/sites/{site}/code-deploy/{source}/{target}", ss.workflowCrossEnvCodeDeployHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/dbs/{db}/db-copy/{source}/{target}", ss.workflowCrossEnvDbCopyHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/domain-move/{source}/{target}", ss.workflowMoveDomainAcrossEnvsHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/envs/{env}/code-deploy", ss.workflowCodeDeployHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/files-copy/{source}/{target}", ss.workflowCopyFilesAcrossEnvsHandler).Methods("POST")
	// VCS users
	router.HandleFunc("/sites/{site}/svnusers", ss.vcsusersListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/svnusers/{svnuserid}", ss.vcsusersRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/svnusers/{svnuserid}", ss.vcsusersDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/svnusers/{username}", ss.vcsusersCreateHandler).Methods("POST")
	// Sites and Environments
	router.HandleFunc("/sites", ss.sitesListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}", ss.sitesRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs", ss.sitesEnvListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}", ss.sitesEnvRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/logstream", ss.sitesEnvLogstreamHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/install/{type}", ss.sitesEnvInstallHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/livedev/{action}", ss.sitesLivedevHandler).Methods("GET")

	return router
}

func (ss *AcquiaServerState) taskListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) drushRcHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) taskRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) domainDeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) domainVarnishPurgeHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) domainListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) domainRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) domainAddHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) serverEnvListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) serverRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) serverPhpProcsHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sshkeyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sshkeyListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sshkeyRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sshkeyAddHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbDeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbCreateHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceBackupListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceBackupCreateHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceBackupRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceBackupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceBackupDownloadHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) dbEnvInstanceBackupRestoreHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) workflowCrossEnvCodeDeployHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) workflowCrossEnvDbCopyHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) workflowMoveDomainAcrossEnvsHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) workflowCodeDeployHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) workflowCopyFilesAcrossEnvsHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) vcsusersListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) vcsusersRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) vcsusersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) vcsusersCreateHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sitesListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sitesRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sitesEnvListHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sitesEnvRecordHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sitesEnvLogstreamHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sitesEnvInstallHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}

func (ss *AcquiaServerState) sitesLivedevHandler(w http.ResponseWriter, r *http.Request) {
	panic("not yet implemented")
}
