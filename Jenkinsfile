pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building the payment service application'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common'
                sh 'bazel build //...'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing the payment service application'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common'
                sh 'bazel test //...'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying the payment service application'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common'
                sh 'bazel run //cmd/server:server'
            }
        }
    }
}