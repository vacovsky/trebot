import calpoly.golang.Go
import calpoly.sonarqube.SonarQube
import calpoly.artifact.Artifact


def go = new Go(env, steps)
def art = new Artifact(env, steps)
def sq = new SonarQube(env, steps)

node {

    stage('Preparation') {
        git repo
    }

    stage('Install Go and Dependencies') {
        def deps = [
            "github.com/go-chat-bot/bot",
                "github.com/vacoj/trebot/trivia",
                "github.com/nlopes/slack",
                "github.com/olekukonko/tablewriter"
        ] as String[]

        go.InstallDeps(deps) 
    }

    stage('Run Unit Tests') {
        go.Test()
    }

    stage('Static Code Analysis') {
        String inclusions = '*.go,trivia/*.go'
        String exclusions = ''
        String testInclusions = '*_test.go,trivia/*_test.go'
        sq.GoScan(inclusions, exclusions, testInclusions, sonarqubeServer, sonarqubeKey, project)
    }
    
    stage('Build Code Artifact') {
        go.Build(name)
    }
   
    stage("Quality Gate") {

    }

    stage('Branch Merge') {

    }

    stage('Store Artifact') {
        art.Publish(arifactoryUser, artifactoryPassword, artifactoryServer, project, name)
    }

    stage('Build Environment-specific Configuration Files') {

    }

    stage('Fact Discovery') {

    }

    stage('Deploy To Development Environment') {
       
    }

}
