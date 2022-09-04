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
                sh 'export GOPRIVATE=github.com/Melon-Network-Inc/common && make test'
            }
        }
        stage('Cleanup') {
            agent any
            when { branch "main" }
            steps {
                script {
                    try {
                        echo 'New release is approved. Clean up previous release.'
                        sh 'screen -XS payment-host quit'
                    } catch (e) {
                        echo 'No need to clean up and proceed to the Release stage.'
                    }
                }
            }
        }
        stage('Release') {
            agent any
            when { branch "main" }
            steps {
                echo 'Deploying the payment service application to Production.'
                sh 'export JENKINS_NODE_COOKIE=dontKillMe; screen -S payment-host  -d -m -c /dev/null -- sh -c "export GOPRIVATE=github.com/Melon-Network-Inc/common; make staging; exec sh"'
            }
        }
    }
    environment {
        EMAIL_TO = 'michaelzhou@melonnetwork.io'
    }
    post {
        success {
            emailext mimeType: 'text/html',
            body: 'Check console output at $BUILD_URL to view the results.', 
            to: "${EMAIL_TO}", 
            subject: 'Build Success in Jenkins: $PROJECT_NAME - #$BUILD_NUMBER'
        }
        failure {
            emailext mimeType: 'text/html',
            body: 'Check console output at $BUILD_URL to view the results.', 
            to: "${EMAIL_TO}", 
            subject: 'Build failed in Jenkins: $PROJECT_NAME - #$BUILD_NUMBER'
        }
        unstable {
            emailext mimeType: 'text/html',
            body: 'Check console output at $BUILD_URL to view the results.', 
            to: "${EMAIL_TO}", 
            subject: 'Unstable build in Jenkins: $PROJECT_NAME - #$BUILD_NUMBER'
        }
        changed {
            emailext mimeType: 'text/html',
            body: 'Check console output at $BUILD_URL to view the results.', 
            to: "${EMAIL_TO}", 
            subject: 'Jenkins build is back to normal: $PROJECT_NAME - #$BUILD_NUMBER'
        }
    }
}