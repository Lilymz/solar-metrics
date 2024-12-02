pipeline {
    agent any
    environment{
        GIT_URL = "https://github.com/Lilymz/solar-metrics.git"
        GIT_CRED = "solar"
        APP_NAME = "solar-metrics"
        VERSION = "V1.0.0"
        GO111MODULE = 'on'
        GOPROXY = 'https://goproxy.cn,direct'
        GOROOT = '/usr/local/go'
        GOPATH = '/usr/local/go/library'
        PATH = "${env.PATH}:${GOROOT}/bin:${GOPATH}/bin"
    }
    parameters {
        string(name: 'domain', defaultValue: 'localhost', description: 'mq地址')
        string(name: 'user', defaultValue: 'admin', description: 'The username for authentication')
        string(name: 'password', defaultValue: 'admin', description: 'The password for authentication')
    }
    stages {
        stage("Dockerfile Before Handler") {
            steps {
                script {
                    sh '''
                        docker rm -f "${APP_NAME}" || true
                        docker rmi "${APP_NAME}":${VERSION} || true
                        docker image prune -f
                        docker network prune -f
                    '''
                }
            }
        }
        stage("Checkout") {
            steps {
                script {
                     retry(3) {
                        git "${GIT_URL}"
                        def amqpUrl = "amqp://${params.user}:${params.password}@${params.domain}:5672"
                        sh """
                        sed -i 's|amqp://admin:admin@localhost:5672|${amqpUrl}|' config/solar-metric.yml
                        """
                     }
                }
            }
        }
        stage("Build") {
            steps {
                script {
                    sh '''
                        go mod tidy
                        go build -o "${APP_NAME}" ./internal/.
                    '''
                }
            }
        }
        stage("Build Image") {
            steps {
                script {
                    sh(script: "docker build --no-cache --compress --pull --tag ${APP_NAME}:${VERSION} .")
                }
            }
        }
        stage("Deploy") {
            steps {
                script {
                    sh 'docker compose up -d'
                }
            }
        }
    }
    post {
        success {
         echo "Build success"
        // 构建成功，将当前镜像推送到私服仓库(可选)
        // 发送构建成功通知（可选）
        }
        failure {
            script {
                // 处理构建失败的情况
                echo "Build failed, check the logs."
            }
        }
    }
}
