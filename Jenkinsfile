pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building the payment service application'
                make build
            }
        }
        stage('Test') {
            steps {
                echo 'Testing the payment service application'
                make test
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying the payment service application'
                make run
            }
        }
    }
}