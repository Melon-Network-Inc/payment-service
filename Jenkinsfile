pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building the payment service application'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel build //...'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing the payment service application'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel test //...'
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