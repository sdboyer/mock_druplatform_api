package acquia

import (
  "github.com/gorilla/mux"
  "net/http"
)

type AcquiaServerState struct {
  Users map[string]string
}

func NewRouter() *mux.Router {
  router := mux.NewRouter()

  // Drush
  router.HandleFunc("/me/drushrc", drushrcHandler).Methods("GET")
  // Tasks
  router.HandleFunc("/sites/{site}/tasks", taskListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/tasks/{task}", taskRecordHandler).Methods("GET")
  // Domains
  router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", domainDeleteHandler).Methods("DELETE")
  router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}/cache", domainVarnishPurgeHandler).Methods("DELETE")
  router.HandleFunc("/sites/{site}/envs/{env}/domains", domainListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", domainRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/domains/{domain}", domainAddHandler).Methods("POST")
  // Servers
  router.HandleFunc("/sites/{site}/envs/{env}/servers", serverEnvListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/servers/{server}", serverRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/servers/{server}/php-procs", serverPhpProcsHandler).Methods("GET")
  // SSH Keys
  router.HandleFunc("/sites/{site}/sshkeys/{sshkeyid}", sshkeyDeleteHandler).Methods("DELETE")
  router.HandleFunc("/sites/{site}/sshkeys", sshkeyListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/sshkeys/{sshkeyid}", sshkeyRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/sshkeys", sshkeyAddHandler).Methods("POST")
  // Databases
  router.HandleFunc("/sites/{site}/dbs/{db}", dbDeleteHandler).Methods("DELETE")
  router.HandleFunc("/sites/{site}/dbs", dbListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/dbs", dbCreateHandler).Methods("POST")
  router.HandleFunc("/sites/{site}/dbs/{db}", dbRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs", dbEnvInstanceListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}", dbEnvInstanceRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups", dbEnvInstanceBackupListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups", dbEnvInstanceBackupCreateEHandler).Methods("POST")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}", dbEnvInstanceBackupRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}", dbEnvInstanceBackupDeleteHandler).Methods("DELETE")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}/download", dbEnvInstanceBackupDownloadHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/dbs/{db}/backups/{backup}/restore", dbEnvInstanceBackupRestoreHandler).Methods("POST")
  // Workflow
  router.HandleFunc("/sites/{site}/code-deploy/{source}/{target}", workflowCrossEnvCodeDeployHandler).Methods("POST")
  router.HandleFunc("/sites/{site}/dbs/{db}/db-copy/{source}/{target}", workflowCrossEnvDbCopyHandler).Methods("POST")
  router.HandleFunc("/sites/{site}/domain-move/{source}/{target}", workflowMoveDomainAcrossEnvsHandler).Methods("POST")
  router.HandleFunc("/sites/{site}/envs/{env}/code-deploy", workflowCodeDeployHandler).Methods("POST")
  router.HandleFunc("/sites/{site}/files-copy/{source}/{target}", workflowCopyFilesAcrossEnvsHandler).Methods("POST")
  // VCS users
  router.HandleFunc("/sites/{site}/svnusers", vcsusersListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/svnusers/{svnuserid}", vcsusersRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/svnusers/{svnuserid}", vcsusersDeleteHandler).Methods("DELETE")
  router.HandleFunc("/sites/{site}/svnusers/{username}", vcsusersDeleteHandler).Methods("POST")
  // Sites and Environments
  router.HandleFunc("/sites", sitesListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}", sitesRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs", sitesEnvListHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}", sitesEnvRecordHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/logstream", sitesEnvLogstreamHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/install/{type}", sitesEnvInstallHandler).Methods("GET")
  router.HandleFunc("/sites/{site}/envs/{env}/livedev/{action}", sitesLivedevHandler).Methods("GET")

  return router
}

func ServerListHandler(w http.ResponseWriter, r *http.Request) {

}
