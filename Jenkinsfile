pipeline {
    agent any

    options {
        disableConcurrentBuilds()
        timestamps()
    }

    parameters {
        choice(
            name: 'SERVICE',
            choices: [
                'adservice',
                'cartservice',
                'checkoutservice',
                'currencyservice',
                'emailservice',
                'frontend',
                'loadgenerator',
                'paymentservice',
                'productcatalogservice',
                'recommendationservice',
                'shippingservice'
            ],
            description: 'Service to test, scan, and build'
        )
        booleanParam(
            name: 'PUBLISH_IMAGE',
            defaultValue: false,
            description: 'Push the validated image to the configured registry'
        )
    }

    environment {
        REGISTRY_NAMESPACE = 'tnt850910'
    }

    stages {
        stage('Prepare') {
            steps {
                script {
                    env.SERVICE_DIR = "services/${params.SERVICE}"
                    env.BUILD_CONTEXT = params.SERVICE == 'cartservice' ? "${env.SERVICE_DIR}/src" : env.SERVICE_DIR
                    env.IMAGE = "${env.REGISTRY_NAMESPACE}/${params.SERVICE}:1.${env.BUILD_NUMBER}"
                }
            }
        }

        stage('Test service') {
            steps {
                sh './scripts/test-service.sh "$SERVICE"'
            }
        }

        stage('Scan source') {
            steps {
                sh 'grype "dir:$SERVICE_DIR" --config .grype.yaml --only-fixed --fail-on low'
            }
        }

        stage('Build image') {
            steps {
                sh 'docker build --tag "$IMAGE" "$BUILD_CONTEXT"'
            }
        }

        stage('Scan image') {
            steps {
                sh 'grype "$IMAGE" --config .grype.yaml --only-fixed --fail-on low'
            }
        }

        stage('Publish image') {
            when {
                expression { params.PUBLISH_IMAGE }
            }
            steps {
                script {
                    withDockerRegistry(credentialsId: 'docker-cred', toolName: 'docker') {
                        sh 'docker push "$IMAGE"'
                    }
                }
            }
        }
    }
}
