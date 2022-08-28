pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building the payment service application'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel build //...'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying the payment service application'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel run //cmd/server:server'
            }
        }
    }
}