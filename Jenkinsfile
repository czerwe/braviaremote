node ('docker') {
    stage('prepare Env') {
        deleteDir()
    }
    stage('build application') {
        withDockerContainer(image: 'golang:1.9.0') {

            dir ('src/braviaremote'){
                git credentialsId: '1db449c1-a2af-4972-b32b-f7cdd65473e8', url: 'git@github.com:czerwe/braviaremote.git'
            }
            sh 'export GOPATH=$(pwd) && cd src/braviaremote && make dependencys && make build'
        }
    }
    stage('Build Docker Image ') {
        dir ('src/braviaremote'){
            docker.build("braviaremote:${BUILD_NUMBER}")
        }
    }
}