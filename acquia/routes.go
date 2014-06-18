package acquia

import (
	"github.com/gorilla/mux"
	"net/http"
)

type AcquiaServerState struct {
	Subscriptions []*Subscription
	Tasks TaskList
}

func (ss *AcquiaServerState) Version() string {
	return "1.0"
}

type Subscription struct {
	Name string
	Environments []*Environment
	Databases []*Database
	Users []*User
}

type Database struct {
	Name string
}

type Environment struct {
	Name string
	Domains []string
	CodeVersion string
}

type User struct {
	Name     string
	Email    string
	Password string
	Key      string
}

type Task struct {
	Id int
	Created int
	Started int
	Description string
	Completed bool
	Logs string
	Queue string
	Result string
	Sender string
	State int
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
		Tasks: make([]*Task, 0),
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

func NewRouter(srv *AcquiaServerState) *mux.Router {
	router := mux.NewRouter()

	// Drush
	router.HandleFunc("/me/drushrc", srv.drushRcHandler).Methods("GET")
	// Tasks
	router.HandleFunc("/sites/{site}/tasks", srv.taskListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/tasks/{task}", srv.taskRecordHandler).Methods("GET")
	// Domains
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", srv.domainDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}/cache", srv.domainVarnishPurgeHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/envs/{env}/domains", srv.domainListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", srv.domainRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", srv.domainAddHandler).Methods("POST")
	// Servers
	router.HandleFunc("/sites/{site}/envs/{env}/servers", srv.serverEnvListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/servers/{server}", srv.serverRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/servers/{server}/php-procs", srv.serverPhpProcsHandler).Methods("GET")
	// SSH Keys
	router.HandleFunc("/sites/{site}/sshkeys/{sshkeyid}", srv.sshkeyDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/sshkeys", srv.sshkeyListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/sshkeys/{sshkeyid}", srv.sshkeyRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/sshkeys", srv.sshkeyAddHandler).Methods("POST")
	// Databases
	router.HandleFunc("/sites/{site}/dbs/{db}", srv.dbDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/dbs", srv.dbListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/dbs", srv.dbCreateHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/dbs/{db}", srv.dbRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs", srv.dbEnvInstanceListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}", srv.dbEnvInstanceRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups", srv.dbEnvInstanceBackupListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups", srv.dbEnvInstanceBackupCreateHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}", srv.dbEnvInstanceBackupRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}", srv.dbEnvInstanceBackupDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}/download", srv.dbEnvInstanceBackupDownloadHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}/restore", srv.dbEnvInstanceBackupRestoreHandler).Methods("POST")
	// Workflow
	router.HandleFunc("/sites/{site}/code-deploy/{source}/{target}", srv.workflowCrossEnvCodeDeployHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/dbs/{db}/db-copy/{source}/{target}", srv.workflowCrossEnvDbCopyHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/domain-move/{source}/{target}", srv.workflowMoveDomainAcrossEnvsHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/envs/{env}/code-deploy", srv.workflowCodeDeployHandler).Methods("POST")
	router.HandleFunc("/sites/{site}/files-copy/{source}/{target}", srv.workflowCopyFilesAcrossEnvsHandler).Methods("POST")
	// VCS users
	router.HandleFunc("/sites/{site}/svnusers", srv.vcsusersListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/svnusers/{svnuserid}", srv.vcsusersRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/svnusers/{svnuserid}", srv.vcsusersDeleteHandler).Methods("DELETE")
	router.HandleFunc("/sites/{site}/svnusers/{username}", srv.vcsusersCreateHandler).Methods("POST")
	// Sites and Environments
	router.HandleFunc("/sites", srv.sitesListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}", srv.sitesRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs", srv.sitesEnvListHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}", srv.sitesEnvRecordHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/logstream", srv.sitesEnvLogstreamHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/install/{type}", srv.sitesEnvInstallHandler).Methods("GET")
	router.HandleFunc("/sites/{site}/envs/{env}/livedev/{action}", srv.sitesLivedevHandler).Methods("GET")

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

