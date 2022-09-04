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
                echo 'New release is approved. Clean up previous release.'
                sh 'screen -XS payment-host quit'
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE')
                {
                    echo 'No need to clean up and proceed to the Release stage.'
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
}