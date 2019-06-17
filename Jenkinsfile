pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'echo build'
                sh 'go build -o client.exe'
                sh 'cd Services && go build -o services.exe'
            }
        }
        stage('Test') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'echo run tests...'
                sh 'cd Services/cinemahall && go test'
                sh 'cd Services/movie && go test'
                sh 'cd Services/reseration && go test'
                sh 'cd Services/showing && go test'
                sh 'cd Services/user && go test'
            }
        }
        stage('Lint') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }   
            steps {
                sh 'golangci-lint run --deadline 20m --enable-all'
            }
        }
        stage('Build Docker Image') {
            agent any
            steps {
                sh "echo build docker"
            }
        }
    }
}
