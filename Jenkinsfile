pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building the payment service application'
                sh 'go build cmd/server/main.go'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying the payment service application'
                sh 'go run cmd/server/main.go'
            }
        }
    }
}