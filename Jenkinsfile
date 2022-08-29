pipeline {
    agent none

    stages {
        stage('Build') {
            agent any
            steps {
                echo 'Run bazel build on payment service target'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel build //...'
            }
        }
        stage('Test') {
            agent any
            steps {
                echo 'Run bazel test on payment service target'
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && bazel test //...'
            }
        }
        stage('Cleanup') {
            agent any
            steps {
                echo 'New release is approved. Clean up previous release.'
                sh 'screen -XS payment-host quit'
            }
        }
        stage('Release') {
            agent any
            steps {
                echo 'Deploying the payment service application to Production.'
                sh 'screen -S payment-host  -d -m -c /dev/null -- sh -c "cd ~/payment-service; export GOPRIVATE=github.com/Melon-Network-Inc/common; make run; exec sh"'
                sh 'JENKINS_NODE_COOKIE=dontKillMe ./start.sh'
            }
        }
    }
}