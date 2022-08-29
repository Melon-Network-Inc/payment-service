pipeline {
    agent none

    stages {
        stage('Build') {
            agent any
            when { branch "main" }
            steps {
                echo 'Run bazel build on payment service target'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel build //...'
            }
        }
        stage('Test') {
            agent any
            when { branch "main" }
            steps {
                echo 'Run bazel test on payment service target'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel test //...'
            }
        }
        stage('Cleanup') {
            agent any
            when { branch "main" }
            steps {
                echo 'New release is approved. Clean up previous release.'
                sh 'screen -XS payment-host quit'
            }
        }
        stage('Release') {
            agent any
            when { branch "main" }
            steps {
                echo 'Deploying the payment service application to Production.'
                sh 'export JENKINS_NODE_COOKIE=dontKillMe; screen -S payment-host  -d -m -c /dev/null -- sh -c "export GOPRIVATE=github.com/Melon-Network-Inc/common; make run; exec sh"'
            }
        }
    }
}